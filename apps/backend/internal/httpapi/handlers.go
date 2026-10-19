package httpapi

import (
	"errors"
	"io"
	"net/http"
	"time"

	"piaps/backend/internal/armazenamento"
	"piaps/backend/internal/dominio"
)

type erroCampo = dominio.ErroValidacao

type rascunhoJSON struct {
	ID           string            `json:"id"`
	CriadoEm     time.Time         `json:"criadoEm"`
	AtualizadoEm time.Time         `json:"atualizadoEm"`
	Etapa        int               `json:"etapa"`
	Respostas    dominio.Respostas `json:"respostas"`
}

type concluidaJSON struct {
	ID          string            `json:"id"`
	ConcluidaEm time.Time         `json:"concluidaEm"`
	Respostas   dominio.Respostas `json:"respostas"`
	Resultado   dominio.Resultado `json:"resultado"`
	Resumo      *string           `json:"resumo,omitempty"`
}

func paraRascunhoJSON(r armazenamento.Rascunho) rascunhoJSON {
	return rascunhoJSON{
		ID:           r.ID,
		CriadoEm:     r.CriadoEm.UTC(),
		AtualizadoEm: r.AtualizadoEm.UTC(),
		Etapa:        r.Etapa,
		Respostas:    r.Respostas,
	}
}

func paraConcluidaJSON(c armazenamento.Concluida) concluidaJSON {
	return concluidaJSON{
		ID:          c.ID,
		ConcluidaEm: c.ConcluidaEm.UTC(),
		Respostas:   c.Respostas,
		Resultado:   c.Resultado,
	}
}

func (s *Servidor) mapearErroArmazenamento(w http.ResponseWriter, err error) bool {
	switch {
	case err == nil:
		return false
	case errors.Is(err, armazenamento.ErrNaoEncontrado):
		s.erro(w, http.StatusNotFound, "registro não encontrado")
	case errors.Is(err, armazenamento.ErrConflito):
		s.erro(w, http.StatusConflict, "registro em estado conflitante")
	default:
		s.erroInterno(w, err)
	}
	return true
}

func (s *Servidor) listarRascunhos(w http.ResponseWriter, r *http.Request, dispositivo string) {
	lista, err := s.arm.ListarRascunhos(r.Context(), dispositivo)
	if s.mapearErroArmazenamento(w, err) {
		return
	}
	out := make([]rascunhoJSON, 0, len(lista))
	for _, item := range lista {
		out = append(out, paraRascunhoJSON(item))
	}
	s.json(w, http.StatusOK, map[string]any{"rascunhos": out})
}

func (s *Servidor) criarRascunho(w http.ResponseWriter, r *http.Request, dispositivo string) {
	respostas := dominio.RespostasVazias()
	var corpo struct {
		Respostas *dominio.Respostas `json:"respostas"`
	}
	if err := decodificarCorpo(r, &corpo); err != nil && !errors.Is(err, io.EOF) {
		s.erro(w, http.StatusBadRequest, "corpo inválido")
		return
	}
	if corpo.Respostas != nil {
		respostas = *corpo.Respostas
	}
	if erros := dominio.ValidarEstrutura(&respostas); len(erros) > 0 {
		s.erroValidacao(w, erros)
		return
	}
	rascunho, err := s.arm.CriarRascunho(r.Context(), dispositivo, respostas)
	if s.mapearErroArmazenamento(w, err) {
		return
	}
	s.json(w, http.StatusCreated, paraRascunhoJSON(rascunho))
}

func (s *Servidor) obterRascunho(w http.ResponseWriter, r *http.Request, dispositivo string) {
	rascunho, err := s.arm.ObterRascunho(r.Context(), dispositivo, r.PathValue("id"))
	if s.mapearErroArmazenamento(w, err) {
		return
	}
	s.json(w, http.StatusOK, paraRascunhoJSON(rascunho))
}

