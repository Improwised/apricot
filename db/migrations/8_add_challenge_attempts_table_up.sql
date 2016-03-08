create table challenge_attempts (
  id serial primary key NOT NULL,
  sessionId int references sessions(id),
  input text,
  output text,
  solution text,
  status int,
  created timestamp,
  modified timestamp
-- all the response from hankerrank
)
