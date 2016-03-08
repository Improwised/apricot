create table challenge_answers (
  id serial primary key NOT NULL,
  sessionId int references sessions(id),
  answer text,
  attempts int,
  created timestamp,
  modified timestamp
)
