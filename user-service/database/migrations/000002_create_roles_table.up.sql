create table if not exists "roles"(
    id serial primary key,
    name varchar(255) unique not null,
    created_at timestamp not null default current_timestamp,
    updated_at timestamp null,
    deleted_at timestamp null
)