create table questions (
  id SERIAL PRIMARY KEY NOT NULL,
  description text,
  deleted timestamp,
  sequence int,
  created timestamp,
  modified timestamp
)
