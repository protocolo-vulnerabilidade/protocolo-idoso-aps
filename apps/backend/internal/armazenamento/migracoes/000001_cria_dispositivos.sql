CREATE TABLE dispositivos (
  id UUID PRIMARY KEY,
  primeiro_acesso_em TIMESTAMPTZ NOT NULL DEFAULT now()
);
