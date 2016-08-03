#!/bin/bash

source ../shell_scripts/before.sh

psql -U $USER -d iims_test -c "insert into challenges(description) values('Rmlyc3QgVGVzdGluZyBDaGFsbGVuZ2U=')"