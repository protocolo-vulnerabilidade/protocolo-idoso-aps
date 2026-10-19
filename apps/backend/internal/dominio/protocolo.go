package dominio

import "fmt"

type Respostas struct {
	Idade         *int    `json:"idade"`
	Autopercepcao *string `json:"autopercepcao"`
	Aivd          []*bool `json:"aivd"`
	Abvd          *bool   `json:"abvd"`
	Cognicao      []*bool `json:"cognicao"`
	Humor         []*bool `json:"humor"`
	Alcance       []*bool `json:"alcance"`
	Aerobica      *bool   `json:"aerobica"`
	Marcha        []*bool `json:"marcha"`
	Comunicacao   []*bool `json:"comunicacao"`
	Comorbidades  *bool   `json:"comorbidades"`
	Easfi         []*int  `json:"easfi"`
	Sentinelas    []*bool `json:"sentinelas"`
}

func RespostasVazias() Respostas {
	return Respostas{
		Aivd:        make([]*bool, 3),
		Cognicao:    make([]*bool, 3),
		Humor:       make([]*bool, 2),
		Alcance:     make([]*bool, 2),
		Marcha:      make([]*bool, 3),
		Comunicacao: make([]*bool, 2),
		Easfi:       make([]*int, 5),
		Sentinelas:  make([]*bool, 5),
	}
}

type ErroValidacao struct {
	Campo    string `json:"campo"`
	Mensagem string `json:"mensagem"`
}

func ValidarEstrutura(r *Respostas) []ErroValidacao {
	var erros []ErroValidacao
	if r.Idade != nil && (*r.Idade < 60 || *r.Idade > 130) {
		erros = append(erros, ErroValidacao{Campo: "idade", Mensagem: "deve estar entre 60 e 130 anos"})
	}
	if r.Autopercepcao != nil && *r.Autopercepcao != "boa" && *r.Autopercepcao != "ruim" {
		erros = append(erros, ErroValidacao{Campo: "autopercepcao", Mensagem: "deve ser \"boa\" ou \"ruim\""})
	}
	erros = append(erros, validarVetorBool("aivd", r.Aivd, 3)...)
	erros = append(erros, validarVetorBool("cognicao", r.Cognicao, 3)...)
	erros = append(erros, validarVetorBool("humor", r.Humor, 2)...)
	erros = append(erros, validarVetorBool("alcance", r.Alcance, 2)...)
	erros = append(erros, validarVetorBool("marcha", r.Marcha, 3)...)
	erros = append(erros, validarVetorBool("comunicacao", r.Comunicacao, 2)...)
	if len(r.Easfi) != 5 {
		erros = append(erros, ErroValidacao{Campo: "easfi", Mensagem: "esperados 5 domínios"})
	} else {
		for i, v := range r.Easfi {
			if v != nil && (*v < 1 || *v > 5) {
				erros = append(erros, ErroValidacao{
					Campo:    fmt.Sprintf("easfi[%d]", i),
					Mensagem: "deve estar entre 1 e 5",
				})
			}
		}
	}
	erros = append(erros, validarVetorBool("sentinelas", r.Sentinelas, 5)...)
	return erros
}

