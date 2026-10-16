# apps/backend

API e núcleo de domínio do **Protocolo Idoso APS** — regras de negócio, persistência e exposição de serviços para o frontend.

## Papel

- **Regras de domínio centrais**: cálculo e classificação do IVCF-20 (0–40) e EASFI-BR (5–25), cruzamento na **matriz MIV-CFS** (níveis 1–4), **marcadores sentinela** com reclassificação direta, prazos de reavaliação por nível
- Persistência de dados sensíveis de saúde (conformidade LGPD — ver [docs/visao-geral.md](../../docs/visao-geral.md#privacidade-e-lgpd))
- Controle de acesso por papel profissional e trilha de auditoria de mudanças de classificação

## Status

**Pendente de implementação.** A stack será definida e registrada em ADR (`docs/decisoes/`).

## Documentação

- [Visão geral do sistema](../../docs/visao-geral.md)
- [Arquitetura](../../docs/arquitetura.md)
