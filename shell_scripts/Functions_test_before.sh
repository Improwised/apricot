#!/bin/bash

source ../shell_scripts/before.sh

psql -U $USER -d iims_test -c "insert into challenges(description) values('this is first challenge...')"

psql -U $USER -d iims_test -c "insert into challenges(description) values('this is second challenge...')"

psql -U $USER -d iims_test -c "insert into candidates(email, name, contact, degree, college, yearOfCompletion) values('ashvin+test@improwised.com', 'Ashvin', '+91 9712186012', 'B.Tech', 'RK University', '2016')"

psql -U $USER -d iims_test -c "insert into sessions(hash, candidateid, challengeid, status) values('16745e402eefd6f41082fbd68cfe1835ea8fd2b1', '1', '1', '1')"