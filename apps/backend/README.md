# apps/backend

API e núcleo de domínio do **Protocolo Idoso APS** — regras de estratificação, persistência anônima e serviços para o frontend.

## Papel

- **Regras de domínio centrais**: cálculo e classificação do IVCF-20 (0–40) e EASFI-BR (5–25), cruzamento na **matriz MIV-CFS** (níveis 1–4), **marcadores sentinela** com reclassificação direta, conduta/prazo por nível e geração do **resumo copiável** — em `internal/dominio`, pacote puro testável
- **Persistência anônima**: rascunhos e histórico de avaliações **sem qualquer identificação da pessoa idosa**, vinculados a um identificador anônimo de dispositivo ([ADR-0002](../../docs/decisoes/ADR-0002-avaliacao-assistida-anonima.md))
- **Eliminação**: operações de exclusão individual e total do histórico como parte nuclear da API
- **Logs sem conteúdo clínico**: a infraestrutura registra apenas método/caminho/status/duração

## Stack

**[ADR-0004](../../docs/decisoes/ADR-0004-stack-backend.md)**: Go + API JSON REST + pgx com SQL explícito (sem ORM). **[ADR-0006](../../docs/decisoes/ADR-0006-roteador-http-stdlib.md)**: roteador `net/http` da stdlib.

**Banco ([ADR-0005](../../docs/decisoes/ADR-0005-banco-postgresql-supabase.md))**: PostgreSQL — desenvolvimento local e produção no Supabase (apenas Postgres gerenciado, via connection string).

## Como rodar (desenvolvimento)

```powershell
# banco: PostgreSQL local em localhost:5432 (ajuste DATABASE_URL se necessário)
$env:DATABASE_URL = 'postgresql://piaps_dev:piaps_dev@localhost:5432/protocolo_idoso_dev'
go run ./cmd/api   # aplica migrações e sobe em http://localhost:8080
```

Variáveis de ambiente: `DATABASE_URL`, `PORTA` (padrão `8080`), `CORS_ORIGENS` (padrão `http://localhost:5173,http://localhost:4173`).

Migrations ficam em `internal/armazenamento/migracoes/`, são embutidas no binário (`embed.FS`) e aplicadas no startup com advisory lock.

## API

Dispositivo anônimo: cabeçalho `X-Dispositivo-Id` (UUID gerado pelo frontend) em todas as rotas exceto `/api/saude`.

| Método e caminho | Função |
|---|---|
| `GET /api/saude` | verificação de saúde |
| `GET /api/rascunhos` | listar rascunhos do dispositivo |
| `POST /api/rascunhos` | criar rascunho (respostas vazias por padrão) |
| `GET /api/rascunhos/{id}` | obter rascunho |
| `PUT /api/rascunhos/{id}` | auto-save (`{etapa, respostas}`) |
| `DELETE /api/rascunhos/{id}` | descartar rascunho |
| `POST /api/rascunhos/{id}/concluir` | valida completude, **calcula o resultado no servidor** e move para o histórico (mesmo id) |
| `GET /api/historico` | listar avaliações concluídas |
| `GET /api/historico/{id}` | detalhe com resultado e **resumo copiável** |
| `DELETE /api/historico/{id}` | apagar avaliação concluída |
| `POST /api/historico/{id}/reabrir` | reabrir concluída como rascunho (mesmo id, etapa 0) |
| `DELETE /api/historico` | apagar todo o histórico do dispositivo |

Erros: `400` JSON inválido/campos desconhecidos, `404` não encontrado (ou de outro dispositivo), `409` estado conflitante, `422` validação com `detalhes [{campo, mensagem}]`.

## Testes

```powershell
go test ./...
```

`internal/dominio` tem testes de mesa das regras (fronteiras IVCF 6/7 e 14/15, EASFI 9/10 e 14/15, teto AIVD = 4, máximo = 40, matriz completa, sentinelas); `internal/httpapi` testa os handlers com armazenamento em memória (fluxo completo, isolamento entre dispositivos, eliminação, CORS).

## Documentação

- [Visão geral do sistema](../../docs/visao-geral.md)
- [Arquitetura](../../docs/arquitetura.md)
- [ADR-0002 — Avaliação assistida anônima](../../docs/decisoes/ADR-0002-avaliacao-assistida-anonima.md)
- [ADR-0004 — Stack do backend](../../docs/decisoes/ADR-0004-stack-backend.md)
- [ADR-0005 — Banco PostgreSQL/Supabase](../../docs/decisoes/ADR-0005-banco-postgresql-supabase.md)
- [ADR-0006 — Roteador HTTP stdlib](../../docs/decisoes/ADR-0006-roteador-http-stdlib.md)
