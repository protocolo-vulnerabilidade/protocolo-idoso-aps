# ADR-0005 — Banco de dados: PostgreSQL (dev local) com deploy no Supabase

- **Status**: Aceita
- **Data**: 2026-10-16

## Contexto

O sistema persiste **rascunhos** e **histórico de avaliações** — todos anônimos, vinculados a um identificador de dispositivo ([ADR-0002](ADR-0002-avaliacao-assistida-anonima.md)). O desenvolvimento acontecerá localmente na máquina do desenvolvedor e, em seguida, o sistema será publicado.

## Decisão

- **PostgreSQL** como banco de dados único, nos dois ambientes:
  - **Desenvolvimento**: instância PostgreSQL local (recomenda-se container Docker com `docker compose`)
  - **Produção/deploy**: **Supabase**, utilizado **apenas como PostgreSQL gerenciado**, acessado via connection string padrão (`postgresql://...`)
- **Sem uso dos serviços adicionais do Supabase** (Auth, Storage, Realtime, Edge Functions): a identificação por dispositivo anônimo já está decidida no ADR-0002 e não há arquivos a armazenar
- **Migrations idênticas** nos dois ambientes, versionadas no repositório (via pgx/SQL explícito — [ADR-0004](ADR-0004-stack-backend.md))

## Consequências

**Positivas**
- **Zero lock-in funcional**: como se usa apenas o Postgres padrão, trocar o Supabase por outro provedor (ou voltar a self-hosted) é apenas mudar a connection string
- Paridade dev/prod: mesmo banco nos dois ambientes reduz surpresas de comportamento
- Supabase remove o fardo operacional de backups/disponibilidade em produção

**Negativas / trade-offs**
- Funcionalidades do Supabase ficam deliberadamente não usadas (custo/benefício aceito pela portabilidade)
- Desenvolvimento depende de instância Postgres local (mitigado por container Docker)

## Alternativas consideradas

1. **SQLite local + Postgres em produção** — rejeitada: quebra de paridade entre ambientes.
2. **Firebase/Supabase como backend completo (BaaS)** — rejeitada: tiraria as regras de domínio do backend Go e criaria lock-in.
3. **PostgreSQL self-hosted em produção** — rejeitado por ora: custo operacional desnecessário para o escopo; revisível via novo ADR.
