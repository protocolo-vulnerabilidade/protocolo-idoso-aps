import { useState } from "react";
import { createFileRoute, Link } from "@tanstack/react-router";
import { Botao, Cartao, Pagina, corNivel } from "@/components/ui-base";
import { DialogoSenhaExclusao } from "@/components/dialogo-senha";
import { useApagarConcluida, useApagarHistorico, useHistorico } from "@/lib/armazenamento";
import { NIVEIS } from "@/lib/protocolo";
import { cn } from "@/lib/utils";

export const Route = createFileRoute("/historico")({
  head: () => ({
    meta: [
      { title: "Histórico de avaliações — Protocolo Idoso APS" },
      {
        name: "description",
        content: "Avaliações concluídas neste dispositivo, com exclusão individual ou total.",
      },
      { property: "og:title", content: "Histórico — Protocolo Idoso APS" },
      { property: "og:description", content: "Avaliações concluídas neste dispositivo." },
    ],
  }),
  component: Historico,
});

type AlvoExclusao = { tipo: "uma"; id: string } | { tipo: "tudo" } | null;

function Historico() {
  const { data, isPending } = useHistorico();
  const apagar = useApagarConcluida();
  const apagarTudo = useApagarHistorico();
  const lista = data?.historico ?? [];
  const [alvo, setAlvo] = useState<AlvoExclusao>(null);
  const [erroExclusao, setErroExclusao] = useState<string | null>(null);

  const pendente = apagar.isPending || apagarTudo.isPending;

  function confirmarExclusao(senha: string) {
    if (!alvo) return;
    setErroExclusao(null);
    const callbacks = {
      onSuccess: () => setAlvo(null),
      onError: (erro: unknown) =>
        setErroExclusao(erro instanceof Error ? erro.message : "não foi possível apagar"),
    };
    if (alvo.tipo === "uma") {
      apagar.mutate({ id: alvo.id, senha }, callbacks);
    } else {
      apagarTudo.mutate(senha, callbacks);
    }
  }

  function fecharDialogo() {
    setAlvo(null);
    setErroExclusao(null);
  }

  return (
    <Pagina>
      <div className="flex flex-wrap items-center justify-between gap-3">
        <h1 className="text-2xl font-bold tracking-tight">Histórico</h1>
        <Botao
          variante="perigo"
          disabled={!lista.length || pendente}
          onClick={() => setAlvo({ tipo: "tudo" })}
        >
          Apagar tudo
        </Botao>
      </div>
      <p className="mt-1 text-sm text-muted-foreground">
        Dados salvos apenas neste dispositivo, sem identificação.
      </p>
      {isPending && <p className="mt-6 text-muted-foreground">Carregando…</p>}
      {!isPending && !lista.length && (
        <p className="mt-6 text-muted-foreground">Nenhuma avaliação concluída.</p>
      )}
      <div className="mt-4 space-y-3">
        {lista.map((a) => (
          <Cartao
            key={a.id}
            className="flex flex-col gap-3 sm:flex-row sm:items-center sm:justify-between"
          >
            <div className="flex items-center gap-3">
              <span
                className={cn(
                  "grid h-12 w-12 shrink-0 place-items-center rounded-md text-xl font-bold",
                  corNivel[a.resultado.nivel],
                )}
              >
                {a.resultado.nivel}
              </span>
              <div>
                <p className="font-medium">
                  Nível {a.resultado.nivel} — {NIVEIS[a.resultado.nivel].nome}
                </p>
                <p className="text-sm text-muted-foreground">
                  {new Date(a.concluidaEm).toLocaleString("pt-BR")} · IVCF {a.resultado.ivcf} ·
                  EASFI {a.resultado.easfi}
                </p>
              </div>
            </div>
            <div className="flex gap-2">
              <Link to="/resultado/$id" params={{ id: a.id }}>
                <Botao variante="secundario">Abrir</Botao>
              </Link>
              <Botao
                variante="perigo"
                disabled={pendente}
                onClick={() => setAlvo({ tipo: "uma", id: a.id })}
              >
                Apagar
              </Botao>
            </div>
          </Cartao>
        ))}
      </div>
      <DialogoSenhaExclusao
        aberto={alvo !== null}
        titulo={alvo?.tipo === "tudo" ? "Apagar todas as avaliações?" : "Apagar esta avaliação?"}
        descricao={
          alvo?.tipo === "tudo"
            ? "Todas as avaliações concluídas neste dispositivo serão apagadas. Esta ação não pode ser desfeita."
            : "A avaliação selecionada será apagada. Esta ação não pode ser desfeita."
        }
        pendente={pendente}
        erro={erroExclusao}
        onConfirmar={confirmarExclusao}
        onFechar={fecharDialogo}
      />
    </Pagina>
  );
}
