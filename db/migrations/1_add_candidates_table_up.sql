create table candidates (
  id serial primary key NOT NULL,
  email text,
  name text,
  contact text,
  degree text,
  college text,
  yearOfCompletion int,
  created timestamp,
  modified timestamp
)
