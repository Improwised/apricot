create table challenge_cases (
  id SERIAL PRIMARY KEY NOT NULL,
  challengeId int references challenges(id),
  input text,
  output text,
  defaultCase boolean,
  created timestamp,
  modified timestamp
)
