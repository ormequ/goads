CREATE TABLE ads (
    id SERIAL CONSTRAINT ads_pkey PRIMARY KEY,
    author_id INT CONSTRAINT ads_author_id_key NOT NULL REFERENCES users(id),
    published BOOL,
    title VARCHAR(99),
    text VARCHAR(499),
    create_date DATE,
    update_date DATE
);