import { createFileRoute, Link, useNavigate } from "@tanstack/react-router";
import { useEffect, useRef, useState, type ReactNode } from "react";
import { Botao, Cartao, Pagina } from "@/components/ui-base";
import { useConcluir, useRascunho, useSalvarRascunho } from "@/lib/armazenamento";
import { SENTINELAS, type Respostas, type YN } from "@/lib/protocolo";
import { cn } from "@/lib/utils";

export const Route = createFileRoute("/avaliacao/$id")({
  head: () => ({
    meta: [
      { title: "Avaliação em andamento — Protocolo Idoso APS" },
      {
        name: "description",
        content: "Questionário passo a passo IVCF-20, EASFI-BR e marcadores sentinela.",
      },
      { property: "og:title", content: "Avaliação — Protocolo Idoso APS" },
      { property: "og:description", content: "Questionário passo a passo IVCF-20 e EASFI-BR." },
      { name: "robots", content: "noindex" },
    ],
  }),
  component: Wizard,
});

const ETAPAS = [
  "Idade e autopercepção",
  "Atividades de vida diária (AVD)",
  "Cognição e humor",
  "Mobilidade e comunicação",
  "Comorbidades",
  "Avaliação social (EASFI-BR)",
  "Marcadores sentinela",
];

const EASFI = [
  {
    t: "Situação familiar",
    o: [
      "Vive com cônjuge/família",
      "Vive com cônjuge e filhos",
      "Vive com parentes sem dependência direta",
      "Vive sozinho sem parentes próximos",
      "Família monoparental com sobrecarga/conflito grave",
    ],
  },
  {
    t: "Situação econômica",
    o: [
      "Renda acima de 3 salários mínimos",
      "De 2 a 3 salários mínimos",
      "De 1 a 2 salários mínimos",
      "Menos de 1 salário mínimo / BPC",
      "Sem renda própria / extrema pobreza",
    ],
  },
  {
    t: "Moradia",
    o: [
      "Própria, adaptada e segura",
      "Adequada sem adaptações",
      "Barreiras arquitetônicas moderadas",
      "Precária / saneamento inadequado",
      "Insalubre / risco de despejo / situação de rua",
    ],
  },
  {
    t: "Relações sociais",
    o: [
      "Vida social ativa e grupos comunitários",
      "Relações ocasionais",
      "Contato apenas com vizinhos/família",
      "Saídas raras do domicílio",
      "Isolamento completo / confinamento",
    ],
  },
  {
    t: "Rede de suporte social",
    o: [
      "Apoio familiar e comunitário constante",
      "Apoio familiar suficiente em emergências",
      "Apoio informal limitado",
      "Dependência de apoio formal público",
      "Ausência total de rede / abandono",
    ],
  },
];

function Opcao({
  ativo,
  onClick,
  children,
}: {
  ativo: boolean;
  onClick: () => void;
  children: ReactNode;
}) {
  return (
    <button
      type="button"
      onClick={onClick}
      aria-pressed={ativo}
      className={cn(
        "min-h-12 rounded-md border px-4 py-2 text-left text-base transition-colors",
        ativo
          ? "border-primary bg-primary text-primary-foreground"
          : "border-input bg-card hover:bg-accent",
      )}
    >
      {children}
    </button>
  );
}

function SimNao({
  texto,
  valor,
  onChange,
  erro,
  ariaTexto,
  sub,
}: {
  texto: string;
  valor: YN;
  onChange: (v: boolean) => void;
  erro: boolean;
  ariaTexto?: string;
  sub?: { sim: string; nao: string };
}) {
  return (
    <div
      role="group"
      aria-label={ariaTexto ?? texto}
      className={cn(
        "border-b border-border py-4 last:border-0",
        erro && valor === null && "text-destructive",
      )}
    >
      <p className="mb-3 font-medium" aria-hidden={ariaTexto ? true : undefined}>
        {texto}
      </p>
      <div className={cn("grid grid-cols-2 gap-2", sub ? "sm:max-w-xl" : "sm:w-64")}>
        <Opcao ativo={valor === true} onClick={() => onChange(true)}>
          <span className="block font-semibold">Sim</span>
          {sub && <span className="block text-sm opacity-80">{sub.sim}</span>}
        </Opcao>
        <Opcao ativo={valor === false} onClick={() => onChange(false)}>
          <span className="block font-semibold">Não</span>
          {sub && <span className="block text-sm opacity-80">{sub.nao}</span>}
        </Opcao>
      </div>
    </div>
  );
}

