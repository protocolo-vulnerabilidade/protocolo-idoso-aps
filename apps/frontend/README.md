# apps/frontend

Interface web do **Protocolo Idoso APS** — formulário assistido de avaliação da vulnerabilidade da pessoa idosa para uso pela equipe da UBS (ACS, enfermagem, medicina, eMulti e coordenação).

## Papel

- **Formulário assistido** com os dados nominais e numéricos dos instrumentos IVCF-20 e EASFI-BR — sem nenhum campo de identificação da pessoa idosa
- Exibição do **resultado**: escores, classificações, nível MIV-CFS (1–4), marcadores sentinela e conduta/prazo de reavaliação recomendados
- **Resumo copiável** formatado para transferência à ficha no sistema oficial da UBS
- **Rascunhos e histórico**: retomar avaliação de onde parou; listar, reabrir e **apagar** avaliações
- Mantém o **identificador anônimo de dispositivo** no navegador (sem contas de usuário)
- Consome a **API do backend** — não implementa regras de estratificação

## Status

**Pendente de implementação.** A stack será definida e registrada em ADR (`docs/decisoes/`).

## Documentação

- [Visão geral do sistema](../../docs/visao-geral.md)
- [Arquitetura](../../docs/arquitetura.md)
- [ADR-0002 — Avaliação assistida anônima](../../docs/decisoes/ADR-0002-avaliacao-assistida-anonima.md)
