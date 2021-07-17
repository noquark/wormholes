-- enable uuid
create extension if not exists "uuid-ossp";

-- users
create table if not exists wh_users (
  id uuid default uuid_generate_v4 () primary key,
  is_admin bool default false,
  email text not null unique,
  hashed_password text not null,
  created_at timestamptz not null default now()
);

-- links
create table if not exists wh_links (
  id text primary key,
  tag text,
  target text,
  created_at timestamptz not null default now()
);

-- clicks
create table if not exists wh_clicks (
  time timestamptz not null,
  link text not null,
  tag text not null,
  cookie text not null,
  is_mobile boolean,
  is_bot boolean,
  browser text,
  browser_version text,
  os text,
  os_version text,
  platform text,
  ip text,
  lat double precision,
  long double precision,
  city text,
  region text,
  country text,
  continent text
);

-- create hypertable
select
  create_hypertable ('wh_clicks', 'time', if_not_exists => true);

