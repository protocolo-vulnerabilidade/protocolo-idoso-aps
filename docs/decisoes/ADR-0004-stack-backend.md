# ADR-0004 — Stack do backend: Go + API JSON REST + pgx/SQL explícito

- **Status**: Aceita (com ponto aberto)
- **Data**: 2026-10-16

## Contexto

O backend concentra as **regras de domínio** (cálculo e classificação do IVCF-20 e EASFI-BR, matriz MIV-CFS, marcadores sentinela, reclassificação de nível, geração do resumo copiável) e a **persistência anônima** de rascunhos e histórico por dispositivo ([ADR-0002](ADR-0002-avaliacao-assistida-anonima.md)).

As regras do protocolo são **normativas** — calculam prioridade assistencial — e devem ser testáveis e verificáveis contra os casos do documento-fonte. A equipe é pequena e o deploy deve ser simples.

## Decisão

- **Go** como linguagem do backend
- **API JSON REST** como contrato com o frontend
- **pgx + SQL explícito** como camada de acesso ao PostgreSQL (sem ORM) — queries e migrations versionadas no repositório
- Regras de domínio em **pacotes puros** (sem dependência de HTTP ou banco), testáveis isoladamente

### Ponto aberto

O **framework/roteador HTTP** (p.ex. Gin, Echo, chi + stdlib) será decidido na implementação e registrado em ADR próprio quando fechado. Este ADR fixa apenas a linguagem, o estilo de API e a camada de dados.

## Consequências

**Positivas**
- Tipagem forte e testes rápidos para as regras da matriz MIV-CFS e sentinelas
- Binário único e imagem de deploy minimalista; consumo de recursos baixo
- SQL explícito com pgx: controle total das queries, sem mágica de ORM sobre dados de saúde
- Regras de domínio puras facilitam testes de mesa contra os casos do protocolo

**Negativas / trade-offs**
- Sem ORM: migrations e SQL ficam sob responsabilidade direta da equipe (disciplina exigida)
- Backend em linguagem distinta do frontend (Go ≠ TypeScript) — contratos de API precisam ser mantidos manualmente

## Alternativas consideradas

1. **Node.js/TypeScript (monolinguagem)** — rejeitada: perderia a robustez de tipagem em runtime e o deploy trivial do binário Go.
2. **Python/FastAPI** — rejeitada: desempenho e empacotamento de deploy menos atrativos para este escopo.
3. **GORM (ORM)** — rejeitado: abstração desnecessária para um schema pequeno; SQL explícito favorece auditoria das queries.
