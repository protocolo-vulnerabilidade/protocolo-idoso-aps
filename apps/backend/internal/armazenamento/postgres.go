package armazenamento

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"piaps/backend/internal/dominio"
)

type Postgres struct {
	pool *pgxpool.Pool
}

func NovoPostgres(pool *pgxpool.Pool) *Postgres {
	return &Postgres{pool: pool}
}

func (p *Postgres) RegistrarDispositivo(ctx context.Context, id string) error {
	_, err := p.pool.Exec(ctx, `INSERT INTO dispositivos (id) VALUES ($1) ON CONFLICT (id) DO NOTHING`, id)
	return err
}

func codificarRespostas(r dominio.Respostas) []byte {
	b, err := json.Marshal(r)
	if err != nil {
		panic(fmt.Sprintf("respostas impossíveis de serializar: %v", err))
	}
	return b
}

func decodificarRespostas(b []byte) (dominio.Respostas, error) {
	var r dominio.Respostas
	if err := json.Unmarshal(b, &r); err != nil {
		return r, err
	}
	return r, nil
}

func paraInt16(vals []int) []int16 {
	out := make([]int16, len(vals))
	for i, v := range vals {
		out[i] = int16(v)
	}
	return out
}

func paraInt(vals []int16) []int {
	out := make([]int, len(vals))
	for i, v := range vals {
		out[i] = int(v)
	}
	return out
}

const colunasRascunho = `id, etapa, respostas, criado_em, atualizado_em`

func escanearRascunho(row pgx.Row) (Rascunho, error) {
	var r Rascunho
	var respostasJSON []byte
	if err := row.Scan(&r.ID, &r.Etapa, &respostasJSON, &r.CriadoEm, &r.AtualizadoEm); err != nil {
		return r, err
	}
	respostas, err := decodificarRespostas(respostasJSON)
	if err != nil {
		return r, err
	}
	r.Respostas = respostas
	return r, nil
}

func (p *Postgres) CriarRascunho(ctx context.Context, dispositivoID string, respostas dominio.Respostas) (Rascunho, error) {
	row := p.pool.QueryRow(ctx, `
		INSERT INTO avaliacoes (id, dispositivo_id, situacao, etapa, respostas)
		VALUES ($1, $2, 'rascunho', 0, $3::jsonb)
		RETURNING `+colunasRascunho,
		uuid.NewString(), dispositivoID, codificarRespostas(respostas))
	return escanearRascunho(row)
}

func (p *Postgres) ObterRascunho(ctx context.Context, dispositivoID, id string) (Rascunho, error) {
	row := p.pool.QueryRow(ctx, `
		SELECT `+colunasRascunho+` FROM avaliacoes
		WHERE id = $1 AND dispositivo_id = $2 AND situacao = 'rascunho'`, id, dispositivoID)
	r, err := escanearRascunho(row)
	if errors.Is(err, pgx.ErrNoRows) {
		return Rascunho{}, ErrNaoEncontrado
	}
	return r, err
}

func (p *Postgres) SalvarRascunho(ctx context.Context, dispositivoID, id string, etapa int, respostas dominio.Respostas) (Rascunho, error) {
	row := p.pool.QueryRow(ctx, `
		UPDATE avaliacoes
		SET etapa = $3, respostas = $4::jsonb, atualizado_em = now()
		WHERE id = $1 AND dispositivo_id = $2 AND situacao = 'rascunho'
		RETURNING `+colunasRascunho,
		id, dispositivoID, etapa, codificarRespostas(respostas))
	r, err := escanearRascunho(row)
	if errors.Is(err, pgx.ErrNoRows) {
		return Rascunho{}, ErrNaoEncontrado
	}
	return r, err
}

func (p *Postgres) ApagarRascunho(ctx context.Context, dispositivoID, id string) error {
	tag, err := p.pool.Exec(ctx, `
		DELETE FROM avaliacoes
		WHERE id = $1 AND dispositivo_id = $2 AND situacao = 'rascunho'`, id, dispositivoID)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return ErrNaoEncontrado
	}
	return nil
}

const colunasConcluida = `id, respostas, escore_ivcf, classe_ivcf, escore_easfi, classe_easfi,
	nivel_matriz, nivel_final, COALESCE(sentinelas_positivas, '{}'), concluida_em`

func escanearConcluida(row pgx.Row) (Concluida, error) {
	var c Concluida
	var respostasJSON []byte
	var escoreIvcf, escoreEasfi, nivelMatriz, nivelFinal int16
	var classeIvcf, classeEasfi string
	var sentinelas []int16
	if err := row.Scan(&c.ID, &respostasJSON, &escoreIvcf, &classeIvcf, &escoreEasfi, &classeEasfi,
		&nivelMatriz, &nivelFinal, &sentinelas, &c.ConcluidaEm); err != nil {
		return c, err
	}
	respostas, err := decodificarRespostas(respostasJSON)
	if err != nil {
		return c, err
	}
	c.Respostas = respostas
	c.Resultado = dominio.Resultado{
		Ivcf:                int(escoreIvcf),
		ClasseIVCF:          dominio.ClasseIVCF(classeIvcf),
		Easfi:               int(escoreEasfi),
		ClasseEASFI:         dominio.ClasseEASFI(classeEasfi),
		NivelMatriz:         dominio.Nivel(nivelMatriz),
		Nivel:               dominio.Nivel(nivelFinal),
		SentinelasPositivas: paraInt(sentinelas),
	}
	return c, nil
}

