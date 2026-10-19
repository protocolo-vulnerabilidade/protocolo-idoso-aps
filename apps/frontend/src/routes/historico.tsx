import { createFileRoute, Link } from "@tanstack/react-router";
import { Botao, Cartao, Pagina, corNivel } from "@/components/ui-base";
import { store, useDados } from "@/lib/armazenamento";
import { NIVEIS } from "@/lib/protocolo";
import { cn } from "@/lib/utils";

export const Route = createFileRoute("/historico")({
  head: () => ({
    meta: [
      { title: "Histórico de avaliações — Protocolo Idoso APS" },
      { name: "description", content: "Avaliações concluídas neste dispositivo, com exclusão individual ou total." },
      { property: "og:title", content: "Histórico — Protocolo Idoso APS" },
      { property: "og:description", content: "Avaliações concluídas neste dispositivo." },
    ],
  }),
  component: Historico,
});

function Historico() {
  const dados = useDados();
  const lista = dados?.historico ?? [];
  return (
    <Pagina>
      <div className="flex flex-wrap items-center justify-between gap-3">
        <h1 className="text-2xl font-bold tracking-tight">Histórico</h1>
        <Botao variante="perigo" disabled={!lista.length}
          onClick={() => confirm("Apagar TODAS as avaliações do histórico? Esta ação não pode ser desfeita.") && store.apagarHistorico()}>
          Apagar tudo
        </Botao>
      </div>
      <p className="mt-1 text-sm text-muted-foreground">Dados salvos apenas neste dispositivo, sem identificação.</p>
      {dados && !lista.length && <p className="mt-6 text-muted-foreground">Nenhuma avaliação concluída.</p>}
      <div className="mt-4 space-y-3">
        {lista.map((a) => (
          <Cartao key={a.id} className="flex flex-col gap-3 sm:flex-row sm:items-center sm:justify-between">
            <div className="flex items-center gap-3">
              <span className={cn("grid h-12 w-12 shrink-0 place-items-center rounded-md text-xl font-bold", corNivel[a.resultado.nivel])}>
                {a.resultado.nivel}
              </span>
              <div>
                <p className="font-medium">Nível {a.resultado.nivel} — {NIVEIS[a.resultado.nivel].nome}</p>
                <p className="text-sm text-muted-foreground">
                  {new Date(a.concluidaEm).toLocaleString("pt-BR")} · IVCF {a.resultado.ivcf} · EASFI {a.resultado.easfi}
                </p>
              </div>
            </div>
            <div className="flex gap-2">
              <Link to="/resultado/$id" params={{ id: a.id }}><Botao variante="secundario">Abrir</Botao></Link>
              <Botao variante="perigo" onClick={() => confirm("Apagar esta avaliação?") && store.apagarConcluida(a.id)}>Apagar</Botao>
            </div>
          </Cartao>
        ))}
      </div>
    </Pagina>
  );
}
