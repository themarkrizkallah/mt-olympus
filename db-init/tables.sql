create extension if not exists "uuid-ossp";

create table if not exists users(
   id uuid primary key default uuid_generate_v1(),
   email text unique not null,
   password text not null,
   created_at timestamp not null default now()
);

create table if not exists api_keys(
   id uuid primary key default uuid_generate_v1(),
   user_id uuid references users(id),
   created_at timestamp not null default now()
);

create table if not exists assets(
    id serial primary key,
    name text unique not null,
    tick varchar(12) unique not null
);

create table if not exists accounts(
    id uuid primary key default uuid_generate_v1(),
    user_id uuid references users(id),
    asset_id serial references assets(id),
    balance bigint not null default 0,
    holds bigint not null default 0
);