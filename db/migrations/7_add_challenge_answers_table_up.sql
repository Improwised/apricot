create table challenge_answers (
  id SERIAL PRIMARY KEY NOT NULL,
  sessionId int references sessions(id),
  answer text,
  attempts int,
  created timestamp,
  modified timestamp
)
