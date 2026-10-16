# Protocolo Idoso APS

Ferramenta web de **avaliação assistida da vulnerabilidade clínico-funcional e sociofamiliar da pessoa idosa** na Atenção Primária à Saúde (APS): aplicação dos instrumentos **IVCF-20** e **EASFI-BR** com cálculo automático da **Matriz Integrada de Vulnerabilidade (MIV-CFS)**.

> **Status: fase de setup** — documentação e decisões arquiteturais concluídas; implementação a iniciar. Stacks definidas em ADRs: frontend **React + Vite + TypeScript + Tailwind + shadcn/ui** · backend **Go + pgx** · banco **PostgreSQL** (dev local / Supabase em produção). Ver [docs/arquitetura.md](docs/arquitetura.md).

## Sobre o projeto

O sistema apoia a equipe da UBS na estratificação da vulnerabilidade da pessoa idosa (≥ 60 anos), a partir do protocolo assistencial MIV-CFS (Unoeste/Curso de Medicina, 2026), que combina dois instrumentos validados:

| Instrumento | Dimensão | Pontuação | Estratificação |
|---|---|---|---|
| **IVCF-20** | Clínico-funcional (8 domínios) | 0–40 | Baixo (0–6) · Moderado (7–14) · Alto (15–40) |
| **EASFI-BR** | Sociofamiliar (5 domínios) | 5–25 | Adequada (5–9) · Risco social (10–14) · Problema social (≥ 15) |

Os escores são cruzados na **matriz MIV-CFS**, que classifica o idoso em **4 níveis de prioridade assistencial**; **marcadores sociais sentinela** reclassificam diretamente para os níveis mais graves.

### Como o sistema é usado

1. O profissional **retoma um rascunho** ou inicia uma nova avaliação
2. Preenche o **formulário assistido** com os dados nominais e numéricos dos instrumentos
3. O sistema calcula os escores, cruza a **matriz MIV-CFS**, avalia os **marcadores sentinela** e apresenta o **nível de vulnerabilidade** com a conduta e o prazo de reavaliação recomendados
4. O profissional **copia o resumo do resultado** para a ficha do idoso **no sistema oficial** da UBS
5. O cálculo fica registrado no **histórico anônimo** do dispositivo — que pode ser consultado ou **apagado** a qualquer momento

### Características centrais

- **Anônimo por design** — nenhum dado de identificação da pessoa idosa é informado ou armazenado
- **Rascunhos persistentes** — é possível sair da tela e voltar exatamente de onde parou
- **Histórico gerenciável** — cálculos anteriores ficam guardados e podem ser abertos ou apagados (individualmente ou todos)
- **Resumo copiável** — saída formatada para transferir o resultado à ficha no sistema oficial

## Estrutura do monorepo

```
protocolo-idoso-aps/
├── README.md               # Este arquivo
├── LICENSE                 # Licença proprietária
├── CONTRIBUTING.md         # Como contribuir (commits, branches, PRs)
├── docs/                   # Documentação formal
│   ├── index.md            # Hub de navegação da documentação
│   ├── visao-geral.md      # Escopo do sistema baseado no protocolo
│   ├── arquitetura.md      # Decisões arquiteturais e stacks
│   └── decisoes/           # Registro de decisões (ADRs 0001–0005)
└── apps/
    ├── frontend/           # React + Vite + TypeScript + Tailwind + shadcn/ui
    └── backend/            # Go + pgx · API JSON REST
```

## Documentação

| Documento | Quando ler |
|---|---|
| [docs/visao-geral.md](docs/visao-geral.md) | Entender o domínio, o escopo e os requisitos do sistema |
| [docs/arquitetura.md](docs/arquitetura.md) | Entender a estrutura do monorepo, stacks e decisões técnicas |
| [docs/decisoes/](docs/decisoes/) | Histórico de decisões arquiteturais (ADRs) |
| [CONTRIBUTING.md](CONTRIBUTING.md) | Contribuir com o projeto |

## Privacidade e dados

- O sistema é **anônimo por design**: nenhum campo de identificação da pessoa idosa (nome, nº de ficha, etc.) é coletado — a associação entre o resultado e a pessoa idosa permanece no sistema oficial da UBS
- As respostas e resultados são armazenados **sem identificação**, vinculados a um **identificador de dispositivo** (sem contas de usuário)
- O histórico é **apagável** pelo próprio usuário (direito de eliminação)
- Consulte [docs/visao-geral.md](docs/visao-geral.md#privacidade-e-lgpd) para o tratamento completo de LGPD

## Contribuindo

Consulte [CONTRIBUTING.md](CONTRIBUTING.md). Commits seguem o padrão **Conventional Commits** e decisões técnicas relevantes devem ser registradas como ADRs em `docs/decisoes/`.

## Licença

Software **proprietário**. Todos os direitos reservados — consulte o arquivo [LICENSE](LICENSE).
