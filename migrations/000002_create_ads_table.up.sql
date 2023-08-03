CREATE TABLE ads
(
    id          BIGINT PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
    author_id   BIGINT
        CONSTRAINT ads_author_id_key NOT NULL REFERENCES users (id) ON DELETE CASCADE,
    published   BOOL,
    title       TEXT,
    text        TEXT,
    create_date DATE,
    update_date DATE
);
