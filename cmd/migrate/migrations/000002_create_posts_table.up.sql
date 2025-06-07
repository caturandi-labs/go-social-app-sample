CREATE TABLE IF NOT EXISTS posts (
    id bigserial PRIMARY KEY,
    title text not null,
    user_id bigint not null,
    content text not null,
    tags varchar(300) [],
    created_at timestamp(0) with time zone not null default now(),
    updated_at timestamp(0) with time zone null default now()
)
