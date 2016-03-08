 create table challenges (
  id serial primary key NOT NULL,
  description text,
  deleted timestamp,
  created timestamp,
  modified timestamp
)