func ValidarCompleta(r *Respostas) []ErroValidacao {
	erros := ValidarEstrutura(r)
	if r.Idade == nil {
		erros = append(erros, ErroValidacao{Campo: "idade", Mensagem: "resposta obrigatória"})
	}
	if r.Autopercepcao == nil {
		erros = append(erros, ErroValidacao{Campo: "autopercepcao", Mensagem: "resposta obrigatória"})
	}
	erros = append(erros, exigirTodosBool("aivd", r.Aivd)...)
	erros = append(erros, exigirTodosBool("cognicao", r.Cognicao)...)
	erros = append(erros, exigirTodosBool("humor", r.Humor)...)
	erros = append(erros, exigirTodosBool("alcance", r.Alcance)...)
	erros = append(erros, exigirTodosBool("marcha", r.Marcha)...)
	erros = append(erros, exigirTodosBool("comunicacao", r.Comunicacao)...)
	erros = append(erros, exigirTodosBool("sentinelas", r.Sentinelas)...)
	if r.Abvd == nil {
		erros = append(erros, ErroValidacao{Campo: "abvd", Mensagem: "resposta obrigatória"})
	}
	if r.Aerobica == nil {
		erros = append(erros, ErroValidacao{Campo: "aerobica", Mensagem: "resposta obrigatória"})
	}
	if r.Comorbidades == nil {
		erros = append(erros, ErroValidacao{Campo: "comorbidades", Mensagem: "resposta obrigatória"})
	}
	if len(r.Easfi) == 5 {
		for i, v := range r.Easfi {
			if v == nil {
				erros = append(erros, ErroValidacao{
					Campo:    fmt.Sprintf("easfi[%d]", i),
					Mensagem: "resposta obrigatória",
				})
			}
		}
	}
	return erros
}

func validarVetorBool(campo string, arr []*bool, esperado int) []ErroValidacao {
	if len(arr) != esperado {
		return []ErroValidacao{{
			Campo:    campo,
			Mensagem: fmt.Sprintf("esperados %d itens", esperado),
		}}
	}
	return nil
}

func exigirTodosBool(campo string, arr []*bool) []ErroValidacao {
	if len(arr) == 0 {
		return nil
	}
	var erros []ErroValidacao
	for i, v := range arr {
		if v == nil {
			erros = append(erros, ErroValidacao{
				Campo:    fmt.Sprintf("%s[%d]", campo, i),
				Mensagem: "resposta obrigatória",
			})
		}
	}
	return erros
}

func sim(v *bool) bool { return v != nil && *v }

func somaBool(arr []*bool, pesos []int) int {
	total := 0
	for i, v := range arr {
		if i >= len(pesos) {
			break
		}
		if sim(v) {
			total += pesos[i]
		}
	}
	return total
}

func PontosIdade(idade int) int {
	if idade >= 85 {
		return 3
	}
	if idade >= 75 {
		return 1
	}
	return 0
}

type DominiosIVCF struct {
	Idade         int
	Autopercepcao int
	Aivd          int
	Abvd          int
	Cognicao      int
	Humor         int
	Alcance       int
	Aerobica      int
	Marcha        int
	Comunicacao   int
	Comorbidades  int
}

func PontuarDominiosIVCF(r *Respostas) DominiosIVCF {
	d := DominiosIVCF{}
	if r.Idade != nil {
		d.Idade = PontosIdade(*r.Idade)
	}
	if r.Autopercepcao != nil && *r.Autopercepcao == "ruim" {
		d.Autopercepcao = 1
	}
	d.Aivd = min(somaBool(r.Aivd, []int{4, 4, 4}), 4)
	if sim(r.Abvd) {
		d.Abvd = 6
	}
	d.Cognicao = somaBool(r.Cognicao, []int{1, 1, 2})
	d.Humor = somaBool(r.Humor, []int{2, 2})
	d.Alcance = somaBool(r.Alcance, []int{1, 1})
	if sim(r.Aerobica) {
		d.Aerobica = 2
	}
	d.Marcha = somaBool(r.Marcha, []int{2, 2, 2})
	d.Comunicacao = somaBool(r.Comunicacao, []int{2, 2})
	if sim(r.Comorbidades) {
		d.Comorbidades = 4
	}
	return d
}

func EscoreIVCF(r *Respostas) int {
	d := PontuarDominiosIVCF(r)
	return d.Idade + d.Autopercepcao + d.Aivd + d.Abvd + d.Cognicao + d.Humor +
		d.Alcance + d.Aerobica + d.Marcha + d.Comunicacao + d.Comorbidades
}

type ClasseIVCF string

