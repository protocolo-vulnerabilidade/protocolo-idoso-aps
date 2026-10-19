package armazenamento

import (
	"context"
	"sort"
	"sync"
	"time"

	"github.com/google/uuid"

	"piaps/backend/internal/dominio"
)

type registroMemoria struct {
	id          string
	dispositivo string
	rascunho    bool
	etapa       int
	respostas   dominio.Respostas
	resultado   *dominio.Resultado
	criadoEm    time.Time
	atualizado  time.Time
	concluidaEm time.Time
}

type Memoria struct {
	mu        sync.Mutex
	registros map[string]*registroMemoria
}

func NovoMemoria() *Memoria {
	return &Memoria{registros: map[string]*registroMemoria{}}
}

func (m *Memoria) RegistrarDispositivo(ctx context.Context, id string) error {
	return nil
}

func (m *Memoria) CriarRascunho(ctx context.Context, dispositivoID string, respostas dominio.Respostas) (Rascunho, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	agora := time.Now().UTC()
	reg := &registroMemoria{
		id:          uuid.NewString(),
		dispositivo: dispositivoID,
		rascunho:    true,
		etapa:       0,
		respostas:   respostas,
		criadoEm:    agora,
		atualizado:  agora,
	}
	m.registros[reg.id] = reg
	return reg.paraRascunho(), nil
}

func (m *Memoria) buscar(dispositivoID, id string, rascunho bool) (*registroMemoria, error) {
	reg, ok := m.registros[id]
	if !ok || reg.dispositivo != dispositivoID || reg.rascunho != rascunho {
		return nil, ErrNaoEncontrado
	}
	return reg, nil
}

func (m *Memoria) ObterRascunho(ctx context.Context, dispositivoID, id string) (Rascunho, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	reg, err := m.buscar(dispositivoID, id, true)
	if err != nil {
		return Rascunho{}, err
	}
	return reg.paraRascunho(), nil
}

func (m *Memoria) SalvarRascunho(ctx context.Context, dispositivoID, id string, etapa int, respostas dominio.Respostas) (Rascunho, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	reg, err := m.buscar(dispositivoID, id, true)
	if err != nil {
		return Rascunho{}, err
	}
	reg.etapa = etapa
	reg.respostas = respostas
	reg.atualizado = time.Now().UTC()
	return reg.paraRascunho(), nil
}

func (m *Memoria) ApagarRascunho(ctx context.Context, dispositivoID, id string) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	if _, err := m.buscar(dispositivoID, id, true); err != nil {
		return err
	}
	delete(m.registros, id)
	return nil
}

func (m *Memoria) Concluir(ctx context.Context, dispositivoID, id string, respostas dominio.Respostas, resultado dominio.Resultado) (Concluida, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	reg, err := m.buscar(dispositivoID, id, true)
	if err != nil {
		if _, err2 := m.buscar(dispositivoID, id, false); err2 == nil {
			return Concluida{}, ErrConflito
		}
		return Concluida{}, err
	}
	agora := time.Now().UTC()
	reg.rascunho = false
	reg.respostas = respostas
	reg.resultado = &resultado
	reg.concluidaEm = agora
	reg.atualizado = agora
	return reg.paraConcluida(), nil
}

func (m *Memoria) Reabrir(ctx context.Context, dispositivoID, id string) (Rascunho, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	reg, err := m.buscar(dispositivoID, id, false)
	if err != nil {
		return Rascunho{}, err
	}
	reg.rascunho = true
	reg.etapa = 0
	reg.resultado = nil
	reg.concluidaEm = time.Time{}
	reg.atualizado = time.Now().UTC()
	return reg.paraRascunho(), nil
}

func (m *Memoria) ListarRascunhos(ctx context.Context, dispositivoID string) ([]Rascunho, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	var lista []Rascunho
	for _, reg := range m.registros {
		if reg.dispositivo == dispositivoID && reg.rascunho {
			lista = append(lista, reg.paraRascunho())
		}
	}
	sort.Slice(lista, func(i, j int) bool { return lista[i].AtualizadoEm.After(lista[j].AtualizadoEm) })
	return lista, nil
}

func (m *Memoria) ListarHistorico(ctx context.Context, dispositivoID string) ([]Concluida, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	var lista []Concluida
	for _, reg := range m.registros {
		if reg.dispositivo == dispositivoID && !reg.rascunho {
			lista = append(lista, reg.paraConcluida())
		}
	}
	sort.Slice(lista, func(i, j int) bool { return lista[i].ConcluidaEm.After(lista[j].ConcluidaEm) })
	return lista, nil
}

func (m *Memoria) ObterConcluida(ctx context.Context, dispositivoID, id string) (Concluida, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	reg, err := m.buscar(dispositivoID, id, false)
	if err != nil {
		return Concluida{}, err
	}
	return reg.paraConcluida(), nil
}

func (m *Memoria) ApagarConcluida(ctx context.Context, dispositivoID, id string) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	if _, err := m.buscar(dispositivoID, id, false); err != nil {
		return err
	}
	delete(m.registros, id)
	return nil
}

func (m *Memoria) ApagarHistorico(ctx context.Context, dispositivoID string) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	for id, reg := range m.registros {
		if reg.dispositivo == dispositivoID && !reg.rascunho {
			delete(m.registros, id)
		}
	}
	return nil
}

func (r *registroMemoria) paraRascunho() Rascunho {
	return Rascunho{
		ID:           r.id,
		Etapa:        r.etapa,
		Respostas:    r.respostas,
		CriadoEm:     r.criadoEm,
		AtualizadoEm: r.atualizado,
	}
}

func (r *registroMemoria) paraConcluida() Concluida {
	c := Concluida{
		ID:          r.id,
		Respostas:   r.respostas,
		ConcluidaEm: r.concluidaEm,
	}
	if r.resultado != nil {
		c.Resultado = *r.resultado
	}
	return c
}