func (s *Servidor) salvarRascunho(w http.ResponseWriter, r *http.Request, dispositivo string) {
	var corpo struct {
		Etapa     int                `json:"etapa"`
		Respostas *dominio.Respostas `json:"respostas"`
	}
	if err := decodificarCorpo(r, &corpo); err != nil {
		s.erro(w, http.StatusBadRequest, "corpo inválido")
		return
	}
	if corpo.Etapa < 0 || corpo.Etapa > 6 {
		s.erro(w, http.StatusBadRequest, "etapa deve estar entre 0 e 6")
		return
	}
	if corpo.Respostas == nil {
		s.erro(w, http.StatusBadRequest, "campo \"respostas\" obrigatório")
		return
	}
	if erros := dominio.ValidarEstrutura(corpo.Respostas); len(erros) > 0 {
		s.erroValidacao(w, erros)
		return
	}
	rascunho, err := s.arm.SalvarRascunho(r.Context(), dispositivo, r.PathValue("id"), corpo.Etapa, *corpo.Respostas)
	if s.mapearErroArmazenamento(w, err) {
		return
	}
	s.json(w, http.StatusOK, paraRascunhoJSON(rascunho))
}

func (s *Servidor) apagarRascunho(w http.ResponseWriter, r *http.Request, dispositivo string) {
	if s.mapearErroArmazenamento(w, s.arm.ApagarRascunho(r.Context(), dispositivo, r.PathValue("id"))) {
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (s *Servidor) concluir(w http.ResponseWriter, r *http.Request, dispositivo string) {
	var corpo struct {
		Respostas *dominio.Respostas `json:"respostas"`
	}
	if err := decodificarCorpo(r, &corpo); err != nil {
		s.erro(w, http.StatusBadRequest, "corpo inválido")
		return
	}
	if corpo.Respostas == nil {
		s.erro(w, http.StatusBadRequest, "campo \"respostas\" obrigatório")
		return
	}
	if erros := dominio.ValidarCompleta(corpo.Respostas); len(erros) > 0 {
		s.erroValidacao(w, erros)
		return
	}
	resultado := dominio.Calcular(corpo.Respostas)
	c, err := s.arm.Concluir(r.Context(), dispositivo, r.PathValue("id"), *corpo.Respostas, resultado)
	if s.mapearErroArmazenamento(w, err) {
		return
	}
	out := paraConcluidaJSON(c)
	resumo := dominio.GerarResumo(c.Resultado, c.ConcluidaEm)
	out.Resumo = &resumo
	s.json(w, http.StatusOK, out)
}

func (s *Servidor) listarHistorico(w http.ResponseWriter, r *http.Request, dispositivo string) {
	lista, err := s.arm.ListarHistorico(r.Context(), dispositivo)
	if s.mapearErroArmazenamento(w, err) {
		return
	}
	out := make([]concluidaJSON, 0, len(lista))
	for _, item := range lista {
		out = append(out, paraConcluidaJSON(item))
	}
	s.json(w, http.StatusOK, map[string]any{"historico": out})
}

func (s *Servidor) obterConcluida(w http.ResponseWriter, r *http.Request, dispositivo string) {
	c, err := s.arm.ObterConcluida(r.Context(), dispositivo, r.PathValue("id"))
	if s.mapearErroArmazenamento(w, err) {
		return
	}
	out := paraConcluidaJSON(c)
	resumo := dominio.GerarResumo(c.Resultado, c.ConcluidaEm)
	out.Resumo = &resumo
	s.json(w, http.StatusOK, out)
}

func (s *Servidor) apagarConcluida(w http.ResponseWriter, r *http.Request, dispositivo string) {
	if s.mapearErroArmazenamento(w, s.arm.ApagarConcluida(r.Context(), dispositivo, r.PathValue("id"))) {
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (s *Servidor) apagarHistorico(w http.ResponseWriter, r *http.Request, dispositivo string) {
	if s.mapearErroArmazenamento(w, s.arm.ApagarHistorico(r.Context(), dispositivo)) {
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (s *Servidor) reabrir(w http.ResponseWriter, r *http.Request, dispositivo string) {
	rascunho, err := s.arm.Reabrir(r.Context(), dispositivo, r.PathValue("id"))
	if s.mapearErroArmazenamento(w, err) {
		return
	}
	s.json(w, http.StatusOK, paraRascunhoJSON(rascunho))
}
