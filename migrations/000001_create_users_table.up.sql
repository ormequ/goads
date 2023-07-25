CREATE TABLE users
(
    id    SERIAL
        CONSTRAINT users_pkey PRIMARY KEY,
    email VARCHAR(320)
        CONSTRAINT users_email_key UNIQUE NOT NULL,
    name  VARCHAR(99)
        CONSTRAINT users_name_key NOT NULL,
    password VARCHAR(72)
        CONSTRAINT users_password_key NOT NULL
);
