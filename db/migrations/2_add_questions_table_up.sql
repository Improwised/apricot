create table questions (
  id serial primary key NOT NULL,
  description text,
  deleted timestamp,
  sequence int,
  created timestamp,
  modified timestamp
)
