import { createFileRoute, Link, useNavigate } from "@tanstack/react-router";
import { Botao, Cartao, Pagina } from "@/components/ui-base";
import { useApagarRascunho, useCriarRascunho, useRascunhos } from "@/lib/armazenamento";

export const Route = createFileRoute("/")({
  head: () => ({
    meta: [
      { title: "Protocolo Idoso APS — Avaliação de vulnerabilidade" },
      {
        name: "description",
        content:
          "Estratificação anônima da vulnerabilidade da pessoa idosa na Atenção Primária (IVCF-20 e EASFI-BR).",
      },
      { property: "og:title", content: "Protocolo Idoso APS" },
      {
        property: "og:description",
        content: "Estratificação anônima da vulnerabilidade da pessoa idosa na APS.",
      },
    ],
  }),
  component: Inicio,
});

function Inicio() {
  const { data, isPending } = useRascunhos();
  const criar = useCriarRascunho();
  const apagar = useApagarRascunho();
  const nav = useNavigate();
  const nova = async () => {
    try {
      const rascunho = await criar.mutateAsync();
      nav({ to: "/avaliacao/$id", params: { id: rascunho.id } });
    } catch {
      alert("Não foi possível iniciar a avaliação. Verifique se o servidor está disponível.");
    }
  };

  return (
    <Pagina>
      <h1 className="text-2xl font-bold tracking-tight sm:text-3xl">Avaliação da pessoa idosa</h1>
      <p className="mt-2 text-muted-foreground">
        Aplique o IVCF-20 e o EASFI-BR para obter o nível de vulnerabilidade. Nenhum dado de
        identificação é solicitado ou armazenado.
      </p>
      <div className="mt-6 flex flex-col gap-3 sm:flex-row">
        <Botao
          onClick={() => {
            void nova();
          }}
          disabled={criar.isPending}
          className="sm:min-w-56"
        >
          {criar.isPending ? "Iniciando…" : "Nova avaliação"}
        </Botao>
        <Link to="/historico">
          <Botao variante="secundario" className="w-full">
            Ver histórico
          </Botao>
        </Link>
      </div>

      <h2 className="mt-10 text-lg font-semibold">Rascunhos</h2>
      {isPending && <p className="mt-2 text-sm text-muted-foreground">Carregando…</p>}
      {!isPending && data?.rascunhos.length === 0 && (
        <p className="mt-2 text-sm text-muted-foreground">Nenhum rascunho salvo.</p>
      )}
      <div className="mt-3 space-y-3">
        {data?.rascunhos.map((r) => (
          <Cartao
            key={r.id}
            className="flex flex-col gap-3 sm:flex-row sm:items-center sm:justify-between"
          >
            <div>
              <p className="font-medium">Continuar de onde parou</p>
              <p className="text-sm text-muted-foreground">
                Etapa {r.etapa + 1} de 7 · atualizado em{" "}
                {new Date(r.atualizadoEm).toLocaleString("pt-BR")}
              </p>
            </div>
            <div className="flex gap-2">
              <Link to="/avaliacao/$id" params={{ id: r.id }}>
                <Botao>Continuar</Botao>
              </Link>
              <Botao
                variante="perigo"
                disabled={apagar.isPending && apagar.variables === r.id}
                onClick={() => confirm("Descartar este rascunho?") && apagar.mutate(r.id)}
              >
                Descartar
              </Botao>
            </div>
          </Cartao>
        ))}
      </div>
    </Pagina>
  );
}
