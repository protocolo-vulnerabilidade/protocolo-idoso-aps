package httpapi

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"piaps/backend/internal/armazenamento"
)

const dispA = "11111111-1111-4111-8111-111111111111"
const dispB = "22222222-2222-4222-8222-222222222222"

func novoServidorTeste() http.Handler {
	return NovoHandler(armazenamento.NovoMemoria(), []string{"http://localhost:5173"})
}

func requisicao(t *testing.T, h http.Handler, metodo, caminho, dispositivo string, corpo any) *httptest.ResponseRecorder {
	t.Helper()
	var buf *bytes.Buffer
	if corpo != nil {
		b, err := json.Marshal(corpo)
		if err != nil {
			t.Fatal(err)
		}
		buf = bytes.NewBuffer(b)
	} else {
		buf = bytes.NewBuffer(nil)
	}
	req := httptest.NewRequest(metodo, caminho, buf)
	if dispositivo != "" {
		req.Header.Set("X-Dispositivo-Id", dispositivo)
	}
	if corpo != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, req)
	return rec
}

func decodificar(t *testing.T, rec *httptest.ResponseRecorder, v any) {
	t.Helper()
	if err := json.Unmarshal(rec.Body.Bytes(), v); err != nil {
		t.Fatalf("resposta não é JSON válido (%d): %s — %v", rec.Code, rec.Body.String(), err)
	}
}

func respostasCompletas() map[string]any {
	nao := []any{false, false, false}
	nao5 := []any{false, false, false, false, false}
	return map[string]any{
		"idade":         80,
		"autopercepcao": "boa",
		"aivd":          []any{false, false, false},
		"abvd":          false,
		"cognicao":      nao,
		"humor":         []any{false, false},
		"alcance":       []any{false, false},
		"aerobica":      false,
		"marcha":        nao,
		"comunicacao":   []any{false, false},
		"comorbidades":  false,
		"easfi":         []any{1, 1, 1, 1, 1},
		"sentinelas":    nao5,
	}
}

func criarRascunho(t *testing.T, h http.Handler, disp string) string {
	t.Helper()
	rec := requisicao(t, h, http.MethodPost, "/api/rascunhos", disp, nil)
	if rec.Code != http.StatusCreated {
		t.Fatalf("criar rascunho = %d: %s", rec.Code, rec.Body.String())
	}
	var out struct {
		ID string `json:"id"`
	}
	decodificar(t, rec, &out)
	if out.ID == "" {
		t.Fatal("id vazio")
	}
	return out.ID
}

func TestSaude(t *testing.T) {
	h := novoServidorTeste()
	rec := requisicao(t, h, http.MethodGet, "/api/saude", "", nil)
	if rec.Code != http.StatusOK {
		t.Fatalf("saude = %d", rec.Code)
	}
}

func TestDispositivoInvalido(t *testing.T) {
	h := novoServidorTeste()
	for _, disp := range []string{"", "abc", "não-uuid"} {
		rec := requisicao(t, h, http.MethodGet, "/api/rascunhos", disp, nil)
		if rec.Code != http.StatusBadRequest {
			t.Fatalf("dispositivo %q: esperava 400, obteve %d", disp, rec.Code)
		}
	}
}

