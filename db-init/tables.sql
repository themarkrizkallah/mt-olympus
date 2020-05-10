-- noinspection SqlNoDataSourceInspectionForFile

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
    id uuid primary key default uuid_generate_v1(),
    name varchar(100) unique not null,
    tick varchar(12) unique not null,
    created_at timestamp not null default now()
);

create table if not exists products(
   id varchar(25) primary key,
   base_id uuid references assets(id) not null,
   quote_id uuid references assets(id) not null,
   base_tick varchar(12) references assets(tick) not null,
   quote_tick varchar(12) references assets(tick) not null,
   created_at timestamp not null default now()
);

create table if not exists accounts(
    id uuid primary key default uuid_generate_v1(),
    user_id uuid references users(id),
    asset_id uuid references assets(id),
    balance bigint not null default 0,
    holds bigint not null default 0,
    created_at timestamp not null default now()
);

create table if not exists orders(
     id uuid unique not null,
     product_id varchar(25) references products(id) not null,
     user_id uuid references users(id) not null,
     amount bigint not null,
     price bigint not null default 0,          -- 0 if it's a market order
     type varchar(6) not null default 'limit', -- one of {LIMIT, MARKET, STOP}
     side boolean not null,                    -- Buy = True, Sell = False
     status varchar(16) not null,              -- one of {Filled, Partially Filled, Confirmed}
     created_at timestamp not null
);

-- Initial asset setup
insert into assets(name, tick) values('US Dollar', 'USD');
insert into assets(name, tick) values('Bitcoin', 'BTC');

-- Initial product setup
with usd as (select id, tick from assets where tick = 'USD' limit 1),
     btc as (select id, tick from assets where tick = 'BTC' limit 1)
insert into products(id, base_id, quote_id, base_tick, quote_tick)
values (
    concat((select tick from btc), '-', (select tick from usd)),
    (select id from btc),
    (select id from usd),
    (select tick from btc),
    (select tick from usd)
);