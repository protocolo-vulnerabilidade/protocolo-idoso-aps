# Visão Geral do Sistema

Documento de escopo do sistema web **Protocolo Idoso APS**, derivado do protocolo assistencial *Identificação, Estratificação e Manejo da Vulnerabilidade Clínico-Funcional e Sociofamiliar da Pessoa Idosa na APS (MIV-CFS)* (Unoeste/Curso de Medicina, set/2026).

## 1. Contexto e problema

O envelhecimento populacional brasileiro desafia a APS/SUS na identificação precoce de condições que comprometem a autonomia, a independência e a qualidade de vida das pessoas idosas. Diagnósticos isolados de doenças crônicas não quantificam com precisão o risco de perda funcional, institucionalização ou mortalidade — a vulnerabilidade é **multidimensional** (clínico-funcional **e** sociofamiliar).

O protocolo padroniza, na UBS, os processos de identificação, avaliação, estratificação, acompanhamento e manejo das pessoas idosas (≥ 60 anos). O sistema dá suporte à etapa de **avaliação e estratificação** dessa linha de cuidado.

## 2. O que o sistema é — e o que não é

**É**: uma ferramenta web de **avaliação assistida**, para uso pontual pelo profissional durante o encontro (consulta ou visita). O profissional informa os dados nominais e numéricos dos instrumentos, e o sistema calcula a estratificação de vulnerabilidade como **valor de referência** para a decisão clínica, gerando um resumo copiável para a ficha da pessoa idosa.

**Não é**:

- **Prontuário ou sistema de registro** — o registro formal do resultado é feito pelo profissional no sistema oficial da UBS; este sistema não substitui nem integra esse registro
- **Sistema de cadastro e acompanhamento longitudinal** — não mantém "ficha" da pessoa idosa, não agenda reavaliações, não gerencia PTS nem indicadores
- **Identificador de pacientes** — nenhum dado de identificação da pessoa idosa é coletado ou armazenado

## 3. Fluxo de uso do profissional

1. **Retomar ou iniciar** — o profissional retoma um rascunho salvo automaticamente ou inicia nova avaliação
2. **Preencher** — formulário assistido com os itens dos instrumentos (idade e demais dados nominais/numéricos dos domínios do IVCF-20 e EASFI-BR) e verificação de marcadores sentinela
3. **Resultado** — escores e classificações do IVCF-20 e EASFI-BR, cruzamento na matriz MIV-CFS, **nível de vulnerabilidade (1–4)** com conduta recomendada e prazo de reavaliação
4. **Transferir** — o profissional copia o **resumo formatado** do resultado e o registra na ficha da pessoa idosa no sistema oficial
5. **Histórico** — a avaliação fica guardada anonimamente no histórico do dispositivo, podendo ser reaberta ou **apagada** a qualquer momento

## 4. Núcleo do domínio (fonte: protocolo assistencial)

### 4.1 Instrumentos

| Instrumento | Dimensão | Domínios | Pontuação | Classificação |
|---|---|---|---|---|
| **IVCF-20** (Moraes et al., 2016) | Clínico-funcional | Idade, autopercepção de saúde, AVDs, cognição, humor/depressão, mobilidade, comunicação, comorbidades | 0–40 | 0–6 baixo risco · 7–14 risco moderado · 15–40 alto risco |
| **EASFI-BR** (Zanela, 2022; derivada da Escala de Gijón) | Sociofamiliar | Situação familiar, econômica, moradia, relações sociais, rede de suporte social | 5–25 | 5–9 adequada · 10–14 risco social · ≥ 15 problema social |

### 4.2 Matriz Integrada MIV-CFS

Os escores **não são somados** (construtos distintos). A condição social **modifica o nível de prioridade assistencial**, conforme o princípio da equidade do SUS:

| IVCF-20 \ EASFI-BR | Adequada (5–9) | Risco social (10–14) | Problema social (≥ 15) |
|---|---|---|---|
| **Baixo risco (0–6)** | Nível 1 (Baixa) | Nível 2 (Moderada) | Nível 3 (Alta) |
| **Risco moderado (7–14)** | Nível 2 (Moderada) | Nível 3 (Alta) | Nível 4 (Muito alta) |
| **Alto risco (15–40)** | Nível 3 (Alta) | Nível 4 (Muito alta) | Nível 4 (Muito alta) |

### 4.3 Marcadores sociais sentinela (reclassificação direta)

A presença de **qualquer** marcador abaixo reclassifica diretamente para **Nível 3 ou 4**, independentemente dos escores:

- Situação de rua, extrema pobreza ou insegurança alimentar grave
- Abandono, suspeita/confirmação de violência doméstica/intrafamiliar e exploração financeira
- Ausência total de cuidador para idoso com dependência nas AVDs
- Sobrecarga grave do cuidador com risco iminente de ruptura do cuidado
- Interrupção de tratamento por impossibilidade financeira ou falta de acesso habitacional/geográfico

### 4.4 Plano de manejo por nível

| Nível | Conduta resumida | Reavaliação | PTS |
|---|---|---|---|
| **1 – Baixa** | Rotina da APS, promoção da saúde, vacinação, prevenção de quedas, revisão de medicamentos, Caderneta do Idoso | Anual | Não obrigatório |
| **2 – Moderada** | Avaliação clínica ampliada, revisão farmacoterapêutica, intervenção sobre fatores reversíveis | Até 6 meses | Indicado se problemas múltiplos/baixa adesão |
| **3 – Alta** | Discussão interdisciplinar, avaliação aprofundada, profissional de referência, visita domiciliar, articulação eMulti/CRAS/CREAS | A cada 3 meses | **Obrigatório** |
| **4 – Muito alta** | PTS prioritário intensivo e compartilhado, suporte matricial, plano de contingentação familiar, articulação intersetorial urgente (Rede de Proteção Social/MP se couber) | 30 a 90 dias | **Prioritário** |

