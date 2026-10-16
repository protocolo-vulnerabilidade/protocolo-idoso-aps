# ADR-0002 — Avaliação assistida anônima com histórico por dispositivo

- **Status**: Aceita
- **Data**: 2026-10-16
- **Relação**: Amplia o [ADR-0001](ADR-0001-monorepo.md) (monorepo frontend+backend mantido)

## Contexto

O escopo do sistema foi deliberadamente reduzido: ele **não é um prontuário nem um sistema de acompanhamento longitudinal**. O profissional da saúde aplica os instrumentos (IVCF-20 e EASFI-BR), o sistema calcula a estratificação de vulnerabilidade (matriz MIV-CFS, níveis 1–4, marcadores sentinela) e o resultado é **registrado manualmente pelo profissional na ficha da pessoa idosa no sistema oficial da UBS**.

Duas necessidades surgiram desse escopo:

1. **Histórico** das avaliações realizadas, consultável e **apagável**;
2. **Retomada de rascunho**: o profissional pode sair da tela no meio da avaliação e voltar exatamente de onde parou.

Ambas exigem persistência — o sistema deixou de ser stateless. Foi preciso decidir **o que** é persistido, **onde** e **a quem pertence**.

## Decisão

1. **Anônimo por design**: o sistema não possui campos de identificação da pessoa idosa (nome, nº de ficha, CNS, etc.). Persistem apenas as respostas dos instrumentos, os resultados calculados e o instante da avaliação. A associação entre resultado e pessoa idosa permanece no sistema oficial da UBS.
2. **Persistência no backend com banco de dados** (não no navegador): histórico centralizado, disponível em qualquer dispositivo da UBS que acesse o sistema.
3. **Vínculo por dispositivo, sem contas de usuário**: rascunhos e histórico pertencem a um **identificador anônimo de dispositivo** gerado no primeiro acesso e mantido no navegador. Não há login, cadastro ou perfis de acesso.
4. **Eliminação como funcionalidade de primeira classe**: o histórico pode ser apagado total ou individualmente, a qualquer momento, direto na interface.
5. A interface é um **formulário assistido completo** (dados nominais e numéricos dos domínios dos instrumentos) — não uma tela de calculadora — e o resultado é apresentado com **resumo copiável** para transferência ao sistema oficial.

## Consequências

**Positivas**
- Privacidade maximizada: sem identificadores diretos, o dado persistido tem risco reduzido de exposição; LGPD substancialmente simplificada
- Baixa fricção de adoção: sem cadastro de contas, senha ou administração de usuários
- Histórico e rascunhos funcionam de forma uniforme em qualquer dispositivo da UBS

**Negativas / trade-offs (aceitas conscientemente)**
- **Computador compartilhado**: sem contas, todos que usam o mesmo dispositivo compartilham o mesmo histórico — não há separação por profissional
- **Órfãos**: se o navegador for limpo (ou o token do dispositivo perdido), os rascunhos e o histórico ficam inacessíveis por falta de vínculo
- **Sem sincronização**: o histórico de um dispositivo não aparece em outro
- **Risco residual de reidentificação**: em territórios pequenos, o cruzamento de dados clínicos/sociais detalhados pode permitir inferências — mitigado pela ausência de identificadores diretos, minimização de coleta e eliminabilidade a qualquer momento

## Alternativas consideradas

1. **Histórico apenas no navegador (localStorage)** — rejeitada: não atravessa dispositivos da UBS e expõe dados clínicos em máquinas compartilhadas sem controle central de eliminação.
2. **Contas de usuário por profissional (login)** — rejeitada por ora: agrega complexidade (autenticação, gestão de contas, recuperação de senha) desproporcional ao escopo; pode ser adotada futuramente via novo ADR.
3. **Campo opcional de identificação do idoso (nome/nº da ficha) no cálculo** — rejeitada: tornaria o sistema tratamento de dado pessoal de saúde com armazenamento, obrigando controles de LGPD plenos e desvirtuando o propósito de ferramenta de referência anônima.
4. **Sistema stateless puro (sem histórico, sem rascunho)** — rejeitada: perde a necessidade expressa de histórico e retomada de rascunho.

## Referências

- [docs/visao-geral.md](../visao-geral.md) — escopo, capacidades e fora de escopo
- [docs/arquitetura.md](../arquitetura.md) — restrições transversais (anonimato, eliminação, logs)