func TestFluxoCompleto(t *testing.T) {
	h := novoServidorTeste()
	id := criarRascunho(t, h, dispA)

	rec := requisicao(t, h, http.MethodPut, "/api/rascunhos/"+id, dispA, map[string]any{
		"etapa":     3,
		"respostas": respostasCompletas(),
	})
	if rec.Code != http.StatusOK {
		t.Fatalf("salvar = %d: %s", rec.Code, rec.Body.String())
	}

	rec = requisicao(t, h, http.MethodGet, "/api/rascunhos/"+id, dispA, nil)
	if rec.Code != http.StatusOK {
		t.Fatalf("obter = %d", rec.Code)
	}
	var rasc struct {
		Etapa int `json:"etapa"`
	}
	decodificar(t, rec, &rasc)
	if rasc.Etapa != 3 {
		t.Errorf("etapa = %d, queria 3", rasc.Etapa)
	}

	rec = requisicao(t, h, http.MethodPost, "/api/rascunhos/"+id+"/concluir", dispA,
		map[string]any{"respostas": respostasCompletas()})
	if rec.Code != http.StatusOK {
		t.Fatalf("concluir = %d: %s", rec.Code, rec.Body.String())
	}
	var concl map[string]any
	decodificar(t, rec, &concl)
	resultado := concl["resultado"].(map[string]any)
	if resultado["ivcf"].(float64) != 1 {
		t.Errorf("ivcf = %v, queria 1 (idade 80)", resultado["ivcf"])
	}
	if resultado["easfi"].(float64) != 5 {
		t.Errorf("easfi = %v, queria 5", resultado["easfi"])
	}
	if resultado["nivel"].(float64) != 1 {
		t.Errorf("nivel = %v, queria 1", resultado["nivel"])
	}
	resumo, ok := concl["resumo"].(string)
	if !ok || resumo == "" {
		t.Errorf("resumo ausente na conclusão")
	}

	rec = requisicao(t, h, http.MethodGet, "/api/rascunhos/"+id, dispA, nil)
	if rec.Code != http.StatusNotFound {
		t.Fatalf("rascunho deveria ter sumido após concluir, obteve %d", rec.Code)
	}

	rec = requisicao(t, h, http.MethodGet, "/api/historico", dispA, nil)
	if rec.Code != http.StatusOK {
		t.Fatalf("listar histórico = %d", rec.Code)
	}
	var lista struct {
		Historico []map[string]any `json:"historico"`
	}
	decodificar(t, rec, &lista)
	if len(lista.Historico) != 1 {
		t.Fatalf("histórico com %d itens, queria 1", len(lista.Historico))
	}

	rec = requisicao(t, h, http.MethodPost, "/api/historico/"+id+"/reabrir", dispA, nil)
	if rec.Code != http.StatusOK {
		t.Fatalf("reabrir = %d: %s", rec.Code, rec.Body.String())
	}
	var reaberto struct {
		ID    string `json:"id"`
		Etapa int    `json:"etapa"`
	}
	decodificar(t, rec, &reaberto)
	if reaberto.ID != id || reaberto.Etapa != 0 {
		t.Errorf("reaberto = %+v", reaberto)
	}

	rec = requisicao(t, h, http.MethodGet, "/api/historico", dispA, nil)
	decodificar(t, rec, &lista)
	if len(lista.Historico) != 0 {
		t.Errorf("histórico deveria estar vazio após reabrir, tem %d", len(lista.Historico))
	}

	rec = requisicao(t, h, http.MethodPost, "/api/rascunhos/"+id+"/concluir", dispA,
		map[string]any{"respostas": respostasCompletas()})
	if rec.Code != http.StatusOK {
		t.Fatalf("reconcluir = %d: %s", rec.Code, rec.Body.String())
	}

	rec = requisicao(t, h, http.MethodDelete, "/api/historico/"+id, dispA, nil)
	if rec.Code != http.StatusNoContent {
		t.Fatalf("apagar concluída = %d", rec.Code)
	}
}

func TestConcluirIncompleto(t *testing.T) {
	h := novoServidorTeste()
	id := criarRascunho(t, h, dispA)
	respostas := respostasCompletas()
	respostas["idade"] = nil
	respostas["easfi"] = []any{1, nil, 1, 1, 1}
	rec := requisicao(t, h, http.MethodPost, "/api/rascunhos/"+id+"/concluir", dispA,
		map[string]any{"respostas": respostas})
	if rec.Code != http.StatusUnprocessableEntity {
		t.Fatalf("esperava 422, obteve %d: %s", rec.Code, rec.Body.String())
	}
	var corpo struct {
		Detalhes []struct {
			Campo string `json:"campo"`
		} `json:"detalhes"`
	}
	decodificar(t, rec, &corpo)
	campos := map[string]bool{}
	for _, d := range corpo.Detalhes {
		campos[d.Campo] = true
	}
	if !campos["idade"] || !campos["easfi[1]"] {
		t.Errorf("esperava erros em idade e easfi[1], obteve %v", corpo.Detalhes)
	}
}

func TestConcluirDuplicado(t *testing.T) {
	h := novoServidorTeste()
	id := criarRascunho(t, h, dispA)
	corpo := map[string]any{"respostas": respostasCompletas()}
	if rec := requisicao(t, h, http.MethodPost, "/api/rascunhos/"+id+"/concluir", dispA, corpo); rec.Code != http.StatusOK {
		t.Fatalf("primeira conclusão = %d", rec.Code)
	}
	rec := requisicao(t, h, http.MethodPost, "/api/rascunhos/"+id+"/concluir", dispA, corpo)
	if rec.Code != http.StatusConflict {
		t.Fatalf("segunda conclusão: esperava 409, obteve %d", rec.Code)
	}
}

func TestIsolamentoEntreDispositivos(t *testing.T) {
	h := novoServidorTeste()
	id := criarRascunho(t, h, dispA)
	rec := requisicao(t, h, http.MethodGet, "/api/rascunhos/"+id, dispB, nil)
	if rec.Code != http.StatusNotFound {
		t.Fatalf("dispositivo B não deveria acessar rascunho de A: %d", rec.Code)
	}
	rec = requisicao(t, h, http.MethodDelete, "/api/rascunhos/"+id, dispB, nil)
	if rec.Code != http.StatusNotFound {
		t.Fatalf("dispositivo B não deveria apagar rascunho de A: %d", rec.Code)
	}
}

