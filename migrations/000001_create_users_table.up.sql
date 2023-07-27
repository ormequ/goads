CREATE TABLE users
(
    id    SERIAL PRIMARY KEY,
    email TEXT
        CONSTRAINT users_email_key UNIQUE NOT NULL,
    name  TEXT NOT NULL,
    password VARCHAR(72) NOT NULL
);
