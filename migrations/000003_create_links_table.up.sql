CREATE TABLE links
(
    id        SERIAL PRIMARY KEY,
    alias     TEXT
        CONSTRAINT links_alias_key UNIQUE NOT NULL,
    url       TEXT                        NOT NULL,
    author_id BIGINT
        CONSTRAINT links_author_id_key NOT NULL REFERENCES users (id) ON DELETE CASCADE
);
