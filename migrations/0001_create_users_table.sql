CREATE TABLE users (
    id SERIAL CONSTRAINT users_pkey PRIMARY KEY,
    email VARCHAR(320) CONSTRAINT users_email_key UNIQUE NOT NULL,
    name VARCHAR(50) CONSTRAINT users_name_key NOT NULL
);
