CREATE TABLE IF NOT EXISTS snippets (
    id SERIAL PRIMARY KEY,
    title TEXT,
    content TEXT NOT NULL,
    created TIMESTAMPTZ DEFAULT NOW(),
    expires TIMESTAMPTZ
);