# ADR-0003 — Stack do frontend: React + Vite + TypeScript + Tailwind CSS + shadcn/ui

- **Status**: Aceita
- **Data**: 2026-10-16

## Contexto

O frontend é a única interface do sistema: um **formulário assistido** com dados nominais e numéricos (domínios do IVCF-20 e EASFI-BR), verificação de marcadores sentinela, exibição do resultado com conduta, resumo copiável, rascunhos persistentes e histórico apagável. É usado por profissionais de saúde em computadores da UBS — a interface precisa ser clara, acessível e responsiva, e será mantida por uma equipe pequena.

## Decisão

- **React** (com componentes funcionais e hooks) como biblioteca de UI
- **Vite** como bundler e servidor de desenvolvimento
- **TypeScript** em modo estrito para contratos de dados com o backend
- **Tailwind CSS** para estilização utilitária
- **shadcn/ui** (componentes sobre Radix UI) para os blocos de interface (formulários, diálogos, tabelas, botões)

O frontend consome a API do backend e **não implementa regras de estratificação** (ver [arquitetura](../arquitetura.md)).

## Consequências

**Positivas**
- shadcn/ui entrega componentes **acessíveis** (Radix) prontos — crítico para uso clínico e formulários longos
- Tailwind + Vite: iteração rápida de UI sem CSS ad-hoc
- TypeScript dá segurança nos contratos de API (escores, classificações, níveis)
- Ecossistema enorme e material de apoio abundante

**Negativas / trade-offs**
- Dependência do ecossistema Radix/Tailwind (evolução ativa, breaking changes ocasionais)
- shadcn/ui copia código para o repositório (mais código próprio para manter, em troca de controle total)

## Alternativas consideradas

1. **Angular** — rejeitado: curva de aprendizado e verbosidade desproporcionais ao escopo.
2. **Vue/Nuxt** — rejeitado: ecossistema menor para componentes acessíveis prontos no contexto da equipe.
3. **Next.js (SSR)** — rejeitado: não há necessidade de SSR/SEO para ferramenta de uso interno; SPA simplifica.
