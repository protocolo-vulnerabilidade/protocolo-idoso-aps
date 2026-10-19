import { Link } from "@tanstack/react-router";
import type { ButtonHTMLAttributes, ReactNode } from "react";
import { cn } from "@/lib/utils";
import type { Nivel } from "@/lib/protocolo";

type Variante = "primario" | "secundario" | "fantasma" | "perigo";
const estilos: Record<Variante, string> = {
  primario: "bg-primary text-primary-foreground hover:bg-primary/90",
  secundario: "border border-input bg-card text-foreground hover:bg-accent",
  fantasma: "text-foreground hover:bg-accent",
  perigo: "border border-destructive/40 bg-card text-destructive hover:bg-destructive/10",
};

export function Botao({
  variante = "primario",
  className,
  ...p
}: ButtonHTMLAttributes<HTMLButtonElement> & { variante?: Variante }) {
  return (
    <button
      {...p}
      className={cn(
        "inline-flex min-h-12 items-center justify-center gap-2 rounded-md px-5 text-base font-semibold transition-colors focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring disabled:opacity-50",
        estilos[variante],
        className,
      )}
    />
  );
}

export function Cabecalho() {
  return (
    <header className="border-b border-border bg-card">
      <div className="mx-auto flex max-w-3xl items-center justify-between px-4 py-3">
        <Link to="/" className="flex items-center gap-2">
          <span className="grid h-8 w-8 place-items-center rounded-md bg-primary text-sm font-bold text-primary-foreground">
            +
          </span>
          <span className="font-semibold tracking-tight">Protocolo Idoso APS</span>
        </Link>
        <Link to="/historico" className="text-sm font-medium text-primary hover:underline">
          Histórico
        </Link>
      </div>
    </header>
  );
}

export function Pagina({ children }: { children: ReactNode }) {
  return (
    <div className="min-h-screen bg-background">
      <Cabecalho />
      <main className="mx-auto max-w-3xl px-4 py-6 pb-24">{children}</main>
    </div>
  );
}

export function Cartao({ children, className }: { children: ReactNode; className?: string }) {
  return (
    <section className={cn("rounded-lg border border-border bg-card p-5 shadow-sm", className)}>
      {children}
    </section>
  );
}

export const corNivel: Record<Nivel, string> = {
  1: "bg-nivel-1 text-nivel-foreground",
  2: "bg-nivel-2 text-nivel-foreground",
  3: "bg-nivel-3 text-nivel-foreground",
  4: "bg-nivel-4 text-nivel-foreground",
};
