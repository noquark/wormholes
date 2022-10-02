select
  id,
  target,
  tag
from
  links
where
  id = $1
