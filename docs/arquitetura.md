# Arquitetura

Decisões arquiteturais de alto nível do monorepo **Protocolo Idoso APS**. Decisões pontuais e justificadas são registradas como ADRs em [decisoes/](decisoes/).

## Visão geral

Sistema web cliente-servidor com fronteira clara entre apresentação (frontend) e regras de negócio/persistência (backend):

```
┌───────────────────────┐          ┌─────────────────────┐         ┌──────────────────────────┐
│  apps/frontend        │   HTTP   │  apps/backend       │         │  Banco de dados          │
│  React + Vite + TS    │ ───────► │  Go (API JSON REST) │ ──────► │  PostgreSQL              │
│  Tailwind + shadcn/ui │          │  pgx + SQL explícito│         │  dev: local · prod:      │
│                       │          │                     │         │  Supabase (gerenciado)   │
└───────────────────────┘          └─────────────────────┘         └──────────────────────────┘
```

- **Regras de domínio centrais** (cálculo dos escores IVCF-20/EASFI-BR, matriz MIV-CFS, marcadores sentinela, reclassificação de nível, geração do resumo copiável) pertencem ao **backend**, em pacotes puros testáveis. O frontend coleta, exibe e transfere — nunca recalcula sozinho o que define prioridade assistencial.
- **Persistência**: o backend armazena **rascunhos** e **histórico de avaliações** — todos **sem identificação da pessoa idosa**, vinculados a um identificador anônimo de dispositivo ([ADR-0002](decisoes/ADR-0002-avaliacao-assistida-anonima.md)).

## Estrutura do monorepo

| Caminho | Papel |
|---|---|
| `apps/frontend/` | Aplicação web (formulário assistido, resultado, resumo copiável, histórico) |
| `apps/backend/` | API, regras de domínio (IVCF-20, EASFI-BR, MIV-CFS), persistência anônima |
| `docs/` | Documentação formal (visão, arquitetura, ADRs) |

## Stack tecnológica

| Camada | Stack | Registro |
|---|---|---|
| Frontend | React + Vite + TypeScript (estrito) + Tailwind CSS + shadcn/ui | [ADR-0003](decisoes/ADR-0003-stack-frontend.md) |
| Backend | Go + API JSON REST + pgx com SQL explícito | [ADR-0004](decisoes/ADR-0004-stack-backend.md) |
| Banco | PostgreSQL — dev local (container) · produção no Supabase (apenas Postgres gerenciado) | [ADR-0005](decisoes/ADR-0005-banco-postgresql-supabase.md) |
| Domínio | Regras MIV-CFS em pacotes Go puros, testáveis contra os casos do protocolo | ADR-0004 |

**Ponto aberto**: framework/roteador HTTP do backend (Gin/Echo/chi) — a decidir na implementação e registrar em ADR próprio.

Demais decisões institucionais:

| Tema | Decisão | Registro |
|---|---|---|
| Monorepo com `apps/frontend` + `apps/backend` | Sim | [ADR-0001](decisoes/ADR-0001-monorepo.md) |
| Natureza do sistema: avaliação assistida anônima com histórico | Sim — sem identificação do idoso, vínculo por dispositivo, histórico apagável | [ADR-0002](decisoes/ADR-0002-avaliacao-assistida-anonima.md) |
| Licença | Proprietária | [LICENSE](../LICENSE) |
| Idioma da documentação | Português (BR) | — |

## Ambientes

| Ambiente | Banco | Observações |
|---|---|---|
| Desenvolvimento | PostgreSQL local (recomendado: container Docker) | Migrations versionadas no repositório, idênticas às de produção |
| Produção | Supabase (PostgreSQL gerenciado via connection string) | Sem uso de Auth/Storage/Realtime do Supabase |

## Restrições e requisitos transversais

- **Anonimato por design**: nenhuma entidade persistida pode conter identificação de pessoa idosa (nome, nº de ficha, CNS, data de nascimento identificadora); validado em revisão e testes
- **Eliminação como cidadã de primeira classe**: apagar avaliações do histórico (individual e total) é funcionalidade nuclear, testada e sempre disponível na interface
- **Logs sem conteúdo clínico**: a infraestrutura não deve registrar respostas, escores ou resultados em logs de aplicação ou acesso
- **Reprodutibilidade das regras**: a matriz MIV-CFS e os marcadores sentinela são normas do protocolo — sua implementação deve ser testável e verificável contra os casos do documento-fonte (ver [visao-geral.md](visao-geral.md))
- **Rascunhos confiáveis**: salvar e retomar avaliação de onde parou deve ser comportamento verificável (estado consistente entre sessões)
