-- DROP TABLE IF EXISTS requests

CREATE TABLE IF NOT EXISTS requests (
    id      serial PRIMARY KEY,
    method  varchar NOT NULL,
    url    varchar NOT NULL,
    host    varchar NOT NULL,
    headers varchar NOT NULL,
    body    varchar NOT NULL
)
