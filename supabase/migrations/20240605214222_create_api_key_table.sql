CREATE TABLE IF NOT EXISTS public."api_key" (
  user_id UUID NOT NULL REFERENCES auth.users,
  api_key CHAR(44) NOT NULL
);

ALTER TABLE public."api_key" ENABLE ROW LEVEL SECURITY;

CREATE POLICY "Restrict access to user api keys to their own api keys"
	ON public."api_key"
	FOR SELECT
	TO authenticated
	USING (user_id = auth.uid());
