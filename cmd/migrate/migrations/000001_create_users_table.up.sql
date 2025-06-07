CREATE EXTENSION IF NOT EXISTS citext;
CREATE TABLE IF NOT EXISTS users(
    id bigserial PRIMARY KEY,
    username varchar(255) UNIQUE NOT NULL,
    email citext UNIQUE NOT NULL,
    password varchar NOT NULL,
    created_at timestamp(0) WITH time zone not null default NOW(),
    updated_at timestamp(0) NULL default NOW()
)