#!/bin/bash

source ../shell_scripts/before.sh

psql -U $USER -d iims_test -c "insert into questions(description, created) values('Testing Question', NOW())"

psql -U $USER -d iims_test -c "insert into challenges(description, created) values('Rmlyc3QgVGVzdGluZyBDaGFsbGVuZ2U=', NOW())"

psql -U $USER -d iims_test -c "insert into challenges(description, created) values('U2Vjb25kIFRlc3RpbmcgQ2hhbGxlbmdl', NOW())"

psql -U $USER -d iims_test -c "insert into candidates(email, name, contact, degree, college, yearofcompletion, created) values('ashvin+test@improwised.com', 'Ashvin', '+91 9909970574', 'B.Tech', 'RK University', '2016', NOW())"

psql -U $USER -d iims_test -c "insert into sessions(hash, candidateid, challengeid, status) values('16745e402eefd6f41082fbd68cfe1835ea8fd2b1', '1', '1', '0')"

psql -U $USER -d iims_test -c "insert into challenge_answers(sessionid, answer, language, attempts) values('1', 'cHJpbnQgNQ==', 'python', 1)"

psql -U $USER -d iims_test -c "insert into questions_answers(candidateid, questionsid, answer, created) values('1', '1', 'Answer of First Question', NOW())"