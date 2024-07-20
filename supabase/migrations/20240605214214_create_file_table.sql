CREATE TABLE IF NOT EXISTS public."file" (
	id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
	user_id UUID NOT NULL,
	absolute_path VARCHAR(4096) NOT NULL,
	contents TEXT,
	timestamp TIMESTAMPTZ NOT NULL
);

ALTER TABLE public."file" ENABLE ROW LEVEL SECURITY;

CREATE POLICY "Restrict access to user files to their own files"
	ON public."file"
	FOR SELECT
	TO authenticated
	USING (user_id = auth.uid());
