create table questions_answers (
  id serial primary key NOT NULL,
  candidateId int references candidates(id),
  questionsId int references questions(id),
  answer text,
  created timestamp,
  modified timestamp
)
