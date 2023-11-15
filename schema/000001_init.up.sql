create table if not exists messages
(
    id              serial primary key,
    public_key_from varchar not null,
    public_key_to   varchar not null,
    message         bytea    not null,
    created_at      timestamp with time zone default now(),
    updated_at      timestamp with time zone default now()
)