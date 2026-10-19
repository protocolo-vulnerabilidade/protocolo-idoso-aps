# ADR-0006 — Roteador HTTP: `net/http` da biblioteca padrão

- **Status**: Aceita
- **Data**: 2026-10-19
- **Relação**: Fecha o ponto aberto do [ADR-0004](ADR-0004-stack-backend.md) (framework/roteador a decidir na implementação)

## Contexto

O [ADR-0004](ADR-0004-stack-backend.md) fixou Go + API JSON REST + pgx/SQL explícito, mas deixou em aberto o framework/roteador HTTP (Gin, Echo, chi ou stdlib). A API é pequena e de superfície estável: dois recursos (`rascunhos`, `historico`), rota de saúde, um middleware de identificação por dispositivo e CORS para desenvolvimento.

## Decisão

**`net/http` da biblioteca padrão**, usando o roteamento por padrões disponível desde o Go 1.22 (`mux.HandleFunc("POST /api/rascunhos/{id}/concluir", …)`, `r.PathValue("id")`).

Sem dependências de framework. Middlewares (log sem conteúdo clínico, CORS, dispositivo) são funções wrapper comuns.

## Consequências

**Positivas**

- Zero dependências além de `pgx` e `google/uuid`: binário único enxuto, superfície de atualização de segurança mínima
- Roteamento por padrões da stdlib cobre todo o necessário (métodos exatos, wildcards de caminho, precedência determinística)
- Compatível com `http.Handler`/`httptest`, facilitando testes de handlers sem servidor

**Negativas / trade-offs**

- Sem geradores de OpenAPI, binding sofisticado ou middlewares de ecossistema — aceito: o contrato JSON é pequeno e mantido manualmente (trade-off já registrado no ADR-0004)
- Recursos como rate limiting ou métricas exigiriam código próprio ou biblioteca adicional se um dia forem necessários

## Alternativas consideradas

1. **chi** — compatível com stdlib e leve, mas não agrega nada que os padrões do Go 1.22+ já não cubram neste escopo.
2. **Gin / Echo** — frameworks completos com binding/validação próprios; peso e opinão de design desproporcionais para ~12 rotas.
3. **gRPC / tRPC** — rejeitados de partida: o ADR-0004 fixa API JSON REST como contrato com o frontend.
