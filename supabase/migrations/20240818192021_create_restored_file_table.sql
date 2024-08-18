CREATE TABLE IF NOT EXISTS public."restored_file" (
	id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
	user_id UUID NOT NULL,
	absolute_path VARCHAR(4096) NOT NULL,
	contents TEXT,
	original_timestamp TIMESTAMPTZ NOT NULL
);

ALTER TABLE public."restored_file" ENABLE ROW LEVEL SECURITY;

CREATE POLICY "Restrict access to user files to their own files"
	ON public."restored_file"
	FOR SELECT
	TO authenticated
	USING (user_id = auth.uid());

CREATE POLICY "Users can create a profile."
	ON public."restored_file"
	FOR INSERT
	TO authenticated
	WITH CHECK (auth.uid() = user_id);
