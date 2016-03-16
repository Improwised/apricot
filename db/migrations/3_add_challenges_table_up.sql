 create table challenges (
  id SERIAL PRIMARY KEY NOT NULL,
  description text,
  deleted timestamp,
  created timestamp,
  modified timestamp
)
