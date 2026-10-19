package armazenamento

import (
	"context"
	"embed"
	"fmt"
	"io/fs"
	"sort"
	"strings"

	"github.com/jackc/pgx/v5/pgxpool"
)

//go:embed migracoes/*.sql
var migracoesFS embed.FS

func Migrar(ctx context.Context, pool *pgxpool.Pool) error {
	if _, err := pool.Exec(ctx, `
		CREATE TABLE IF NOT EXISTS schema_migrations (
			versao TEXT PRIMARY KEY,
			aplicada_em TIMESTAMPTZ NOT NULL DEFAULT now()
		)`); err != nil {
		return fmt.Errorf("criar schema_migrations: %w", err)
	}

	entradas, err := fs.ReadDir(migracoesFS, "migracoes")
	if err != nil {
		return err
	}
	nomes := make([]string, 0, len(entradas))
	for _, e := range entradas {
		if !e.IsDir() && strings.HasSuffix(e.Name(), ".sql") {
			nomes = append(nomes, e.Name())
		}
	}
	sort.Strings(nomes)

	tx, err := pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	if _, err := tx.Exec(ctx, `SELECT pg_advisory_xact_lock($1)`, 0x7069617073); err != nil {
		return fmt.Errorf("lock de migração: %w", err)
	}

	for _, nome := range nomes {
		var existe bool
		if err := tx.QueryRow(ctx,
			`SELECT EXISTS (SELECT 1 FROM schema_migrations WHERE versao = $1)`, nome).Scan(&existe); err != nil {
			return err
		}
		if existe {
			continue
		}
		sql, err := migracoesFS.ReadFile("migracoes/" + nome)
		if err != nil {
			return err
		}
		if _, err := tx.Exec(ctx, string(sql)); err != nil {
			return fmt.Errorf("migração %s: %w", nome, err)
		}
		if _, err := tx.Exec(ctx,
			`INSERT INTO schema_migrations (versao) VALUES ($1)`, nome); err != nil {
			return err
		}
	}
	return tx.Commit(ctx)
}
