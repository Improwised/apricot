#!/bin/bash

psql -U $USER -c "DROP DATABASE IF EXISTS iims_test"

psql -U $USER -c "CREATE DATABASE iims_test WITH OWNER $USER"

export GO_ENV2=testing

go run ../db/migration.go
