package dominio

import (
	"fmt"
	"strings"
	"time"
	_ "time/tzdata"
)

var RotuloIVCF = map[ClasseIVCF]string{
	ClasseIVCFBaixo:    "Baixo risco",
	ClasseIVCFModerado: "Risco moderado",
	ClasseIVCFAlto:     "Alto risco",
}

var RotuloEASFI = map[ClasseEASFI]string{
	ClasseEASFIAdequada: "Situação social adequada",
	ClasseEASFIRisco:    "Risco social",
	ClasseEASFIProblema: "Problema social",
}

type NivelInfo struct {
	Nome        string
	Conduta     string
	Reavaliacao string
}

var Niveis = map[Nivel]NivelInfo{
	1: {
		Nome:        "Baixa",
		Conduta:     "Acompanhamento de rotina e promoção da saúde.",
		Reavaliacao: "Anual",
	},
	2: {
		Nome:        "Moderada",
		Conduta:     "Avaliação clínica ampliada e revisão farmacoterapêutica.",
		Reavaliacao: "Em até 6 meses",
	},
	3: {
		Nome:        "Alta",
		Conduta:     "Discussão interdisciplinar, PTS obrigatório, visita domiciliar e articulação eMulti/CRAS.",
		Reavaliacao: "A cada 3 meses",
	},
	4: {
		Nome:        "Muito Alta",
		Conduta:     "PTS prioritário com acompanhamento intensivo e articulação intersetorial urgente.",
		Reavaliacao: "Em 30 a 90 dias",
	},
}

func fusoBrasil() *time.Location {
	local, err := time.LoadLocation("America/Sao_Paulo")
	if err != nil {
		return time.UTC
	}
	return local
}

func GerarResumo(res Resultado, data time.Time) string {
	n := Niveis[res.Nivel]
	linhas := []string{
		"PROTOCOLO IDOSO APS — Estratificação de vulnerabilidade",
		fmt.Sprintf("Data/hora: %s", data.In(fusoBrasil()).Format("02/01/2006 15:04:05")),
		fmt.Sprintf("IVCF-20: %d/40 — %s", res.Ivcf, RotuloIVCF[res.ClasseIVCF]),
		fmt.Sprintf("EASFI-BR: %d/25 — %s", res.Easfi, RotuloEASFI[res.ClasseEASFI]),
		fmt.Sprintf("Nível pela matriz MIV-CFS: %d", res.NivelMatriz),
	}
	if len(res.SentinelasPositivas) > 0 {
		linhas = append(linhas, "Marcadores sentinela presentes:")
		for _, i := range res.SentinelasPositivas {
			if i >= 0 && i < len(Sentinelas) {
				linhas = append(linhas, fmt.Sprintf("  - %s", Sentinelas[i]))
			}
		}
	}
	linhas = append(linhas,
		fmt.Sprintf("NÍVEL DE VULNERABILIDADE: %d — %s", res.Nivel, n.Nome),
		fmt.Sprintf("Conduta: %s", n.Conduta),
		fmt.Sprintf("Reavaliação: %s", n.Reavaliacao),
	)
	return strings.Join(linhas, "\n")
}
