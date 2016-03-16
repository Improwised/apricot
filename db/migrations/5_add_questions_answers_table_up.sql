create table questions_answers (
  id SERIAL PRIMARY KEY NOT NULL,
  candidateId int references candidates(id),
  questionsId int references questions(id),
  answer text,
  created timestamp,
  modified timestamp
)
