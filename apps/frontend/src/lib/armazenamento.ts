import { useEffect, useState } from "react";
import { calcular, type Respostas, type Resultado } from "./protocolo";

export interface Rascunho {
  id: string;
  criadoEm: string;
  atualizadoEm: string;
  etapa: number;
  respostas: Respostas;
}
export interface Concluida {
  id: string;
  concluidaEm: string;
  respostas: Respostas;
  resultado: Resultado;
}
interface Dados {
  rascunhos: Rascunho[];
  historico: Concluida[];
}

const uid = () =>
  typeof crypto !== "undefined" && "randomUUID" in crypto
    ? crypto.randomUUID()
    : Math.random().toString(36).slice(2) + Date.now().toString(36);

function idAnonimo(): string {
  let id = localStorage.getItem("piaps:anon");
  if (!id) {
    id = uid();
    localStorage.setItem("piaps:anon", id);
  }
  return id;
}
const chave = () => `piaps:dados:${idAnonimo()}`;

function ler(): Dados {
  try {
    const raw = localStorage.getItem(chave());
    if (raw) {
      const d = JSON.parse(raw) as Dados;
      // Descarta dados salvos com a estrutura antiga do IVCF-20 (humor como valor único).
      const valido = (r: { respostas: Respostas }) => Array.isArray(r.respostas?.humor);
      return {
        rascunhos: (d.rascunhos ?? []).filter(valido),
        historico: (d.historico ?? []).filter(valido),
      };
    }
  } catch {
    /* ignora */
  }
  return { rascunhos: [], historico: [] };
}
function gravar(d: Dados) {
  localStorage.setItem(chave(), JSON.stringify(d));
  window.dispatchEvent(new Event("piaps:mudou"));
}

export const store = {
  criarRascunho(respostas: Respostas): string {
    const d = ler();
    const agora = new Date().toISOString();
    const r: Rascunho = { id: uid(), criadoEm: agora, atualizadoEm: agora, etapa: 0, respostas };
    d.rascunhos.unshift(r);
    gravar(d);
    return r.id;
  },
  salvarRascunho(id: string, etapa: number, respostas: Respostas) {
    const d = ler();
    const r = d.rascunhos.find((x) => x.id === id);
    if (!r) return;
    Object.assign(r, { etapa, respostas, atualizadoEm: new Date().toISOString() });
    gravar(d);
  },
  apagarRascunho(id: string) {
    const d = ler();
    d.rascunhos = d.rascunhos.filter((x) => x.id !== id);
    gravar(d);
  },
  concluir(id: string, respostas: Respostas): string {
    const d = ler();
    d.rascunhos = d.rascunhos.filter((x) => x.id !== id);
    d.historico = d.historico.filter((x) => x.id !== id);
    d.historico.unshift({
      id,
      concluidaEm: new Date().toISOString(),
      respostas,
      resultado: calcular(respostas),
    });
    gravar(d);
    return id;
  },
  /** Reabre uma avaliação concluída como rascunho para correção das respostas. */
  editarConcluida(id: string): boolean {
    const d = ler();
    const c = d.historico.find((x) => x.id === id);
    if (!c) return false;
    if (!d.rascunhos.some((x) => x.id === id)) {
      const agora = new Date().toISOString();
      d.rascunhos.unshift({ id, criadoEm: c.concluidaEm, atualizadoEm: agora, etapa: 0, respostas: c.respostas });
    }
    d.historico = d.historico.filter((x) => x.id !== id);
    gravar(d);
    return true;
  },
  apagarConcluida(id: string) {
    const d = ler();
    d.historico = d.historico.filter((x) => x.id !== id);
    gravar(d);
  },
  apagarHistorico() {
    const d = ler();
    d.historico = [];
    gravar(d);
  },
};

/** Lê os dados apenas no cliente; `null` enquanto não montado. */
export function useDados(): Dados | null {
  const [d, setD] = useState<Dados | null>(null);
  useEffect(() => {
    const at = () => setD(ler());
    at();
    window.addEventListener("piaps:mudou", at);
    window.addEventListener("storage", at);
    return () => {
      window.removeEventListener("piaps:mudou", at);
      window.removeEventListener("storage", at);
    };
  }, []);
  return d;
}
