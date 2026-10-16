# apps/backend

API e núcleo de domínio do **Protocolo Idoso APS** — regras de estratificação, persistência anônima e serviços para o frontend.

## Papel

- **Regras de domínio centrais**: cálculo e classificação do IVCF-20 (0–40) e EASFI-BR (5–25), cruzamento na **matriz MIV-CFS** (níveis 1–4), **marcadores sentinela** com reclassificação direta, conduta/prazo por nível e geração do **resumo copiável**
- **Persistência anônima**: rascunhos e histórico de avaliações **sem qualquer identificação da pessoa idosa**, vinculados a um identificador anônimo de dispositivo ([ADR-0002](../../docs/decisoes/ADR-0002-avaliacao-assistida-anonima.md))
- **Eliminação**: operações de exclusão individual e total do histórico como parte nuclear da API
- **Logs sem conteúdo clínico**: a infraestrutura não registra respostas, escores ou resultados

## Stack

**Definida em [ADR-0004](../../docs/decisoes/ADR-0004-stack-backend.md)**: Go + API JSON REST + pgx com SQL explícito (sem ORM). Framework HTTP a decidir na implementação.

**Banco ([ADR-0005](../../docs/decisoes/ADR-0005-banco-postgresql-supabase.md))**: PostgreSQL — desenvolvimento local (container) e produção no Supabase (apenas Postgres gerenciado, via connection string).

## Status

**Pendente de implementação.**

## Documentação

- [Visão geral do sistema](../../docs/visao-geral.md)
- [Arquitetura](../../docs/arquitetura.md)
- [ADR-0002 — Avaliação assistida anônima](../../docs/decisoes/ADR-0002-avaliacao-assistida-anonima.md)
- [ADR-0004 — Stack do backend](../../docs/decisoes/ADR-0004-stack-backend.md)
- [ADR-0005 — Banco PostgreSQL/Supabase](../../docs/decisoes/ADR-0005-banco-postgresql-supabase.md)
