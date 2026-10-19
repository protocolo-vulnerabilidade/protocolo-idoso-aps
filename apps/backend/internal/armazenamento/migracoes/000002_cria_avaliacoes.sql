CREATE TABLE avaliacoes (
  id UUID PRIMARY KEY,
  dispositivo_id UUID NOT NULL REFERENCES dispositivos (id) ON DELETE CASCADE,
  situacao TEXT NOT NULL CHECK (situacao IN ('rascunho', 'concluida')),
  etapa SMALLINT NOT NULL DEFAULT 0 CHECK (etapa BETWEEN 0 AND 6),
  respostas JSONB NOT NULL,
  escore_ivcf SMALLINT CHECK (escore_ivcf IS NULL OR escore_ivcf BETWEEN 0 AND 40),
  classe_ivcf TEXT CHECK (classe_ivcf IS NULL OR classe_ivcf IN ('baixo', 'moderado', 'alto')),
  escore_easfi SMALLINT CHECK (escore_easfi IS NULL OR escore_easfi BETWEEN 5 AND 25),
  classe_easfi TEXT CHECK (classe_easfi IS NULL OR classe_easfi IN ('adequada', 'risco', 'problema')),
  nivel_matriz SMALLINT CHECK (nivel_matriz IS NULL OR nivel_matriz BETWEEN 1 AND 4),
  nivel_final SMALLINT CHECK (nivel_final IS NULL OR nivel_final BETWEEN 1 AND 4),
  sentinelas_positivas SMALLINT[],
  concluida_em TIMESTAMPTZ,
  criado_em TIMESTAMPTZ NOT NULL DEFAULT now(),
  atualizado_em TIMESTAMPTZ NOT NULL DEFAULT now(),
  CONSTRAINT avaliacoes_coerencia CHECK (
    (
      situacao = 'rascunho'
      AND concluida_em IS NULL
      AND escore_ivcf IS NULL
      AND classe_ivcf IS NULL
      AND escore_easfi IS NULL
      AND classe_easfi IS NULL
      AND nivel_matriz IS NULL
      AND nivel_final IS NULL
    )
    OR (
      situacao = 'concluida'
      AND concluida_em IS NOT NULL
      AND escore_ivcf IS NOT NULL
      AND classe_ivcf IS NOT NULL
      AND escore_easfi IS NOT NULL
      AND classe_easfi IS NOT NULL
      AND nivel_matriz IS NOT NULL
      AND nivel_final IS NOT NULL
    )
  )
);

CREATE INDEX idx_avaliacoes_dispositivo ON avaliacoes (dispositivo_id, situacao, atualizado_em DESC);
