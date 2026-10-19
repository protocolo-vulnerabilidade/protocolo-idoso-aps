// Cliente HTTP tipado do backend (apps/backend). Contrato mantido manualmente
// com internal/httpapi (ADR-0004) — os campos JSON são camelCase idênticos.

import type { Respostas, Resultado } from "./protocolo";

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

export interface ConcluidaDetalhe extends Concluida {
  resumo: string;
}

const base = import.meta.env["VITE_API_URL"] ?? "http://localhost:8080";

const UUID_RE = /^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$/i;

function idAnonimo(): string {
  const guardado = localStorage.getItem("piaps:anon");
  if (guardado && UUID_RE.test(guardado)) return guardado;
  const id = crypto.randomUUID();
  localStorage.setItem("piaps:anon", id);
  return id;
}

export class ErroApi extends Error {
  status: number;
  detalhes: { campo: string; mensagem: string }[] | undefined;

  constructor(
    status: number,
    corpo: { erro?: string; detalhes?: { campo: string; mensagem: string }[] } | null,
  ) {
    super(corpo?.erro ?? `erro HTTP ${status}`);
    this.status = status;
    this.detalhes = corpo?.detalhes;
  }
}

async function req<T>(caminho: string, init?: RequestInit): Promise<T> {
  const res = await fetch(`${base}${caminho}`, {
    ...init,
    headers: {
      "Content-Type": "application/json",
      "X-Dispositivo-Id": idAnonimo(),
      ...init?.headers,
    },
  });
  if (!res.ok) {
    const corpo = await res.json().catch(() => null);
    throw new ErroApi(res.status, corpo);
  }
  if (res.status === 204) return undefined as T;
  return (await res.json()) as T;
}

export const api = {
  listarRascunhos: () => req<{ rascunhos: Rascunho[] }>("/api/rascunhos"),
  criarRascunho: () => req<Rascunho>("/api/rascunhos", { method: "POST" }),
  obterRascunho: (id: string) => req<Rascunho>(`/api/rascunhos/${id}`),
  salvarRascunho: (id: string, etapa: number, respostas: Respostas) =>
    req<Rascunho>(`/api/rascunhos/${id}`, {
      method: "PUT",
      body: JSON.stringify({ etapa, respostas }),
    }),
  apagarRascunho: (id: string) => req<void>(`/api/rascunhos/${id}`, { method: "DELETE" }),
  concluir: (id: string, respostas: Respostas) =>
    req<ConcluidaDetalhe>(`/api/rascunhos/${id}/concluir`, {
      method: "POST",
      body: JSON.stringify({ respostas }),
    }),
  listarHistorico: () => req<{ historico: Concluida[] }>("/api/historico"),
  obterConcluida: (id: string) => req<ConcluidaDetalhe>(`/api/historico/${id}`),
  apagarConcluida: (id: string) => req<void>(`/api/historico/${id}`, { method: "DELETE" }),
  reabrir: (id: string) => req<Rascunho>(`/api/historico/${id}/reabrir`, { method: "POST" }),
  apagarHistorico: () => req<void>("/api/historico", { method: "DELETE" }),
};
