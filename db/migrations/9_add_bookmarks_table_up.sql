create table bookmarks (
  id SERIAL PRIMARY KEY NOT NULL,
  sessionId int references sessions(id),
  name text,
  created timestamp,
  modified timestamp
)