func (p *Postgres) Concluir(ctx context.Context, dispositivoID, id string, respostas dominio.Respostas, resultado dominio.Resultado) (Concluida, error) {
	row := p.pool.QueryRow(ctx, `
		UPDATE avaliacoes
		SET situacao = 'concluida',
		    respostas = $3::jsonb,
		    escore_ivcf = $4,
		    classe_ivcf = $5,
		    escore_easfi = $6,
		    classe_easfi = $7,
		    nivel_matriz = $8,
		    nivel_final = $9,
		    sentinelas_positivas = $10,
		    concluida_em = now(),
		    atualizado_em = now()
		WHERE id = $1 AND dispositivo_id = $2 AND situacao = 'rascunho'
		RETURNING `+colunasConcluida,
		id, dispositivoID, codificarRespostas(respostas),
		resultado.Ivcf, string(resultado.ClasseIVCF),
		resultado.Easfi, string(resultado.ClasseEASFI),
		int(resultado.NivelMatriz), int(resultado.Nivel),
		paraInt16(resultado.SentinelasPositivas))
	c, err := escanearConcluida(row)
	if errors.Is(err, pgx.ErrNoRows) {
		var situacao string
		err2 := p.pool.QueryRow(ctx,
			`SELECT situacao FROM avaliacoes WHERE id = $1 AND dispositivo_id = $2`, id, dispositivoID).Scan(&situacao)
		if err2 == nil && situacao == "concluida" {
			return Concluida{}, ErrConflito
		}
		return Concluida{}, ErrNaoEncontrado
	}
	return c, err
}

func (p *Postgres) Reabrir(ctx context.Context, dispositivoID, id string) (Rascunho, error) {
	row := p.pool.QueryRow(ctx, `
		UPDATE avaliacoes
		SET situacao = 'rascunho',
		    etapa = 0,
		    escore_ivcf = NULL,
		    classe_ivcf = NULL,
		    escore_easfi = NULL,
		    classe_easfi = NULL,
		    nivel_matriz = NULL,
		    nivel_final = NULL,
		    sentinelas_positivas = NULL,
		    concluida_em = NULL,
		    atualizado_em = now()
		WHERE id = $1 AND dispositivo_id = $2 AND situacao = 'concluida'
		RETURNING `+colunasRascunho, id, dispositivoID)
	r, err := escanearRascunho(row)
	if errors.Is(err, pgx.ErrNoRows) {
		return Rascunho{}, ErrNaoEncontrado
	}
	return r, err
}

func (p *Postgres) ListarRascunhos(ctx context.Context, dispositivoID string) ([]Rascunho, error) {
	rows, err := p.pool.Query(ctx, `
		SELECT `+colunasRascunho+` FROM avaliacoes
		WHERE dispositivo_id = $1 AND situacao = 'rascunho'
		ORDER BY atualizado_em DESC`, dispositivoID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var lista []Rascunho
	for rows.Next() {
		var r Rascunho
		var respostasJSON []byte
		if err := rows.Scan(&r.ID, &r.Etapa, &respostasJSON, &r.CriadoEm, &r.AtualizadoEm); err != nil {
			return nil, err
		}
		if r.Respostas, err = decodificarRespostas(respostasJSON); err != nil {
			return nil, err
		}
		lista = append(lista, r)
	}
	return lista, rows.Err()
}

func (p *Postgres) ListarHistorico(ctx context.Context, dispositivoID string) ([]Concluida, error) {
	rows, err := p.pool.Query(ctx, `
		SELECT `+colunasConcluida+` FROM avaliacoes
		WHERE dispositivo_id = $1 AND situacao = 'concluida'
		ORDER BY concluida_em DESC`, dispositivoID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var lista []Concluida
	for rows.Next() {
		var c Concluida
		var respostasJSON []byte
		var escoreIvcf, escoreEasfi, nivelMatriz, nivelFinal int16
		var classeIvcf, classeEasfi string
		var sentinelas []int16
		if err := rows.Scan(&c.ID, &respostasJSON, &escoreIvcf, &classeIvcf, &escoreEasfi, &classeEasfi,
			&nivelMatriz, &nivelFinal, &sentinelas, &c.ConcluidaEm); err != nil {
			return nil, err
		}
		if c.Respostas, err = decodificarRespostas(respostasJSON); err != nil {
			return nil, err
		}
		c.Resultado = dominio.Resultado{
			Ivcf:                int(escoreIvcf),
			ClasseIVCF:          dominio.ClasseIVCF(classeIvcf),
			Easfi:               int(escoreEasfi),
			ClasseEASFI:         dominio.ClasseEASFI(classeEasfi),
			NivelMatriz:         dominio.Nivel(nivelMatriz),
			Nivel:               dominio.Nivel(nivelFinal),
			SentinelasPositivas: paraInt(sentinelas),
		}
		lista = append(lista, c)
	}
	return lista, rows.Err()
}

func (p *Postgres) ObterConcluida(ctx context.Context, dispositivoID, id string) (Concluida, error) {
	row := p.pool.QueryRow(ctx, `
		SELECT `+colunasConcluida+` FROM avaliacoes
		WHERE id = $1 AND dispositivo_id = $2 AND situacao = 'concluida'`, id, dispositivoID)
	c, err := escanearConcluida(row)
	if errors.Is(err, pgx.ErrNoRows) {
		return Concluida{}, ErrNaoEncontrado
	}
	return c, err
}

func (p *Postgres) ApagarConcluida(ctx context.Context, dispositivoID, id string) error {
	tag, err := p.pool.Exec(ctx, `
		DELETE FROM avaliacoes
		WHERE id = $1 AND dispositivo_id = $2 AND situacao = 'concluida'`, id, dispositivoID)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return ErrNaoEncontrado
	}
	return nil
}

func (p *Postgres) ApagarHistorico(ctx context.Context, dispositivoID string) error {
	_, err := p.pool.Exec(ctx,
		`DELETE FROM avaliacoes WHERE dispositivo_id = $1 AND situacao = 'concluida'`, dispositivoID)
	return err
}
