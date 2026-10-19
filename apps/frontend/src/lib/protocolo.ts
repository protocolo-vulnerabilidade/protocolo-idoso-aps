// Tipos e dados de exibição do Protocolo Idoso APS.
// O cálculo dos escores, da matriz MIV-CFS e das sentinelas — e a geração do
// resumo copiável — pertencem ao backend (docs/arquitetura.md). Este arquivo
// apenas descreve o contrato e os textos exibidos pela interface.

export type YN = boolean | null;

export interface Respostas {
  idade: number | null;
  autopercepcao: "boa" | "ruim" | null;
  aivd: YN[]; // 3 itens: compras, dinheiro/contas, pequenos trabalhos domésticos
  abvd: YN; // 1 item: banho
  cognicao: YN[]; // 3 itens
  humor: YN[]; // 2 itens: desânimo/tristeza, perda de interesse
  alcance: YN[]; // 2 itens: elevar braços, manusear objetos
  aerobica: YN; // 1 item: perda de peso / IMC / panturrilha / velocidade da marcha
  marcha: YN[]; // 3 itens: caminhar, quedas, incontinência
  comunicacao: YN[]; // 2 itens: visão, audição
  comorbidades: YN; // 1 item: 5+ doenças / 5+ medicamentos / internação recente
  easfi: (number | null)[]; // 5 domínios, 1–5
  sentinelas: YN[]; // 5
}

export const respostasVazias = (): Respostas => ({
  idade: null,
  autopercepcao: null,
  aivd: [null, null, null],
  abvd: null,
  cognicao: [null, null, null],
  humor: [null, null],
  alcance: [null, null],
  aerobica: null,
  marcha: [null, null, null],
  comunicacao: [null, null],
  comorbidades: null,
  easfi: [null, null, null, null, null],
  sentinelas: [null, null, null, null, null],
});

export type ClasseIVCF = "baixo" | "moderado" | "alto";
export const rotuloIVCF: Record<ClasseIVCF, string> = {
  baixo: "Baixo risco",
  moderado: "Risco moderado",
  alto: "Alto risco",
};

export type ClasseEASFI = "adequada" | "risco" | "problema";
export const rotuloEASFI: Record<ClasseEASFI, string> = {
  adequada: "Situação social adequada",
  risco: "Risco social",
  problema: "Problema social",
};

export type Nivel = 1 | 2 | 3 | 4;

export const SENTINELAS = [
  "Situação de rua, extrema pobreza ou insegurança alimentar grave",
  "Abandono, suspeita/confirmação de violência doméstica/intrafamiliar ou exploração financeira",
  "Ausência total de cuidador para idoso dependente nas AVDs",
  "Sobrecarga grave do cuidador com risco iminente de ruptura do cuidado",
  "Interrupção de tratamento por impossibilidade financeira ou falta de acesso habitacional/geográfico",
];

export const NIVEIS: Record<Nivel, { nome: string; conduta: string; reavaliacao: string }> = {
  1: {
    nome: "Baixa",
    conduta: "Acompanhamento de rotina e promoção da saúde.",
    reavaliacao: "Anual",
  },
  2: {
    nome: "Moderada",
    conduta: "Avaliação clínica ampliada e revisão farmacoterapêutica.",
    reavaliacao: "Em até 6 meses",
  },
  3: {
    nome: "Alta",
    conduta:
      "Discussão interdisciplinar, PTS obrigatório, visita domiciliar e articulação eMulti/CRAS.",
    reavaliacao: "A cada 3 meses",
  },
  4: {
    nome: "Muito Alta",
    conduta: "PTS prioritário com acompanhamento intensivo e articulação intersetorial urgente.",
    reavaliacao: "Em 30 a 90 dias",
  },
};

export interface Resultado {
  ivcf: number;
  classeIVCF: ClasseIVCF;
  easfi: number;
  classeEASFI: ClasseEASFI;
  nivelMatriz: Nivel;
  nivel: Nivel;
  sentinelasPositivas: number[];
}