func TestApagarHistoricoTotal(t *testing.T) {
	h := novoServidorTeste()
	for i := 0; i < 2; i++ {
		id := criarRascunho(t, h, dispA)
		rec := requisicao(t, h, http.MethodPost, "/api/rascunhos/"+id+"/concluir", dispA,
			map[string]any{"respostas": respostasCompletas()})
		if rec.Code != http.StatusOK {
			t.Fatalf("concluir %d = %d: %s", i, rec.Code, rec.Body.String())
		}
	}
	rascunhoRestante := criarRascunho(t, h, dispA)

	rec := requisicao(t, h, http.MethodDelete, "/api/historico", dispA, nil)
	if rec.Code != http.StatusNoContent {
		t.Fatalf("apagar tudo = %d", rec.Code)
	}

	var lista struct {
		Historico []map[string]any `json:"historico"`
	}
	rec = requisicao(t, h, http.MethodGet, "/api/historico", dispA, nil)
	decodificar(t, rec, &lista)
	if len(lista.Historico) != 0 {
		t.Errorf("histórico deveria estar vazio, tem %d", len(lista.Historico))
	}
	rec = requisicao(t, h, http.MethodGet, "/api/rascunhos/"+rascunhoRestante, dispA, nil)
	if rec.Code != http.StatusOK {
		t.Errorf("apagar histórico não deveria afetar rascunhos: %d", rec.Code)
	}
}

func TestCamposDesconhecidosRejeitados(t *testing.T) {
	h := novoServidorTeste()
	id := criarRascunho(t, h, dispA)
	respostas := respostasCompletas()
	respostas["nomeDoIdoso"] = "João"
	rec := requisicao(t, h, http.MethodPut, "/api/rascunhos/"+id, dispA,
		map[string]any{"etapa": 0, "respostas": respostas})
	if rec.Code != http.StatusBadRequest {
		t.Fatalf("campo desconhecido deveria gerar 400, obteve %d: %s", rec.Code, rec.Body.String())
	}
}

func TestCORSPreflight(t *testing.T) {
	h := novoServidorTeste()
	req := httptest.NewRequest(http.MethodOptions, "/api/rascunhos", nil)
	req.Header.Set("Origin", "http://localhost:5173")
	req.Header.Set("Access-Control-Request-Method", "POST")
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, req)
	if rec.Code != http.StatusNoContent {
		t.Fatalf("preflight = %d", rec.Code)
	}
	if rec.Header().Get("Access-Control-Allow-Origin") != "http://localhost:5173" {
		t.Fatalf("origem não liberada: %q", rec.Header().Get("Access-Control-Allow-Origin"))
	}

	req = httptest.NewRequest(http.MethodGet, "/api/rascunhos", nil)
	req.Header.Set("Origin", "http://evil.example")
	rec = httptest.NewRecorder()
	h.ServeHTTP(rec, req)
	if rec.Header().Get("Access-Control-Allow-Origin") != "" {
		t.Fatalf("origem não permitida foi liberada")
	}
	if rec.Code != http.StatusBadRequest {
		t.Fatalf("esperava 400 sem dispositivo, obteve %d", rec.Code)
	}
}

func TestCamposJSONCamelCase(t *testing.T) {
	h := novoServidorTeste()
	id := criarRascunho(t, h, dispA)
	rec := requisicao(t, h, http.MethodPost, "/api/rascunhos/"+id+"/concluir", dispA,
		map[string]any{"respostas": respostasCompletas()})
	if rec.Code != http.StatusOK {
		t.Fatalf("concluir = %d: %s", rec.Code, rec.Body.String())
	}
	var bruto map[string]json.RawMessage
	decodificar(t, rec, &bruto)
	var resultado map[string]json.RawMessage
	if err := json.Unmarshal(bruto["resultado"], &resultado); err != nil {
		t.Fatal(err)
	}
	for _, campo := range []string{"ivcf", "classeIVCF", "easfi", "classeEASFI", "nivelMatriz", "nivel", "sentinelasPositivas"} {
		if _, ok := resultado[campo]; !ok {
			t.Errorf("resultado sem o campo %q", campo)
		}
	}
	var respostas map[string]json.RawMessage
	if err := json.Unmarshal(bruto["respostas"], &respostas); err != nil {
		t.Fatal(err)
	}
	for _, campo := range []string{"idade", "autopercepcao", "aivd", "abvd", "cognicao", "humor", "alcance", "aerobica", "marcha", "comunicacao", "comorbidades", "easfi", "sentinelas"} {
		if _, ok := respostas[campo]; !ok {
			t.Errorf("respostas sem o campo %q", campo)
		}
	}
}
