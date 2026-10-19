import { createFileRoute, Link, useNavigate } from "@tanstack/react-router";
import { Botao, Cartao, Pagina } from "@/components/ui-base";
import { store, useDados } from "@/lib/armazenamento";
import { respostasVazias } from "@/lib/protocolo";

export const Route = createFileRoute("/")({
  head: () => ({
    meta: [
      { title: "Protocolo Idoso APS — Avaliação de vulnerabilidade" },
      { name: "description", content: "Estratificação anônima da vulnerabilidade da pessoa idosa na Atenção Primária (IVCF-20 e EASFI-BR)." },
      { property: "og:title", content: "Protocolo Idoso APS" },
      { property: "og:description", content: "Estratificação anônima da vulnerabilidade da pessoa idosa na APS." },
    ],
  }),
  component: Inicio,
});

function Inicio() {
  const dados = useDados();
  const nav = useNavigate();
  const nova = () => {
    const id = store.criarRascunho(respostasVazias());
    nav({ to: "/avaliacao/$id", params: { id } });
  };

  return (
    <Pagina>
      <h1 className="text-2xl font-bold tracking-tight sm:text-3xl">Avaliação da pessoa idosa</h1>
      <p className="mt-2 text-muted-foreground">
        Aplique o IVCF-20 e o EASFI-BR para obter o nível de vulnerabilidade. Nenhum dado de
        identificação é solicitado ou armazenado.
      </p>
      <div className="mt-6 flex flex-col gap-3 sm:flex-row">
        <Botao onClick={nova} className="sm:min-w-56">Nova avaliação</Botao>
        <Link to="/historico">
          <Botao variante="secundario" className="w-full">Ver histórico</Botao>
        </Link>
      </div>

      <h2 className="mt-10 text-lg font-semibold">Rascunhos</h2>
      {dados && dados.rascunhos.length === 0 && (
        <p className="mt-2 text-sm text-muted-foreground">Nenhum rascunho salvo.</p>
      )}
      <div className="mt-3 space-y-3">
        {dados?.rascunhos.map((r) => (
          <Cartao key={r.id} className="flex flex-col gap-3 sm:flex-row sm:items-center sm:justify-between">
            <div>
              <p className="font-medium">Continuar de onde parou</p>
              <p className="text-sm text-muted-foreground">
                Etapa {r.etapa + 1} de 7 · atualizado em {new Date(r.atualizadoEm).toLocaleString("pt-BR")}
              </p>
            </div>
            <div className="flex gap-2">
              <Link to="/avaliacao/$id" params={{ id: r.id }}>
                <Botao>Continuar</Botao>
              </Link>
              <Botao variante="perigo" onClick={() => confirm("Descartar este rascunho?") && store.apagarRascunho(r.id)}>
                Descartar
              </Botao>
            </div>
          </Cartao>
        ))}
      </div>
    </Pagina>
  );
}
