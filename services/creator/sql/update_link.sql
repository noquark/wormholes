update
  links
set
  target = $1,
  tag = $2
where
  id = $3
