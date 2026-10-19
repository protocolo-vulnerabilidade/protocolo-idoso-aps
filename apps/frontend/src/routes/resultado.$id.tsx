import { createFileRoute, Link, useNavigate } from "@tanstack/react-router";
import { useState } from "react";
import { Botao, Cartao, Pagina, corNivel } from "@/components/ui-base";
import { useConcluida, useReabrir } from "@/lib/armazenamento";
import { NIVEIS, SENTINELAS, rotuloEASFI, rotuloIVCF } from "@/lib/protocolo";
import { cn } from "@/lib/utils";

export const Route = createFileRoute("/resultado/$id")({
  head: () => ({
    meta: [
      { title: "Resultado da avaliação — Protocolo Idoso APS" },
      {
        name: "description",
        content: "Nível de vulnerabilidade, conduta recomendada e resumo para o sistema da UBS.",
      },
      { property: "og:title", content: "Resultado — Protocolo Idoso APS" },
      { property: "og:description", content: "Nível de vulnerabilidade e conduta recomendada." },
      { name: "robots", content: "noindex" },
    ],
  }),
  component: ResultadoPagina,
});

function ResultadoPagina() {
  const { id } = Route.useParams();
  const { data, isPending, isError } = useConcluida(id);
  const [copiado, setCopiado] = useState(false);
  const nav = useNavigate();
  const reabrir = useReabrir();
  if (isPending)
    return (
      <Pagina>
        <p className="text-muted-foreground">Carregando…</p>
      </Pagina>
    );
  if (isError || !data)
    return (
      <Pagina>
        <p>Avaliação não encontrada.</p>
        <Link to="/" className="text-primary underline">
          Início
        </Link>
      </Pagina>
    );
  const res = data.resultado;
  const n = NIVEIS[res.nivel];
  const resumo = data.resumo;

  const copiar = async () => {
    await navigator.clipboard.writeText(resumo);
    setCopiado(true);
    setTimeout(() => setCopiado(false), 2500);
  };

  const editar = async () => {
    if (
      !confirm(
        "Reabrir a avaliação para corrigir as respostas? O resultado será recalculado ao concluir novamente.",
      )
    )
      return;
    try {
      await reabrir.mutateAsync(id);
      nav({ to: "/avaliacao/$id", params: { id } });
    } catch {
      alert("Não foi possível reabrir a avaliação. Tente novamente.");
    }
  };

  return (
    <Pagina>
      <p className="text-sm text-muted-foreground">
        {new Date(data.concluidaEm).toLocaleString("pt-BR")}
      </p>
      <div className={cn("mt-2 rounded-lg p-6", corNivel[res.nivel])}>
        <p className="text-sm font-semibold uppercase tracking-wide opacity-90">
          Nível de vulnerabilidade
        </p>
        <p className="mt-1 text-4xl font-bold">
          Nível {res.nivel} — {n.nome}
        </p>
        <p className="mt-3 text-base">{n.conduta}</p>
        <p className="mt-2 text-sm font-semibold">Reavaliação: {n.reavaliacao}</p>
      </div>

      {res.sentinelasPositivas.length > 0 && (
        <div
          role="alert"
          className="mt-4 rounded-lg border-2 border-destructive bg-destructive/5 p-4"
        >
          <p className="font-semibold text-destructive">Alerta: marcador sentinela presente</p>
          <p className="text-sm text-muted-foreground">
            Nível pela matriz: {res.nivelMatriz}. Reclassificado para Nível {res.nivel}.
          </p>
          <ul className="mt-2 list-disc pl-5 text-sm">
            {res.sentinelasPositivas.map((i) => (
              <li key={i}>{SENTINELAS[i]}</li>
            ))}
          </ul>
        </div>
      )}

      <div className="mt-4 grid gap-4 sm:grid-cols-2">
        <Cartao>
          <p className="text-sm text-muted-foreground">IVCF-20</p>
          <p className="text-3xl font-bold">
            {res.ivcf}
            <span className="text-base font-normal text-muted-foreground">/40</span>
          </p>
          <p className="font-medium">{rotuloIVCF[res.classeIVCF]}</p>
        </Cartao>
        <Cartao>
          <p className="text-sm text-muted-foreground">EASFI-BR</p>
          <p className="text-3xl font-bold">
            {res.easfi}
            <span className="text-base font-normal text-muted-foreground">/25</span>
          </p>
          <p className="font-medium">{rotuloEASFI[res.classeEASFI]}</p>
        </Cartao>
      </div>

      <Cartao className="mt-4">
        <p className="mb-2 font-semibold">Resumo para a ficha</p>
        <pre className="whitespace-pre-wrap rounded-md bg-muted p-3 font-mono text-sm">
          {resumo}
        </pre>
        <Botao onClick={copiar} className="mt-4 w-full sm:w-auto">
          {copiado ? "Resumo copiado" : "Copiar resumo"}
        </Botao>
      </Cartao>

      <div className="mt-6 flex flex-wrap items-center gap-3">
        <Link to="/">
          <Botao variante="secundario">Finalizar</Botao>
        </Link>
        <Link to="/historico">
          <Botao variante="fantasma">Histórico</Botao>
        </Link>
        <Botao
          className="ml-auto"
          disabled={reabrir.isPending}
          onClick={() => {
            void editar();
          }}
        >
          Editar respostas
        </Botao>
      </div>
    </Pagina>
  );
}
