# Contribuindo com o Protocolo Idoso APS

Obrigado pelo interesse em contribuir. Este documento define o processo de trabalho do repositório.

## Visão geral do processo

1. Crie uma branch a partir de `main` seguindo a nomenclatura abaixo
2. Faça commits seguindo **Conventional Commits**
3. Abra um Pull Request para `main` descrevendo a mudança

## Nomenclatura de branches

| Padrão | Uso |
|---|---|
| `feat/<escopo>` | Nova funcionalidade |
| `fix/<escopo>` | Correção de bug |
| `docs/<topico>` | Alterações apenas em documentação |
| `chore/<topico>` | Ferramentas, configuração, infraestrutura |
| `refactor/<escopo>` | Refatoração sem mudança de comportamento |

## Conventional Commits

Formato: `tipo(escopo): descrição breve`

Tipos aceitos: `feat`, `fix`, `docs`, `chore`, `refactor`, `test`, `style`.

Exemplos:

```
feat(ivcf20): calcular escore total do instrumento
docs(visao-geral): detalhar matriz MIV-CFS
chore: configurar lint do monorepo
```

O escopo é opcional; quando usado, deve referenciar o app (`frontend`, `backend`) ou a área afetada.

## Regras de documentação

- **Decisões técnicas relevantes** (stack, bibliotecas, arquitetura, trocas de caminho) devem ser registradas como ADRs em `docs/decisoes/`, numeradas sequencialmente (`ADR-0002-...`, `ADR-0003-...`, ...). Veja o modelo em [docs/decisoes/ADR-0001-monorepo.md](docs/decisoes/ADR-0001-monorepo.md).
- **Alterações de comportamento de domínio** (regras do IVCF-20, EASFI-BR, matriz MIV-CFS) devem atualizar [docs/visao-geral.md](docs/visao-geral.md) no mesmo PR.
- O `README.md` de cada app deve refletir o propósito atual daquele app.

## Checklist de Pull Request

- [ ] A descrição explica **o quê** e **por quê**
- [ ] Commits seguem Conventional Commits
- [ ] Documentação atualizada quando aplicável (docs/, READMEs)
- [ ] Nova decisão técnica registrada como ADR (se aplicável)
- [ ] Nenhum dado sensível, credencial ou segredo foi commitado

## Segurança e dados sensíveis

Este projeto trata dados de saúde de pessoas idosas (dado sensível segundo a LGPD). **Nunca** commite:

- Dados reais de pacientes ou dados de teste verossímeis com informações pessoais
- Credenciais, chaves de API, tokens ou conteúdo de arquivos `.env`

Vulnerabilidades de segurança devem ser reportadas privadamente aos mantenedores, não via issues públicas.

## Estrutura do repositório

Consulte o [README.md](README.md) e [docs/arquitetura.md](docs/arquitetura.md).