const (
	ClasseIVCFBaixo    ClasseIVCF = "baixo"
	ClasseIVCFModerado ClasseIVCF = "moderado"
	ClasseIVCFAlto     ClasseIVCF = "alto"
)

func ClassificarIVCF(escore int) ClasseIVCF {
	if escore <= 6 {
		return ClasseIVCFBaixo
	}
	if escore <= 14 {
		return ClasseIVCFModerado
	}
	return ClasseIVCFAlto
}

func EscoreEASFI(r *Respostas) int {
	total := 0
	for _, v := range r.Easfi {
		if v != nil {
			total += *v
		}
	}
	return total
}

type ClasseEASFI string

const (
	ClasseEASFIAdequada ClasseEASFI = "adequada"
	ClasseEASFIRisco    ClasseEASFI = "risco"
	ClasseEASFIProblema ClasseEASFI = "problema"
)

func ClassificarEASFI(escore int) ClasseEASFI {
	if escore <= 9 {
		return ClasseEASFIAdequada
	}
	if escore <= 14 {
		return ClasseEASFIRisco
	}
	return ClasseEASFIProblema
}

type Nivel int

func NivelMatriz(ivcf ClasseIVCF, easfi ClasseEASFI) Nivel {
	if ivcf == ClasseIVCFBaixo {
		switch easfi {
		case ClasseEASFIAdequada:
			return 1
		case ClasseEASFIRisco:
			return 2
		default:
			return 3
		}
	}
	if ivcf == ClasseIVCFModerado {
		switch easfi {
		case ClasseEASFIAdequada:
			return 2
		case ClasseEASFIRisco:
			return 3
		default:
			return 4
		}
	}
	switch easfi {
	case ClasseEASFIAdequada:
		return 3
	default:
		return 4
	}
}

var Sentinelas = [5]string{
	"Situação de rua, extrema pobreza ou insegurança alimentar grave",
	"Abandono, suspeita/confirmação de violência doméstica/intrafamiliar ou exploração financeira",
	"Ausência total de cuidador para idoso dependente nas AVDs",
	"Sobrecarga grave do cuidador com risco iminente de ruptura do cuidado",
	"Interrupção de tratamento por impossibilidade financeira ou falta de acesso habitacional/geográfico",
}

var SentinelasGraves = [2]int{0, 1}

func AplicarSentinelas(nivel Nivel, sentinelas []*bool) Nivel {
	algum := false
	for _, s := range sentinelas {
		if sim(s) {
			algum = true
			break
		}
	}
	if !algum {
		return nivel
	}
	if nivel == 4 {
		return 4
	}
	for _, i := range SentinelasGraves {
		if i < len(sentinelas) && sim(sentinelas[i]) {
			return 4
		}
	}
	if nivel < 3 {
		return 3
	}
	return nivel
}

type Resultado struct {
	Ivcf                int         `json:"ivcf"`
	ClasseIVCF          ClasseIVCF  `json:"classeIVCF"`
	Easfi               int         `json:"easfi"`
	ClasseEASFI         ClasseEASFI `json:"classeEASFI"`
	NivelMatriz         Nivel       `json:"nivelMatriz"`
	Nivel               Nivel       `json:"nivel"`
	SentinelasPositivas []int       `json:"sentinelasPositivas"`
}

func Calcular(r *Respostas) Resultado {
	ivcf := EscoreIVCF(r)
	easfi := EscoreEASFI(r)
	ci := ClassificarIVCF(ivcf)
	ce := ClassificarEASFI(easfi)
	nm := NivelMatriz(ci, ce)
	pos := []int{}
	for i, s := range r.Sentinelas {
		if sim(s) {
			pos = append(pos, i)
		}
	}
	return Resultado{
		Ivcf:                ivcf,
		ClasseIVCF:          ci,
		Easfi:               easfi,
		ClasseEASFI:         ce,
		NivelMatriz:         nm,
		Nivel:               AplicarSentinelas(nm, r.Sentinelas),
		SentinelasPositivas: pos,
	}
}
