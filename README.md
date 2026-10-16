# Protocolo Idoso APS

Sistema web de apoio à **identificação, estratificação e manejo da vulnerabilidade clínico-funcional e sociofamiliar da pessoa idosa** na Atenção Primária à Saúde (APS), baseado no protocolo assistencial MIV-CFS.

> **Status: fase de setup** — este repositório está em estruturação. Ainda não há código de aplicação; a stack do frontend e do backend será definida e registrada em ADR (ver [docs/arquitetura.md](docs/arquitetura.md)).

## Sobre o projeto

O protocolo assistencial (elaborado pela equipe da Unoeste/Curso de Medicina, 2026) padroniza, no âmbito da Unidade Básica de Saúde (UBS), os processos de identificação, avaliação, estratificação, acompanhamento e manejo das pessoas idosas (≥ 60 anos) em situação de vulnerabilidade, por meio da combinação de dois instrumentos validados:

| Instrumento | Dimensão | Pontuação | Estratificação |
|---|---|---|---|
| **IVCF-20** | Clínico-funcional (8 domínios) | 0–40 | Baixo (0–6) · Moderado (7–14) · Alto (15–40) |
| **EASFI-BR** | Sociofamiliar (5 domínios) | 5–25 | Adequada (5–9) · Risco social (10–14) · Problema social (≥ 15) |

Os escores são cruzados na **Matriz Integrada de Vulnerabilidade (MIV-CFS)**, que classifica o idoso em **4 níveis de prioridade assistencial** — a condição social atua como modificador do nível de atenção, sem soma numérica direta. **Marcadores sociais sentinela** (violência, situação de rua, ausência de cuidador, entre outros) reclassificam diretamente para os níveis mais graves.

Detalhes completos em [docs/visao-geral.md](docs/visao-geral.md).

## Estrutura do monorepo

```
protocolo-idoso-aps/
├── README.md               # Este arquivo
├── LICENSE                 # Licença proprietária
├── CONTRIBUTING.md         # Como contribuir (commits, branches, PRs)
├── docs/                   # Documentação formal
│   ├── index.md            # Hub de navegação da documentação
│   ├── visao-geral.md      # Escopo do sistema baseado no protocolo
│   ├── arquitetura.md      # Decisões arquiteturais de alto nível
│   └── decisoes/           # Registro de decisões (ADRs)
└── apps/
    ├── frontend/           # Aplicação web (stack a definir)
    └── backend/            # API e regras de negócio (stack a definir)
```

## Documentação

| Documento | Quando ler |
|---|---|
| [docs/visao-geral.md](docs/visao-geral.md) | Entender o domínio, o escopo e os requisitos do sistema |
| [docs/arquitetura.md](docs/arquitetura.md) | Entender a estrutura do monorepo e decisões técnicas |
| [docs/decisoes/](docs/decisoes/) | Histórico de decisões arquiteturais (ADRs) |
| [CONTRIBUTING.md](CONTRIBUTING.md) | Contribuir com o projeto |

## Aviso de dados sensíveis

Este sistema trata **dados pessoais sensíveis de saúde** (dados de pessoas idosas). Todo desenvolvimento deve observar a **LGPD (Lei nº 13.709/2018)**, o Estatuto da Pessoa Idosa (Lei nº 10.741/2003) e as diretrizes do SUS aplicáveis. Veja [docs/visao-geral.md](docs/visao-geral.md#privacidade-e-lgpd).

## Contribuindo

Consulte [CONTRIBUTING.md](CONTRIBUTING.md). Commits seguem o padrão **Conventional Commits** e decisões técnicas relevantes devem ser registradas como ADRs em `docs/decisoes/`.

## Licença

Software **proprietário**. Todos os direitos reservados — consulte o arquivo [LICENSE](LICENSE).
