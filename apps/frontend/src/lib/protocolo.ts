// Lógica pura de pontuação do Protocolo Idoso APS (IVCF-20, EASFI-BR, MIV-CFS, sentinelas).
// Sem dependências de UI — será substituída por chamadas à API (backend Go).

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
  comunicacao: YN[]; // visão, audição
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

const sim = (v: YN | undefined) => v === true;
const soma = (arr: YN[], pesos: number[]) => arr.reduce((t, v, i) => t + (sim(v) ? (pesos[i] ?? 0) : 0), 0);

export function pontosIdade(idade: number): number {
  if (idade >= 85) return 3;
  if (idade >= 75) return 1;
  return 0;
}

export interface DominiosIVCF {
  idade: number;
  autopercepcao: number;
  aivd: number;
  abvd: number;
  cognicao: number;
  humor: number;
  alcance: number;
  aerobica: number;
  marcha: number;
  comunicacao: number;
  comorbidades: number;
}

export function dominiosIVCF(r: Respostas): DominiosIVCF {
  return {
    idade: r.idade == null ? 0 : pontosIdade(r.idade),
    autopercepcao: r.autopercepcao === "ruim" ? 1 : 0,
    // AVD instrumental: cada Sim vale 4, com teto de 4 no bloco.
    aivd: Math.min(soma(r.aivd, [4, 4, 4]), 4),
    abvd: sim(r.abvd) ? 6 : 0,
    cognicao: soma(r.cognicao, [1, 1, 2]),
    humor: soma(r.humor, [2, 2]),
    alcance: soma(r.alcance, [1, 1]),
    aerobica: sim(r.aerobica) ? 2 : 0,
    marcha: soma(r.marcha, [2, 2, 2]),
    comunicacao: soma(r.comunicacao, [2, 2]),
    comorbidades: sim(r.comorbidades) ? 4 : 0,
  };
}

export function escoreIVCF(r: Respostas): number {
  return Object.values(dominiosIVCF(r)).reduce((a, b) => a + b, 0);
}

export type ClasseIVCF = "baixo" | "moderado" | "alto";
export function classificarIVCF(escore: number): ClasseIVCF {
  if (escore <= 6) return "baixo";
  if (escore <= 14) return "moderado";
  return "alto";
}
export const rotuloIVCF: Record<ClasseIVCF, string> = {
  baixo: "Baixo risco",
  moderado: "Risco moderado",
  alto: "Alto risco",
};

export function escoreEASFI(r: Respostas): number {
  return r.easfi.reduce<number>((t, v) => t + (v ?? 0), 0);
}

export type ClasseEASFI = "adequada" | "risco" | "problema";
export function classificarEASFI(escore: number): ClasseEASFI {
  if (escore <= 9) return "adequada";
  if (escore <= 14) return "risco";
  return "problema";
}
export const rotuloEASFI: Record<ClasseEASFI, string> = {
  adequada: "Situação social adequada",
  risco: "Risco social",
  problema: "Problema social",
};

export type Nivel = 1 | 2 | 3 | 4;

const MATRIZ: Record<ClasseIVCF, Record<ClasseEASFI, Nivel>> = {
  baixo: { adequada: 1, risco: 2, problema: 3 },
  moderado: { adequada: 2, risco: 3, problema: 4 },
  alto: { adequada: 3, risco: 4, problema: 4 },
};

export function nivelMatriz(ivcf: ClasseIVCF, easfi: ClasseEASFI): Nivel {
  return MATRIZ[ivcf][easfi];
}

export const SENTINELAS = [
  "Situação de rua, extrema pobreza ou insegurança alimentar grave",
  "Abandono, suspeita/confirmação de violência doméstica/intrafamiliar ou exploração financeira",
  "Ausência total de cuidador para idoso dependente nas AVDs",
  "Sobrecarga grave do cuidador com risco iminente de ruptura do cuidado",
  "Interrupção de tratamento por impossibilidade financeira ou falta de acesso habitacional/geográfico",
];
/** Índices de marcadores graves (situação de rua, violência/abandono) que levam direto ao Nível 4. */
export const SENTINELAS_GRAVES = [0, 1];

export function aplicarSentinelas(nivel: Nivel, sentinelas: YN[]): Nivel {
  if (!sentinelas.some(sim)) return nivel;
  if (nivel === 4) return 4;
  if (SENTINELAS_GRAVES.some((i) => sim(sentinelas[i]))) return 4;
  return 3 > nivel ? 3 : nivel;
}

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
    conduta:
      "PTS prioritário com acompanhamento intensivo e articulação intersetorial urgente.",
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

export function calcular(r: Respostas): Resultado {
  const ivcf = escoreIVCF(r);
  const easfi = escoreEASFI(r);
  const cI = classificarIVCF(ivcf);
  const cE = classificarEASFI(easfi);
  const nm = nivelMatriz(cI, cE);
  return {
    ivcf,
    classeIVCF: cI,
    easfi,
    classeEASFI: cE,
    nivelMatriz: nm,
    nivel: aplicarSentinelas(nm, r.sentinelas),
    sentinelasPositivas: r.sentinelas.flatMap((v, i) => (sim(v) ? [i] : [])),
  };
}

export function gerarResumo(res: Resultado, data: Date): string {
  const n = NIVEIS[res.nivel];
  const linhas = [
    "PROTOCOLO IDOSO APS — Estratificação de vulnerabilidade",
    `Data/hora: ${data.toLocaleString("pt-BR")}`,
    `IVCF-20: ${res.ivcf}/40 — ${rotuloIVCF[res.classeIVCF]}`,
    `EASFI-BR: ${res.easfi}/25 — ${rotuloEASFI[res.classeEASFI]}`,
    `Nível pela matriz MIV-CFS: ${res.nivelMatriz}`,
  ];
  if (res.sentinelasPositivas.length) {
    linhas.push("Marcadores sentinela presentes:");
    res.sentinelasPositivas.forEach((i) => linhas.push(`  - ${SENTINELAS[i]}`));
  }
  linhas.push(
    `NÍVEL DE VULNERABILIDADE: ${res.nivel} — ${n.nome}`,
    `Conduta: ${n.conduta}`,
    `Reavaliação: ${n.reavaliacao}`,
  );
  return linhas.join("\n");
}
