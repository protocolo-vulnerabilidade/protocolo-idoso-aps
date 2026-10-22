// Persistência via API do backend (ADR-0002: histórico e rascunhos no banco,
// vinculados ao identificador anônimo de dispositivo — sem contas de usuário).
// O localStorage guarda somente esse identificador (ver lib/api.ts).

import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import { api, ErroApi } from "./api";

export type { Rascunho, Concluida, ConcluidaDetalhe } from "./api";

const rascunhos = ["rascunhos"] as const;
const historico = ["historico"] as const;

export function useRascunhos() {
  return useQuery({ queryKey: rascunhos, queryFn: api.listarRascunhos });
}

export function useRascunho(id: string) {
  return useQuery({
    queryKey: [...rascunhos, id],
    queryFn: () => api.obterRascunho(id),
    retry: false,
  });
}

export function useHistorico() {
  return useQuery({ queryKey: historico, queryFn: api.listarHistorico });
}

export function useConcluida(id: string) {
  return useQuery({
    queryKey: [...historico, id],
    queryFn: () => api.obterConcluida(id),
    retry: false,
  });
}

export function useCriarRascunho() {
  const qc = useQueryClient();
  return useMutation({
    mutationFn: api.criarRascunho,
    onSuccess: () => qc.invalidateQueries({ queryKey: rascunhos }),
  });
}

export function useSalvarRascunho() {
  const qc = useQueryClient();
  return useMutation({
    mutationFn: (params: {
      id: string;
      etapa: number;
      respostas: Parameters<typeof api.salvarRascunho>[2];
    }) => api.salvarRascunho(params.id, params.etapa, params.respostas),
    onSettled: () => qc.invalidateQueries({ queryKey: rascunhos }),
  });
}

export function useApagarRascunho() {
  const qc = useQueryClient();
  return useMutation({
    mutationFn: api.apagarRascunho,
    onSuccess: () => qc.invalidateQueries({ queryKey: rascunhos }),
  });
}

export function useConcluir() {
  const qc = useQueryClient();
  return useMutation({
    mutationFn: (params: { id: string; respostas: Parameters<typeof api.concluir>[1] }) =>
      api.concluir(params.id, params.respostas),
    onSuccess: (_dados, vars) => {
      qc.removeQueries({ queryKey: [...rascunhos, vars.id] });
      qc.invalidateQueries({ queryKey: rascunhos });
      qc.invalidateQueries({ queryKey: historico });
    },
  });
}

export function useReabrir() {
  const qc = useQueryClient();
  return useMutation({
    mutationFn: api.reabrir,
    onSuccess: (_dados, id) => {
      qc.removeQueries({ queryKey: [...historico, id] });
      qc.invalidateQueries({ queryKey: historico });
      qc.invalidateQueries({ queryKey: rascunhos });
    },
  });
}

// Não repetir tentativa quando a senha é rejeitada (401/429): cada tentativa
// errada conta para o bloqueio anti-força-bruta do backend.
const repetirSoErroServidor = (_i: number, erro: unknown) =>
  !(erro instanceof ErroApi) || erro.status >= 500;

export function useApagarConcluida() {
  const qc = useQueryClient();
  return useMutation({
    mutationFn: (params: { id: string; senha: string }) =>
      api.apagarConcluida(params.id, params.senha),
    retry: repetirSoErroServidor,
    onSuccess: () => qc.invalidateQueries({ queryKey: historico }),
  });
}

export function useApagarHistorico() {
  const qc = useQueryClient();
  return useMutation({
    mutationFn: (senha: string) => api.apagarHistorico(senha),
    retry: repetirSoErroServidor,
    onSuccess: () => qc.invalidateQueries({ queryKey: historico }),
  });
}
