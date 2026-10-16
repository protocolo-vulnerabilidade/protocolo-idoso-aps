# Arquitetura

Decisões arquiteturais de alto nível do monorepo **Protocolo Idoso APS**. Decisões pontuais e justificadas são registradas como ADRs em [decisoes/](decisoes/).

## Visão geral

Sistema web cliente-servidor com fronteira clara entre apresentação (frontend) e regras de negócio/persistência (backend):

```
┌──────────────┐         ┌──────────────┐
│  frontend    │  HTTP   │   backend    │
│  (apps/      │ ──────► │   (apps/     │
│   frontend)  │  API    │   backend)   │
└──────────────┘         └──────────────┘
```

**Regras de domínio centrais** (cálculo de escores IVCF-20/EASFI-BR, matriz MIV-CFS, marcadores sentinela, reclassificação de nível) pertencem ao **backend**. O frontend consome e apresenta — nunca recalcula sozinho o que define prioridade assistencial.

## Estrutura do monorepo

| Caminho | Papel |
|---|---|
| `apps/frontend/` | Aplicação web (interface de uso pela equipe da UBS) |
| `apps/backend/` | API, regras de domínio (IVCF-20, EASFI-BR, MIV-CFS), persistência |
| `docs/` | Documentação formal (visão, arquitetura, ADRs) |

## Stack tecnológica

**A definir.** A escolha de linguagens, frameworks e ferramentas de cada app será registrada em ADRs dedicados (p.ex. `ADR-0002-stack-frontend.md`, `ADR-0003-stack-backend.md`), contemplando: maturidade, manutenibilidade, adequação a dados sensíveis de saúde e produtividade de uma equipe pequena.

Decisões já tomadas:

| Tema | Decisão | Registro |
|---|---|---|
| Monorepo com `apps/frontend` + `apps/backend` | Sim | [ADR-0001](decisoes/ADR-0001-monorepo.md) |
| Licença | Proprietária | [LICENSE](../LICENSE) |
| Idioma da documentação | Português (BR) | — |
| Stack de frontend/backend | Pendente | ADR futuro |

## Restrições e requisitos transversais

- **LGPD**: dados sensíveis de saúde — acesso por papel profissional, trilha de auditoria, sem dados reais fora de produção
- **Auditoria**: mudanças de classificação de vulnerabilidade (especialmente reclassificações) devem ser rastreáveis (quem, quando, por qual critério)
- **Reprodutibilidade das regras**: a matriz MIV-CFS e os marcadores sentinela são normas do protocolo — sua implementação deve ser testável e verificável contra os casos do documento-fonte
