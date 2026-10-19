package armazenamento

import (
	"context"
	"errors"
	"time"

	"piaps/backend/internal/dominio"
)

var (
	ErrNaoEncontrado = errors.New("registro não encontrado")
	ErrConflito      = errors.New("registro em estado conflitante")
)

type Rascunho struct {
	ID           string
	Etapa        int
	Respostas    dominio.Respostas
	CriadoEm     time.Time
	AtualizadoEm time.Time
}

type Concluida struct {
	ID          string
	Respostas   dominio.Respostas
	Resultado   dominio.Resultado
	ConcluidaEm time.Time
}

type Armazenamento interface {
	RegistrarDispositivo(ctx context.Context, id string) error
	CriarRascunho(ctx context.Context, dispositivoID string, respostas dominio.Respostas) (Rascunho, error)
	ObterRascunho(ctx context.Context, dispositivoID, id string) (Rascunho, error)
	SalvarRascunho(ctx context.Context, dispositivoID, id string, etapa int, respostas dominio.Respostas) (Rascunho, error)
	ApagarRascunho(ctx context.Context, dispositivoID, id string) error
	Concluir(ctx context.Context, dispositivoID, id string, respostas dominio.Respostas, resultado dominio.Resultado) (Concluida, error)
	Reabrir(ctx context.Context, dispositivoID, id string) (Rascunho, error)
	ListarRascunhos(ctx context.Context, dispositivoID string) ([]Rascunho, error)
	ListarHistorico(ctx context.Context, dispositivoID string) ([]Concluida, error)
	ObterConcluida(ctx context.Context, dispositivoID, id string) (Concluida, error)
	ApagarConcluida(ctx context.Context, dispositivoID, id string) error
	ApagarHistorico(ctx context.Context, dispositivoID string) error
}
