import * as React from "react";
import { OTPInputContext } from "input-otp";

import {
  AlertDialog,
  AlertDialogCancel,
  AlertDialogContent,
  AlertDialogDescription,
  AlertDialogFooter,
  AlertDialogHeader,
  AlertDialogTitle,
} from "@/components/ui/alert-dialog";
import { InputOTP, InputOTPGroup } from "@/components/ui/input-otp";
import { Botao } from "@/components/ui-base";
import { cn } from "@/lib/utils";

function SlotMascarada({ index }: { index: number }) {
  const { slots } = React.useContext(OTPInputContext);
  const { char, hasFakeCaret, isActive } = slots[index] ?? {
    char: null,
    hasFakeCaret: false,
    isActive: false,
  };
  return (
    <div
      className={cn(
        "relative flex h-9 w-9 items-center justify-center border-y border-r border-input text-sm shadow-sm transition-all first:rounded-l-md first:border-l last:rounded-r-md",
        isActive && "z-10 ring-1 ring-ring",
      )}
    >
      {char ? "•" : null}
      {hasFakeCaret && (
        <div className="pointer-events-none absolute inset-0 flex items-center justify-center">
          <div className="h-4 w-px animate-caret-blink bg-foreground duration-1000" />
        </div>
      )}
    </div>
  );
}

interface DialogoSenhaExclusaoProps {
  aberto: boolean;
  titulo: string;
  descricao: string;
  pendente: boolean;
  erro: string | null;
  onConfirmar: (senha: string) => void;
  onFechar: () => void;
}

export function DialogoSenhaExclusao({
  aberto,
  titulo,
  descricao,
  pendente,
  erro,
  onConfirmar,
  onFechar,
}: DialogoSenhaExclusaoProps) {
  const [senha, setSenha] = React.useState("");

  React.useEffect(() => {
    if (!aberto) setSenha("");
  }, [aberto]);

  return (
    <AlertDialog
      open={aberto}
      onOpenChange={(abertoAtual) => {
        if (!abertoAtual && !pendente) onFechar();
      }}
    >
      <AlertDialogContent className="max-w-sm">
        <AlertDialogHeader>
          <AlertDialogTitle>{titulo}</AlertDialogTitle>
          <AlertDialogDescription>{descricao}</AlertDialogDescription>
        </AlertDialogHeader>
        <div className="flex flex-col items-center gap-2">
          <InputOTP
            maxLength={4}
            value={senha}
            onChange={(valor) => setSenha(valor.replace(/\D/g, "").slice(0, 4))}
            pattern="^[0-9]*$"
            autoComplete="off"
            inputMode="numeric"
            aria-label="Senha de confirmação"
            disabled={pendente}
          >
            <InputOTPGroup>
              <SlotMascarada index={0} />
              <SlotMascarada index={1} />
              <SlotMascarada index={2} />
              <SlotMascarada index={3} />
            </InputOTPGroup>
          </InputOTP>
          <p className="text-xs text-muted-foreground">
            Digite a senha de confirmação (4 dígitos).
          </p>
          {erro && (
            <p role="alert" className="text-sm font-medium text-destructive">
              {erro}
            </p>
          )}
        </div>
        <AlertDialogFooter>
          <AlertDialogCancel disabled={pendente}>Cancelar</AlertDialogCancel>
          <Botao
            variante="perigo"
            disabled={senha.length < 4 || pendente}
            onClick={() => onConfirmar(senha)}
          >
            {pendente ? "Apagando…" : "Apagar"}
          </Botao>
        </AlertDialogFooter>
      </AlertDialogContent>
    </AlertDialog>
  );
}