const ANCORA = "Por causa de sua saúde ou condição física";
const SUB_AVD = {
  sim: "deixei de fazer por causa da saúde",
  nao: "continuo fazendo (ou não faço por outros motivos)",
};

function BlocoAncora({ titulo, children }: { titulo: string; children: ReactNode }) {
  return (
    <div className="group/ancora rounded-lg border border-border">
      <p className="rounded-t-lg border-b border-border bg-muted px-4 py-3 text-xs font-semibold uppercase tracking-wide text-muted-foreground">
        {titulo}
      </p>
      <p className="bg-accent/40 px-4 py-3 font-bold transition-colors group-hover/ancora:bg-accent group-focus-within/ancora:bg-accent">
        {ANCORA}…
      </p>
      <div className="border-l-4 border-primary/30 ml-4 pl-4 pr-4">{children}</div>
    </div>
  );
}

function Wizard() {
  const { id } = Route.useParams();
  const { data, isPending, isError } = useRascunho(id);
  const nav = useNavigate();
  const [r, setR] = useState<Respostas | null>(null);
  const [etapa, setEtapa] = useState(0);
  const [erro, setErro] = useState(false);
  const [erroServidor, setErroServidor] = useState(false);
  const timerRef = useRef<ReturnType<typeof setTimeout> | null>(null);
  const { mutate: salvarAgora } = useSalvarRascunho();
  const concluir = useConcluir();

  useEffect(() => {
    if (data && !r) {
      setR(data.respostas);
      setEtapa(data.etapa);
    }
  }, [data, r]);

  useEffect(() => {
    if (!r) return;
    if (timerRef.current) clearTimeout(timerRef.current);
    timerRef.current = setTimeout(() => salvarAgora({ id, etapa, respostas: r }), 800);
    return () => {
      if (timerRef.current) clearTimeout(timerRef.current);
    };
  }, [r, etapa, id, salvarAgora]);

  if (isPending || (data && !r))
    return (
      <Pagina>
        <p className="text-muted-foreground">Carregando…</p>
      </Pagina>
    );
  if (isError || !data || !r) {
    return (
      <Pagina>
        <p>Rascunho não encontrado.</p>
        <Link to="/" className="mt-4 inline-block text-primary underline">
          Voltar ao início
        </Link>
      </Pagina>
    );
  }

  const set = (p: Partial<Respostas>) => setR({ ...r, ...p });
  const setArr = (k: keyof Respostas, i: number, v: boolean) => {
    const arr = [...(r[k] as YN[])];
    arr[i] = v;
    set({ [k]: arr } as Partial<Respostas>);
  };
  const yn = (k: keyof Respostas, i: number, texto: string) => (
    <SimNao
      key={`${k}${i}`}
      texto={texto}
      valor={(r[k] as YN[])[i] ?? null}
      erro={erro}
      onChange={(v) => setArr(k, i, v)}
    />
  );

  const completa = (): boolean => {
    const t = (a: YN[]) => a.every((x) => x !== null);
    switch (etapa) {
      case 0:
        return r.idade !== null && r.idade >= 60 && r.idade <= 130 && r.autopercepcao !== null;
      case 1:
        return t(r.aivd) && r.abvd !== null;
      case 2:
        return t(r.cognicao) && t(r.humor);
      case 3:
        return t(r.alcance) && r.aerobica !== null && t(r.marcha) && t(r.comunicacao);
      case 4:
        return r.comorbidades !== null;
      case 5:
        return r.easfi.every((x) => x !== null);
      default:
        return t(r.sentinelas);
    }
  };

  const avancar = () => {
    if (!completa()) {
      setErro(true);
      return;
    }
    setErro(false);
    if (etapa < ETAPAS.length - 1) {
      setEtapa(etapa + 1);
      window.scrollTo(0, 0);
    } else void finalizar();
  };

  const finalizar = async () => {
    if (timerRef.current) clearTimeout(timerRef.current);
    try {
      await concluir.mutateAsync({ id, respostas: r });
      nav({ to: "/resultado/$id", params: { id } });
    } catch {
      setErroServidor(true);
      window.scrollTo(0, 0);
    }
  };

  const pct = Math.round(((etapa + 1) / ETAPAS.length) * 100);

  return (
    <Pagina>
      <div className="mb-6">
        <div className="flex justify-between text-sm text-muted-foreground">
          <span>
            Etapa {etapa + 1} de {ETAPAS.length} ·{" "}
            {etapa < 5 ? "IVCF-20" : etapa === 5 ? "EASFI-BR" : "Sentinelas"}
          </span>
          <span>{pct}%</span>
        </div>
        <div className="mt-2 h-2 overflow-hidden rounded-full bg-muted">
          <div className="h-full bg-primary transition-all" style={{ width: `${pct}%` }} />
        </div>
        <h1 className="mt-4 text-2xl font-bold tracking-tight">{ETAPAS[etapa]}</h1>
      </div>

      <Cartao>
        {etapa === 0 && (
          <div className="space-y-6">
            <label className="block">
              <span className="font-medium">Idade (anos)</span>
              <input
                type="number"
                inputMode="numeric"
                min={60}
                max={130}
                value={r.idade ?? ""}
                onChange={(e) =>
                  set({ idade: e.target.value === "" ? null : Number(e.target.value) })
                }
                className="mt-2 block h-12 w-40 rounded-md border border-input bg-background px-3 text-lg focus:outline-none focus:ring-2 focus:ring-ring"
              />
              {erro && (r.idade === null || r.idade < 60 || r.idade > 130) && (
                <span className="mt-1 block text-sm text-destructive">
                  Informe uma idade entre 60 e 130 anos.
                </span>
              )}
            </label>
            <div>
              <p className={cn("mb-3 font-medium", erro && !r.autopercepcao && "text-destructive")}>
                Em geral, comparando com outras pessoas de sua idade, você diria que sua saúde é:
              </p>
              <div className="grid gap-2 sm:grid-cols-2">
                <Opcao
                  ativo={r.autopercepcao === "boa"}
                  onClick={() => set({ autopercepcao: "boa" })}
                >
                  Excelente, muito boa ou boa
                </Opcao>
                <Opcao
                  ativo={r.autopercepcao === "ruim"}
                  onClick={() => set({ autopercepcao: "ruim" })}
                >
                  Regular ou ruim
                </Opcao>
              </div>
            </div>
          </div>
        )}
        {etapa === 1 && (
          <div className="space-y-6">
            <BlocoAncora titulo="AVD instrumental">
              {[
                "você deixou de fazer compras?",
                "você deixou de controlar seu dinheiro, gastos ou pagar as contas de sua casa?",
                "você deixou de realizar pequenos trabalhos domésticos, como lavar a louça, arrumar a casa ou fazer limpeza leve?",
              ].map((t, i) => (
                <SimNao
                  key={`aivd${i}`}
                  texto={`…${t}`}
                  ariaTexto={`${ANCORA}, ${t}`}
                  sub={SUB_AVD}
                  valor={r.aivd[i] ?? null}
                  erro={erro}
                  onChange={(v) => setArr("aivd", i, v)}
                />
              ))}
            </BlocoAncora>
            <BlocoAncora titulo="AVD básica">
              <SimNao
                texto="…você deixou de tomar banho sozinho?"
                ariaTexto={`${ANCORA}, você deixou de tomar banho sozinho?`}
                sub={SUB_AVD}
                valor={r.abvd}
                erro={erro}
                onChange={(v) => set({ abvd: v })}
              />
            </BlocoAncora>
          </div>
        )}
        {etapa === 2 && (
          <>
            {[
              "Algum familiar ou amigo falou que você está ficando esquecido?",
              "Você tem um esquecimento que está piorando nos últimos meses?",
              "Você tem um esquecimento que está impedindo a realização de alguma atividade do cotidiano?",
            ].map((t, i) => yn("cognicao", i, t))}
            {[
              "No último mês, você ficou com desânimo, tristeza ou desesperança?",
              "No último mês, você perdeu o interesse ou prazer em atividades anteriormente prazerosas?",
            ].map((t, i) => yn("humor", i, t))}
          </>
        )}
        {etapa === 3 && (
          <>
            {[
              "Você é incapaz de elevar os braços acima do nível dos ombros?",
              "Você é incapaz de manusear ou segurar objetos?",
            ].map((t, i) => yn("alcance", i, t))}
            <div
              className={cn(
                "border-b border-border py-4 last:border-0",
                erro && r.aerobica === null && "text-destructive",
              )}
            >
              <p className="mb-2 font-medium">Você tem alguma das quatro condições abaixo?</p>
              <ul className="mb-3 list-disc pl-5 text-sm text-muted-foreground">
                <li>
                  Perda de peso não intencional de 4,5 kg ou 5% do peso corporal no último ano, ou 6
                  kg nos últimos 6 meses, ou 3 kg no último mês
                </li>
                <li>IMC (Índice de Massa Corporal) menor que 22 kg/m²</li>
                <li>Circunferência da panturrilha menor que 31 cm</li>
                <li>Velocidade da marcha: mais de 5 segundos para percorrer 4 metros</li>
              </ul>
              <div className="grid grid-cols-2 gap-2 sm:w-64">
                <Opcao ativo={r.aerobica === true} onClick={() => set({ aerobica: true })}>
                  Sim
                </Opcao>
                <Opcao ativo={r.aerobica === false} onClick={() => set({ aerobica: false })}>
                  Não
                </Opcao>
              </div>
            </div>
            {[
              "Você tem dificuldade para caminhar que possa impedir a realização de alguma atividade do cotidiano?",
              "Você teve duas ou mais quedas no último ano?",
              "Você perde urina ou fezes, sem querer, em algum momento?",
            ].map((t, i) => yn("marcha", i, t))}
            {[
              "Você tem problemas de visão capazes de impedir a realização de alguma atividade do cotidiano? (Permitido o uso de óculos ou lentes de contato)",
              "Você tem problemas de audição capazes de impedir a realização de alguma atividade do cotidiano? (Permitido o uso de aparelho de audição)",
            ].map((t, i) => yn("comunicacao", i, t))}
          </>
        )}
        {etapa === 4 && (
          <div className={cn(erro && r.comorbidades === null && "text-destructive")}>
            <p className="mb-2 font-medium">Você tem alguma das três condições abaixo?</p>
            <ul className="mb-3 list-disc pl-5 text-sm text-muted-foreground">
              <li>Cinco ou mais doenças crônicas</li>
              <li>Uso regular de cinco ou mais medicamentos diferentes, todos os dias</li>
              <li>Internação recente, nos últimos 6 meses</li>
            </ul>
            <div className="grid grid-cols-2 gap-2 sm:w-64">
              <Opcao ativo={r.comorbidades === true} onClick={() => set({ comorbidades: true })}>
                Sim
              </Opcao>
              <Opcao ativo={r.comorbidades === false} onClick={() => set({ comorbidades: false })}>
                Não
              </Opcao>
            </div>
          </div>
        )}
        {etapa === 5 && (
          <div className="space-y-8">
            {EASFI.map((d, i) => (
              <div key={d.t}>
                <p
                  className={cn(
                    "mb-3 font-medium",
                    erro && r.easfi[i] === null && "text-destructive",
                  )}
                >
                  {i + 1}. {d.t}
                </p>
                <div className="grid gap-2">
                  {d.o.map((o, j) => (
                    <Opcao
                      key={o}
                      ativo={r.easfi[i] === j + 1}
                      onClick={() => {
                        const a = [...r.easfi];
                        a[i] = j + 1;
                        set({ easfi: a });
                      }}
                    >
                      <span className="mr-2 font-semibold">{j + 1}</span>
                      {o}
                    </Opcao>
                  ))}
                </div>
              </div>
            ))}
          </div>
        )}
        {etapa === 6 && (
          <>
            <p className="text-sm text-muted-foreground">
              Qualquer marcador positivo reclassifica para Nível 3 ou 4.
            </p>
            {SENTINELAS.map((t, i) => yn("sentinelas", i, t))}
          </>
        )}
        {erro && !completa() && (
          <p
            role="alert"
            className="mt-4 rounded-md bg-destructive/10 p-3 text-sm text-destructive"
          >
            Responda todas as perguntas desta etapa para continuar.
          </p>
        )}
        {erroServidor && (
          <p
            role="alert"
            className="mt-4 rounded-md bg-destructive/10 p-3 text-sm text-destructive"
          >
            Não foi possível concluir a avaliação. Verifique a conexão com o servidor e tente
            novamente.
          </p>
        )}
      </Cartao>

      <div className="fixed inset-x-0 bottom-0 border-t border-border bg-card/95 backdrop-blur">
        <div className="mx-auto flex max-w-3xl gap-3 px-4 py-3">
          <Botao
            variante="secundario"
            className="flex-1"
            disabled={etapa === 0}
            onClick={() => {
              setErro(false);
              setEtapa(etapa - 1);
              window.scrollTo(0, 0);
            }}
          >
            Voltar
          </Botao>
          <Botao className="flex-[2]" disabled={concluir.isPending} onClick={avancar}>
            {concluir.isPending
              ? "Concluindo…"
              : etapa === ETAPAS.length - 1
                ? "Ver resultado"
                : "Próxima etapa"}
          </Botao>
        </div>
      </div>
    </Pagina>
  );
}
