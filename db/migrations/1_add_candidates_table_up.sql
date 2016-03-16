create table candidates (
  id SERIAL PRIMARY KEY NOT NULL,
  email text,
  name text,
  contact text,
  degree text,
  college text,
  yearOfCompletion int,
  created timestamp,
  modified timestamp
)
