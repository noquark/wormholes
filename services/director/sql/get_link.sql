select
  id,
  target,
  tag
from
  wh_links
where
  id = $1;

