CREATE TABLE IF NOT EXISTS private."file" (
	id SERIAL PRIMARY KEY,
	user_id UUID NOT NULL,
	absolute_path VARCHAR(4096) NOT NULL,
	contents TEXT,
	timestamp TIMESTAMPTZ NOT NULL
);
