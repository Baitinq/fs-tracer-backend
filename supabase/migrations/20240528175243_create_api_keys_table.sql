CREATE TABLE IF NOT EXISTS private.api_keys (
  id UUID NOT NULL REFERENCES auth.users,
  api_key CHAR(44) NOT NULL
);
