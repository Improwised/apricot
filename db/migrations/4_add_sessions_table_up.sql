create table sessions (
  id SERIAL PRIMARY KEY NOT NULL,
  hash text,
  candidateId int references candidates(id),
  challengeId int references challenges(id),
  expired timestamp,
  status int,
  created timestamp,
  modified timestamp
)
