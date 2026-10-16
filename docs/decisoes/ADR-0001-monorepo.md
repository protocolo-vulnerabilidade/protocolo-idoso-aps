# ADR-0001 — Monorepo com apps de frontend e backend

- **Status**: Aceita
- **Data**: 2026-10-16

## Contexto

O sistema Protocolo Idoso APS é um sistema web de escopo contido (uma equipe pequena, uma linha de produto), composto por uma aplicação de interface (uso pela equipe da UBS) e um serviço de backend com as regras de domínio (IVCF-20, EASFI-BR, matriz MIV-CFS, marcadores sentinela) e persistência de dados sensíveis.

É necessário decidir como organizar o código-fonte: repositórios separados por aplicação ou um único repositório (monorepo).

## Decisão

Adotar **monorepo** com uma pasta `apps/` contendo os dois aplicativos:

```
apps/
├── frontend/   # interface web
└── backend/    # API, regras de domínio, persistência
```

- Branch padrão única (`main`); mudanças de front e back podem ser versionadas e revisadas em conjunto quando couber.
- A fronteira entre os apps é a **API do backend**: o frontend não implementa regras de estratificação de vulnerabilidade.
- Documentação formal centralizada em `docs/` (não dentro dos apps).

## Consequências

**Positivas**
- Uma única fonte da verdade para docs, issues e histórico; visão integrada do produto
- Refatorações que atravessam front/back em um mesmo PR/commit
- Menos sobrecarga de setup e permissões para uma equipe pequena

**Negativas / trade-offs**
- Ferramentas de CI/CD precisarão de filtros de caminho se os apps usarem stacks distintas
- O repositório cresce com artefatos dos dois apps (aceitável no escopo atual)

## Alternativas consideradas

1. **Multirepo (um repositório por app)** — rejeitada: overhead de coordenação, sincronização de versões de API e documentação duplicada para um projeto de escopo contido.
2. **App único fullstack (sem separação)** — rejeitada por ora: dados sensíveis de saúde e regras normativas do protocolo pedem uma fronteira explícita de API e testes de domínio isolados no backend; a decisão pode ser revista em novo ADR.

## Referências

- [docs/arquitetura.md](../arquitetura.md)
- Boas práticas observadas no monorepo de referência `CEIA-Ag.IA/ai-4sdlc-platform` (pnpm workspaces + Turborepo, fronteiras mecânicas entre pacotes)
