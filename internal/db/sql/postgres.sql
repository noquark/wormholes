-- links
create table if not exists links (
  id text primary key,
  tag text,
  target text,
  created_at timestamptz not null default now()
);

