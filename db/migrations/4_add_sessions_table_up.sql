create table sessions (
  id serial primary key NOT NULL,
  hash text,
  candidateId int references candidates(id),
  challengeId int references challenges(id),
  expired timestamp,
  status int,
  created timestamp,
  modified timestamp
)
