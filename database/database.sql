-- enable uuid
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
-- users
CREATE TABLE IF NOT EXISTS wh_users (
  id uuid DEFAULT uuid_generate_v4(),
  email text NOT NULL UNIQUE,
  hashed_password varchar(255) NOT NULL,
  created_at timestamptz NOT NULL default now(),
  CONSTRAINT wh_users_pk PRIMARY KEY (id)
);
-- links
CREATE TABLE IF NOT EXISTS wh_links (
  id varchar(255),
  tag varchar(255),
  target varchar(255),
  created_at timestamptz NOT NULL default now(),
  CONSTRAINT wh_links_pk PRIMARY KEY (id)
);
