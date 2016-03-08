create table challenge_cases (
  id serial primary key NOT NULL,
  challengeId int references challenges(id),
  imput text,
  output text,
  defaultCase boolean,
  created timestamp,
  modified timestamp
)
