package dominio

import (
	"encoding/json"
	"strings"
	"testing"
	"time"
)

func boolPtr(v bool) *bool    { return &v }
func intPtr(v int) *int       { return &v }
func strPtr(v string) *string { return &v }

func respostasMaximas() Respostas {
	r := RespostasVazias()
	r.Idade = intPtr(90)
	r.Autopercepcao = strPtr("ruim")
	r.Aivd = []*bool{boolPtr(true), boolPtr(true), boolPtr(true)}
	r.Abvd = boolPtr(true)
	r.Cognicao = []*bool{boolPtr(true), boolPtr(true), boolPtr(true)}
	r.Humor = []*bool{boolPtr(true), boolPtr(true)}
	r.Alcance = []*bool{boolPtr(true), boolPtr(true)}
	r.Aerobica = boolPtr(true)
	r.Marcha = []*bool{boolPtr(true), boolPtr(true), boolPtr(true)}
	r.Comunicacao = []*bool{boolPtr(true), boolPtr(true)}
	r.Comorbidades = boolPtr(true)
	return r
}

func TestPontosIdade(t *testing.T) {
	casos := []struct {
		idade int
		quer  int
	}{
		{60, 0}, {74, 0}, {75, 1}, {84, 1}, {85, 3}, {90, 3}, {130, 3},
	}
	for _, c := range casos {
		if ob := PontosIdade(c.idade); ob != c.quer {
			t.Errorf("PontosIdade(%d) = %d, queria %d", c.idade, ob, c.quer)
		}
	}
}

func TestDominiosIVCF(t *testing.T) {
	t.Run("AIVD tem teto de 4 mesmo com três Sim", func(t *testing.T) {
		r := RespostasVazias()
		r.Aivd = []*bool{boolPtr(true), boolPtr(true), boolPtr(true)}
		if ob := PontuarDominiosIVCF(&r).Aivd; ob != 4 {
			t.Errorf("Aivd = %d, queria 4", ob)
		}
	})
	t.Run("ABVD vale 6", func(t *testing.T) {
		r := RespostasVazias()
		r.Abvd = boolPtr(true)
		if ob := PontuarDominiosIVCF(&r).Abvd; ob != 6 {
			t.Errorf("Abvd = %d, queria 6", ob)
		}
	})
	t.Run("cognição pesos 1,1,2", func(t *testing.T) {
		r := RespostasVazias()
		r.Cognicao = []*bool{boolPtr(true), boolPtr(true), boolPtr(true)}
		if ob := PontuarDominiosIVCF(&r).Cognicao; ob != 4 {
			t.Errorf("Cognicao = %d, queria 4", ob)
		}
	})
	t.Run("humor pesos 2,2", func(t *testing.T) {
		r := RespostasVazias()
		r.Humor = []*bool{boolPtr(true), boolPtr(true)}
		if ob := PontuarDominiosIVCF(&r).Humor; ob != 4 {
			t.Errorf("Humor = %d, queria 4", ob)
		}
	})
	t.Run("marcha pesos 2,2,2", func(t *testing.T) {
		r := RespostasVazias()
		r.Marcha = []*bool{boolPtr(true), boolPtr(true), boolPtr(true)}
		if ob := PontuarDominiosIVCF(&r).Marcha; ob != 6 {
			t.Errorf("Marcha = %d, queria 6", ob)
		}
	})
	t.Run("comorbidades pergunta única de 4", func(t *testing.T) {
		r := RespostasVazias()
		r.Comorbidades = boolPtr(true)
		if ob := PontuarDominiosIVCF(&r).Comorbidades; ob != 4 {
			t.Errorf("Comorbidades = %d, queria 4", ob)
		}
	})
	t.Run("autopercepcao ruim vale 1", func(t *testing.T) {
		r := RespostasVazias()
		r.Autopercepcao = strPtr("ruim")
		if ob := PontuarDominiosIVCF(&r).Autopercepcao; ob != 1 {
			t.Errorf("Autopercepcao = %d, queria 1", ob)
		}
	})
}

func TestEscoreMaximoExatamente40(t *testing.T) {
	r := respostasMaximas()
	if ob := EscoreIVCF(&r); ob != 40 {
		t.Errorf("escore máximo = %d, queria 40", ob)
	}
}

func TestEscoreZero(t *testing.T) {
	r := RespostasVazias()
	r.Idade = intPtr(70)
	if ob := EscoreIVCF(&r); ob != 0 {
		t.Errorf("escore mínimo = %d, queria 0", ob)
	}
}

func TestClassificarIVCFFronteiras(t *testing.T) {
	casos := []struct {
		escore int
		quer   ClasseIVCF
	}{
		{0, ClasseIVCFBaixo}, {6, ClasseIVCFBaixo},
		{7, ClasseIVCFModerado}, {14, ClasseIVCFModerado},
		{15, ClasseIVCFAlto}, {40, ClasseIVCFAlto},
	}
	for _, c := range casos {
		if ob := ClassificarIVCF(c.escore); ob != c.quer {
			t.Errorf("ClassificarIVCF(%d) = %q, queria %q", c.escore, ob, c.quer)
		}
	}
}

