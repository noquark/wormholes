-- TimescaleDB
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

