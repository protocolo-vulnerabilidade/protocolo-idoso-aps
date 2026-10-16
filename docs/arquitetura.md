# Arquitetura

Decisões arquiteturais de alto nível do monorepo **Protocolo Idoso APS**. Decisões pontuais e justificadas são registradas como ADRs em [decisoes/](decisoes/).

## Visão geral

Sistema web cliente-servidor com fronteira clara entre apresentação (frontend) e regras de negócio/persistência (backend):

```
┌──────────────┐         ┌──────────────┐         ┌────────────┐
│  frontend    │  HTTP   │   backend    │         │  banco de  │
│  (apps/      │ ──────► │   (apps/     │ ──────► │   dados    │
│   frontend)  │  API    │   backend)   │         └────────────┘
└──────────────┘         └──────────────┘
```

- **Regras de domínio centrais** (cálculo dos escores IVCF-20/EASFI-BR, matriz MIV-CFS, marcadores sentinela, reclassificação de nível, geração do resumo copiável) pertencem ao **backend**. O frontend coleta, exibe e transfere — nunca recalcula sozinho o que define prioridade assistencial.
- **Persistência**: o backend armazena **rascunhos** e **histórico de avaliações** — todos **sem identificação da pessoa idosa**, vinculados a um identificador anônimo de dispositivo ([ADR-0002](decisoes/ADR-0002-avaliacao-assistida-anonima.md)).

## Estrutura do monorepo

| Caminho | Papel |
|---|---|
| `apps/frontend/` | Aplicação web (formulário assistido, resultado, resumo copiável, histórico) |
| `apps/backend/` | API, regras de domínio (IVCF-20, EASFI-BR, MIV-CFS), persistência anônima |
| `docs/` | Documentação formal (visão, arquitetura, ADRs) |

## Stack tecnológica

**A definir.** A escolha de linguagens, frameworks, banco de dados e ferramentas de cada app será registrada em ADRs dedicados (p.ex. `ADR-0003-stack-frontend.md`, `ADR-0004-stack-backend.md`), contemplando: maturidade, manutenabilidade, adequação ao contexto de saúde e produtividade de uma equipe pequena.

Decisões já tomadas:

| Tema | Decisão | Registro |
|---|---|---|
| Monorepo com `apps/frontend` + `apps/backend` | Sim | [ADR-0001](decisoes/ADR-0001-monorepo.md) |
| Natureza do sistema: avaliação assistida anônima com histórico | Sim — sem identificação do idoso, vínculo por dispositivo, histórico apagável | [ADR-0002](decisoes/ADR-0002-avaliacao-assistida-anonima.md) |
| Licença | Proprietária | [LICENSE](../LICENSE) |
| Idioma da documentação | Português (BR) | — |
| Stack de frontend/backend e banco de dados | Pendente | ADR futuro |

## Restrições e requisitos transversais

- **Anonimato por design**: nenhuma entidade persistida pode conter identificação de pessoa idosa (nome, nº de ficha, CNS, data de nascimento identificadora); validado em revisão e testes
- **Eliminação como cidadã de primeira classe**: apagar avaliações do histórico (individual e total) é funcionalidade nuclear, testada e sempre disponível na interface
- **Logs sem conteúdo clínico**: a infraestrutura não deve registrar respostas, escores ou resultados em logs de aplicação ou acesso
- **Reprodutibilidade das regras**: a matriz MIV-CFS e os marcadores sentinela são normas do protocolo — sua implementação deve ser testável e verificável contra os casos do documento-fonte (ver [visao-geral.md](visao-geral.md))
- **Rascunhos confiáveis**: salvar e retomar avaliação de onde parou deve ser comportamento verificável (estado consistente entre sessões)