func TestClassificarEASFIFronteiras(t *testing.T) {
	casos := []struct {
		escore int
		quer   ClasseEASFI
	}{
		{5, ClasseEASFIAdequada}, {9, ClasseEASFIAdequada},
		{10, ClasseEASFIRisco}, {14, ClasseEASFIRisco},
		{15, ClasseEASFIProblema}, {25, ClasseEASFIProblema},
	}
	for _, c := range casos {
		if ob := ClassificarEASFI(c.escore); ob != c.quer {
			t.Errorf("ClassificarEASFI(%d) = %q, queria %q", c.escore, ob, c.quer)
		}
	}
}

func TestMatrizCompleta(t *testing.T) {
	casos := []struct {
		ivcf  ClasseIVCF
		easfi ClasseEASFI
		quer  Nivel
	}{
		{ClasseIVCFBaixo, ClasseEASFIAdequada, 1},
		{ClasseIVCFBaixo, ClasseEASFIRisco, 2},
		{ClasseIVCFBaixo, ClasseEASFIProblema, 3},
		{ClasseIVCFModerado, ClasseEASFIAdequada, 2},
		{ClasseIVCFModerado, ClasseEASFIRisco, 3},
		{ClasseIVCFModerado, ClasseEASFIProblema, 4},
		{ClasseIVCFAlto, ClasseEASFIAdequada, 3},
		{ClasseIVCFAlto, ClasseEASFIRisco, 4},
		{ClasseIVCFAlto, ClasseEASFIProblema, 4},
	}
	for _, c := range casos {
		if ob := NivelMatriz(c.ivcf, c.easfi); ob != c.quer {
			t.Errorf("NivelMatriz(%q,%q) = %d, queria %d", c.ivcf, c.easfi, ob, c.quer)
		}
	}
}

func sentinelasCom(indice int) []*bool {
	s := make([]*bool, 5)
	s[indice] = boolPtr(true)
	return s
}

func TestAplicarSentinelas(t *testing.T) {
	nenhuma := make([]*bool, 5)
	casos := []struct {
		nome  string
		nivel Nivel
		sent  []*bool
		quer  Nivel
	}{
		{"nenhuma mantém", 1, nenhuma, 1},
		{"leve (índice 2) eleva 1→3", 1, sentinelasCom(2), 3},
		{"leve (índice 3) eleva 1→3", 1, sentinelasCom(3), 3},
		{"leve (índice 4) eleva 2→3", 2, sentinelasCom(4), 3},
		{"leve não reduz 4", 4, sentinelasCom(2), 4},
		{"grave (índice 0) leva a 4", 1, sentinelasCom(0), 4},
		{"grave (índice 1) leva a 4", 1, sentinelasCom(1), 4},
		{"grave eleva 3→4", 3, sentinelasCom(0), 4},
	}
	for _, c := range casos {
		if ob := AplicarSentinelas(c.nivel, c.sent); ob != c.quer {
			t.Errorf("%s: AplicarSentinelas(%d) = %d, queria %d", c.nome, c.nivel, ob, c.quer)
		}
	}
}

func TestCalcularCompleto(t *testing.T) {
	r := RespostasVazias()
	r.Idade = intPtr(80)
	r.Easfi = []*int{intPtr(1), intPtr(1), intPtr(1), intPtr(1), intPtr(1)}
	res := Calcular(&r)
	if res.Ivcf != 1 || res.Easfi != 5 || res.Nivel != 1 || res.NivelMatriz != 1 {
		t.Errorf("Calcular = %+v, queria ivcf=1 easfi=5 nivel=1", res)
	}
}

func TestCalcularMaximoComSentinelaGrave(t *testing.T) {
	r := respostasMaximas()
	r.Easfi = []*int{intPtr(5), intPtr(5), intPtr(5), intPtr(5), intPtr(5)}
	r.Sentinelas = sentinelasCom(0)
	res := Calcular(&r)
	if res.Ivcf != 40 || res.ClasseIVCF != ClasseIVCFAlto || res.Easfi != 25 {
		t.Errorf("escores = %+v", res)
	}
	if res.NivelMatriz != 4 || res.Nivel != 4 {
		t.Errorf("níveis = matriz %d final %d, queria 4/4", res.NivelMatriz, res.Nivel)
	}
	if len(res.SentinelasPositivas) != 1 || res.SentinelasPositivas[0] != 0 {
		t.Errorf("sentinelasPositivas = %v", res.SentinelasPositivas)
	}
}

func TestCalcularReclassificacaoMatriz3Para4(t *testing.T) {
	r := RespostasVazias()
	r.Idade = intPtr(60)
	r.Easfi = []*int{intPtr(2), intPtr(2), intPtr(2), intPtr(2), intPtr(1)}
	res := Calcular(&r)
	if res.NivelMatriz != 1 {
		t.Fatalf("nivelMatriz = %d, queria 1", res.NivelMatriz)
	}
	r.Sentinelas = sentinelasCom(1)
	res = Calcular(&r)
	if res.Nivel != 4 {
		t.Errorf("sentinela grave deveria reclassificar para 4, obteve %d", res.Nivel)
	}
}

