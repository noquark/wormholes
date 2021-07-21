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

