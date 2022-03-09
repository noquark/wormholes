-- links
create table if not exists wh_links (
  id text primary key,
  tag text,
  target text,
  created_at timestamptz not null default now()
);