### 4.5 Situações prioritárias e exclusão temporária

**Priorizam a avaliação**: idade ≥ 80 anos, primeiro atendimento, alta hospitalar recente, quedas recorrentes, polifarmácia (≥ 5 medicamentos), multimorbidade, declínio cognitivo/funcional, sintomas depressivos, perda ponderal não explicada; viver sozinho, ausência/insuficiência de cuidador, pobreza ou insegurança alimentar, isolamento social, barreiras habitacionais, dificuldade de acesso a serviços, suspeita/confirmação de negligência ou violência.

**Exclusão temporária** (adiar aplicação até estabilização): emergência clínica, delirium, dor intensa não controlada, incapaidade momentânea de comunicação sem informante secundário confiável.

### 4.6 Indicadores do protocolo (contexto — fora do sistema)

O protocolo define indicadores de monitoramento institucional (cobertura IVCF-20 ≥ 80%, cobertura EASFI-BR ≥ 90%, cobertura de PTS 100% nos Níveis 3–4, reavaliação em prazo ≥ 85%, auditoria de desfechos semestral). **O sistema não os calcula** — são monitorados pela coordenação da UBS a partir dos registros no sistema oficial. A ausência de identificação no sistema é justamente o que impede seu uso para esses fins.

## 5. Capacidades do sistema

> **Nota**: lista preliminar, a validar e detalhar em especificações futuras.

1. **Formulário assistido de avaliação** — coleta dos dados nominais e numéricos dos domínios do IVCF-20 (incluindo idade) e da EASFI-BR, com navegação guiada; **nenhum campo de identificação da pessoa idosa**
2. **Cálculo e classificação** — escore total e classificação do IVCF-20 (0–40) e da EASFI-BR (5–25)
3. **Cruzamento na matriz MIV-CFS** — determinação do nível de vulnerabilidade (1–4)
4. **Marcadores sentinela** — verificação ativa durante a avaliação e **reclassificação automática** para Nível 3/4 quando presente
5. **Conduta recomendada** — exibição do plano de manejo resumido e do prazo de reavaliação conforme o nível
6. **Resumo copiável** — saída formatada (escores, classificações, nível, conduta, data/hora) para copiar e colar na ficha no sistema oficial
7. **Rascunho persistente** — avaliação em andamento é salva automaticamente; ao sair e voltar, o profissional continua de onde parou
8. **Histórico de avaliações** — listagem das avaliações concluídas (por data, escores e nível), com reabertura do resultado e **eliminação individual ou total**
9. **Vínculo por dispositivo** — rascunhos e histórico são vinculados a um identificador anônimo de dispositivo (sem contas de usuário)

## 6. Fora de escopo

- Qualquer **identificação da pessoa idosa** (nome, nº de ficha, data de nascimento identificadora, CNS, etc.)
- Cadastro e acompanhamento longitudinal de pessoas idosas
- Prontuário, gestão de **PTS**, agendamento/lembretes de reavaliação
- Painel de **indicadores** institucionais
- Contas de usuário, perfis de acesso por categoria profissional
- Sincronização do histórico entre dispositivos ou compartilhamento entre profissionais

## 7. Privacidade e LGPD

- **Anônimo por design**: o sistema não possui campos de identificação da pessoa idosa; armazena apenas as respostas dos instrumentos, os resultados calculados e o instante da avaliação
- **Vínculo por dispositivo**: os dados são associados a um identificador anônimo de dispositivo, não a pessoas (profissionais ou idosos)
- **Risco residual de reidentificação**: reconhece-se que, em territórios pequenos, o cruzamento de dados clínicos e sociais detalhados pode permitir inferências. Mitigações: possibilidade de **eliminação a qualquer momento** (individual ou total), ausência de identificadores diretos e minimização de coleta
- **Eliminação**: apagar o histórico é uma funcionalidade de primeira classe do sistema, não um recurso secundário
- **Logs**: a infraestrutura não deve registrar conteúdo clínico das avaliações em logs
- **Registro formal**: a responsabilidade pelo registro do resultado na ficha da pessoa idosa — e por sua proteção — permanece no sistema oficial da UBS
- Aplica-se ainda o **Estatuto da Pessoa Idosa** (Lei nº 10.741/2003), incluindo a notificação compulsória de violência/negligência, que continua sendo feita pelos canais oficiais

## 8. Referências principais

- MORAES, E. N. et al. IVCF-20: reconhecimento rápido do idoso frágil. Rev. Saúde Pública, v. 50, p. 87, 2016. DOI: 10.1590/S1518-8787.2016050006963
- ZANELA, L. S. P. Validação da escala de avaliação sociofamiliar em idosos – versão brasileira. UFMS, 2022.
- BRASIL. Ministério da Saúde. Portaria GM/MS nº 2.528/2006 (PNSPI); Portaria nº 2.436/2017 (PNAB)
- BRASIL. Lei nº 10.741/2003 (Estatuto da Pessoa Idosa, atualizado pela Lei nº 14.423/2022)
- Documento-fonte: `protocolo_manejo_vulnerabilidade_idoso_aps.pdf` (9 págs.)