func TestValidarEstrutura(t *testing.T) {
	r := RespostasVazias()
	r.Idade = intPtr(59)
	r.Autopercepcao = strPtr("ótima")
	r.Aivd = []*bool{nil}
	r.Easfi = []*int{intPtr(1), intPtr(1), intPtr(1), intPtr(1), intPtr(1), intPtr(1)}
	erros := ValidarEstrutura(&r)
	campos := map[string]bool{}
	for _, e := range erros {
		campos[e.Campo] = true
	}
	for _, campo := range []string{"idade", "autopercepcao", "aivd", "easfi"} {
		if !campos[campo] {
			t.Errorf("esperava erro no campo %q; erros=%v", campo, erros)
		}
	}
	if len(erros) != 4 {
		t.Errorf("total de erros = %d, queria 4: %v", len(erros), erros)
	}

	r2 := RespostasVazias()
	r2.Easfi = []*int{intPtr(0), intPtr(6), nil, nil, nil}
	erros = ValidarEstrutura(&r2)
	if len(erros) != 2 || erros[0].Campo != "easfi[0]" || erros[1].Campo != "easfi[1]" {
		t.Errorf("valores EASFI fora de 1–5 deveriam gerar easfi[0]/easfi[1]: %v", erros)
	}
}

func TestValidarEstruturaVaziaOk(t *testing.T) {
	r := RespostasVazias()
	if erros := ValidarEstrutura(&r); len(erros) != 0 {
		t.Errorf("estrutura vazia não deveria ter erros: %v", erros)
	}
}

func TestValidarCompletaExigeTudo(t *testing.T) {
	r := RespostasVazias()
	erros := ValidarCompleta(&r)
	esperados := 1 + 1 + 3 + 1 + 3 + 2 + 2 + 1 + 3 + 2 + 1 + 5 + 5
	if len(erros) != esperados {
		t.Errorf("total = %d, queria %d: %v", len(erros), esperados, erros)
	}
	r = respostasMaximas()
	r.Easfi = []*int{intPtr(3), intPtr(3), intPtr(3), intPtr(3), intPtr(3)}
	r.Sentinelas = []*bool{boolPtr(false), boolPtr(false), boolPtr(false), boolPtr(false), boolPtr(false)}
	if erros := ValidarCompleta(&r); len(erros) != 0 {
		t.Errorf("respostas completas não deveriam ter erros: %v", erros)
	}
}

func TestGerarResumo(t *testing.T) {
	r := RespostasVazias()
	r.Idade = intPtr(80)
	r.Easfi = []*int{intPtr(2), intPtr(2), intPtr(2), intPtr(2), intPtr(2)}
	r.Sentinelas = sentinelasCom(3)
	res := Calcular(&r)
	resumo := GerarResumo(res, time.Date(2026, 10, 19, 14, 30, 45, 0, time.UTC))
	esperados := []string{
		"PROTOCOLO IDOSO APS — Estratificação de vulnerabilidade",
		"Data/hora: 19/10/2026 11:30:45",
		"IVCF-20: 1/40 — Baixo risco",
		"EASFI-BR: 10/25 — Risco social",
		"Nível pela matriz MIV-CFS: 2",
		"Marcadores sentinela presentes:",
		"  - Sobrecarga grave do cuidador com risco iminente de ruptura do cuidado",
		"NÍVEL DE VULNERABILIDADE: 3 — Alta",
		"Conduta: Discussão interdisciplinar, PTS obrigatório, visita domiciliar e articulação eMulti/CRAS.",
		"Reavaliação: A cada 3 meses",
	}
	for _, trecho := range esperados {
		if !strings.Contains(resumo, trecho) {
			t.Errorf("resumo deveria conter %q:\n%s", trecho, resumo)
		}
	}
	if strings.Count(resumo, "\n") != len(esperados)-1 {
		t.Errorf("resumo com número inesperado de linhas:\n%s", resumo)
	}
}

func TestJSONRespostasCamposExatos(t *testing.T) {
	corpo := `{"idade":null,"autopercepcao":null,"aivd":[null,null,null],"abvd":null,"cognicao":[true,null,false],"humor":[null,null],"alcance":[null,null],"aerobica":null,"marcha":[null,null,null],"comunicacao":[null,null],"comorbidades":null,"easfi":[1,null,3,4,5],"sentinelas":[null,null,null,null,null]}`
	var r Respostas
	if err := json.Unmarshal([]byte(corpo), &r); err != nil {
		t.Fatal(err)
	}
	if len(r.Aivd) != 3 || len(r.Cognicao) != 3 || len(r.Easfi) != 5 || len(r.Sentinelas) != 5 {
		t.Errorf("comprimentos inesperados: %+v", r)
	}
	if r.Cognicao[0] == nil || !*r.Cognicao[0] || r.Cognicao[2] == nil || *r.Cognicao[2] {
		t.Errorf("cognicao decodificada errado: %v", r.Cognicao)
	}
	if r.Easfi[1] != nil {
		t.Errorf("easfi[1] deveria ser nil: %v", r.Easfi)
	}
}
