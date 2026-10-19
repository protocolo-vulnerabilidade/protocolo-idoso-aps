# apps/frontend

Interface web do **Protocolo Idoso APS** — formulário assistido de avaliação da vulnerabilidade da pessoa idosa para uso pela equipe da UBS (ACS, enfermagem, medicina, eMulti e coordenação).

## Papel

- **Formulário assistido** com os dados nominais e numéricos dos instrumentos IVCF-20 e EASFI-BR — sem nenhum campo de identificação da pessoa idosa
- Exibição do **resultado**: escores, classificações, nível MIV-CFS (1–4), marcadores sentinela e conduta/prazo de reavaliação recomendados
- **Resumo copiável** formatado para transferência à ficha no sistema oficial da UBS
- **Rascunhos e histórico**: retomar avaliação de onde parou; listar, reabrir e **apagar** avaliações
- Mantém o **identificador anônimo de dispositivo** no navegador (sem contas de usuário)
- Consome a **API do backend** — não implementa regras de estratificação

## Stack

**Definida em [ADR-0003](../../docs/decisoes/ADR-0003-stack-frontend.md)**: React + Vite + TypeScript (estrito) + Tailwind CSS + shadcn/ui.

## Status

**Implementado.** Consome a API REST do backend (`apps/backend`) — o cálculo dos escores, da matriz MIV-CFS, das sentinelas e a geração do resumo copiável acontecem no servidor.

- Cliente HTTP tipado em `src/lib/api.ts` (contrato mantido manualmente com o backend — ADR-0004); dados via **react-query** em `src/lib/armazenamento.ts`
- O `localStorage` guarda apenas o identificador anônimo de dispositivo (`piaps:anon`), enviado no cabeçalho `X-Dispositivo-Id`
- Auto-save do rascunho com debounce (~800 ms) durante o preenchimento

## Como rodar (desenvolvimento)

```powershell
bun install
bun run dev   # http://localhost:5173 — espera o backend em http://localhost:8080
```

A URL base da API pode ser sobrescrita com a variável `VITE_API_URL`.

```powershell
bun run test  # vitest
bun run lint  # eslint
```

## Documentação

- [Visão geral do sistema](../../docs/visao-geral.md)
- [Arquitetura](../../docs/arquitetura.md)
- [ADR-0002 — Avaliação assistida anônima](../../docs/decisoes/ADR-0002-avaliacao-assistida-anonima.md)
- [ADR-0003 — Stack do frontend](../../docs/decisoes/ADR-0003-stack-frontend.md)
