package httpapi

import (
	"crypto/subtle"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/google/uuid"

	"piaps/backend/internal/armazenamento"
)

const (
	maxErrosSenha   = 5
	bloqueioDuracao = time.Minute
)

type tentativasSenha struct {
	erros        int
	bloqueadoAte time.Time
}

type Servidor struct {
	arm           armazenamento.Armazenamento
	origens       map[string]bool
	senhaExclusao string
	mu            sync.Mutex
	tentativas    map[string]*tentativasSenha
}

func NovoHandler(arm armazenamento.Armazenamento, origensPermitidas []string, senhaExclusao string) http.Handler {
	origens := make(map[string]bool, len(origensPermitidas))
	for _, o := range origensPermitidas {
		if o != "" {
			origens[o] = true
		}
	}
	s := &Servidor{
		arm:           arm,
		origens:       origens,
		senhaExclusao: senhaExclusao,
		tentativas:    make(map[string]*tentativasSenha),
	}

	mux := http.NewServeMux()
	mux.HandleFunc("GET /api/saude", s.saude)
	mux.HandleFunc("GET /api/rascunhos", s.exigirDispositivo(s.listarRascunhos))
	mux.HandleFunc("POST /api/rascunhos", s.exigirDispositivo(s.criarRascunho))
	mux.HandleFunc("GET /api/rascunhos/{id}", s.exigirDispositivo(s.obterRascunho))
	mux.HandleFunc("PUT /api/rascunhos/{id}", s.exigirDispositivo(s.salvarRascunho))
	mux.HandleFunc("DELETE /api/rascunhos/{id}", s.exigirDispositivo(s.apagarRascunho))
	mux.HandleFunc("POST /api/rascunhos/{id}/concluir", s.exigirDispositivo(s.concluir))
	mux.HandleFunc("GET /api/historico", s.exigirDispositivo(s.listarHistorico))
	mux.HandleFunc("DELETE /api/historico", s.exigirDispositivo(s.exigirSenha(s.apagarHistorico)))
	mux.HandleFunc("GET /api/historico/{id}", s.exigirDispositivo(s.obterConcluida))
	mux.HandleFunc("DELETE /api/historico/{id}", s.exigirDispositivo(s.exigirSenha(s.apagarConcluida)))
	mux.HandleFunc("POST /api/historico/{id}/reabrir", s.exigirDispositivo(s.reabrir))

	return s.cors(s.registrarLog(mux))
}

func (s *Servidor) saude(w http.ResponseWriter, r *http.Request) {
	s.json(w, http.StatusOK, map[string]string{"status": "ok"})
}

func (s *Servidor) exigirDispositivo(prox func(http.ResponseWriter, *http.Request, string)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := r.Header.Get("X-Dispositivo-Id")
		if _, err := uuid.Parse(id); err != nil {
			s.erro(w, http.StatusBadRequest, "cabeçalho X-Dispositivo-Id ausente ou inválido (UUID)")
			return
		}
		if err := s.arm.RegistrarDispositivo(r.Context(), id); err != nil {
			s.erroInterno(w, err)
			return
		}
		prox(w, r, id)
	}
}

func (s *Servidor) exigirSenha(prox func(http.ResponseWriter, *http.Request, string)) func(http.ResponseWriter, *http.Request, string) {
	return func(w http.ResponseWriter, r *http.Request, dispositivo string) {
		s.mu.Lock()
		reg := s.tentativas[dispositivo]
		agora := time.Now()
		if reg != nil && reg.bloqueadoAte.After(agora) {
			s.mu.Unlock()
			s.erro(w, http.StatusTooManyRequests, "muitas tentativas incorretas; aguarde um minuto e tente novamente")
			return
		}
		senha := r.Header.Get("X-Senha-Exclusao")
		if senha != "" && subtle.ConstantTimeCompare([]byte(senha), []byte(s.senhaExclusao)) == 1 {
			delete(s.tentativas, dispositivo)
			s.mu.Unlock()
			prox(w, r, dispositivo)
			return
		}
		if reg == nil {
			reg = &tentativasSenha{}
			s.tentativas[dispositivo] = reg
		}
		reg.erros++
		if reg.erros >= maxErrosSenha {
			reg.bloqueadoAte = agora.Add(bloqueioDuracao)
			reg.erros = 0
		}
		s.mu.Unlock()
		s.erro(w, http.StatusUnauthorized, "senha inválida")
	}
}

type registrador struct {
	http.ResponseWriter
	status int
}

func (reg *registrador) WriteHeader(status int) {
	reg.status = status
	reg.ResponseWriter.WriteHeader(status)
}

func (s *Servidor) registrarLog(prox http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		inicio := time.Now()
		reg := &registrador{ResponseWriter: w, status: http.StatusOK}
		prox.ServeHTTP(reg, r)
		log.Printf("%s %s %d %s", r.Method, r.URL.Path, reg.status, time.Since(inicio).Round(time.Millisecond))
	})
}

func (s *Servidor) cors(prox http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Vary", "Origin")
		origem := r.Header.Get("Origin")
		if origem != "" && s.origens[origem] {
			w.Header().Set("Access-Control-Allow-Origin", origem)
			if r.Method == http.MethodOptions {
				w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
				w.Header().Set("Access-Control-Allow-Headers", "Content-Type, X-Dispositivo-Id, X-Senha-Exclusao")
				w.WriteHeader(http.StatusNoContent)
				return
			}
		}
		prox.ServeHTTP(w, r)
	})
}

func (s *Servidor) json(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(v); err != nil {
		log.Printf("escrever resposta: %v", err)
	}
}

func (s *Servidor) erro(w http.ResponseWriter, status int, msg string) {
	s.json(w, status, map[string]string{"erro": msg})
}

func (s *Servidor) erroValidacao(w http.ResponseWriter, erros []erroCampo) {
	s.json(w, http.StatusUnprocessableEntity, map[string]any{
		"erro":     "validação das respostas falhou",
		"detalhes": erros,
	})
}

func (s *Servidor) erroInterno(w http.ResponseWriter, err error) {
	log.Printf("erro interno: %v", err)
	s.erro(w, http.StatusInternalServerError, "erro interno")
}

func decodificarCorpo(r *http.Request, v any) error {
	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()
	err := dec.Decode(v)
	if err == io.EOF {
		return io.EOF
	}
	if err != nil {
		return err
	}
	return nil
}
