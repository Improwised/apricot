package main

import "github.com/DavidHuie/gomigrate"
import "database/sql"
import _"github.com/lib/pq"
import "fmt"

func main() {
  dbinfo := fmt.Sprintf("user=mohit dbname=iims sslmode=disable")
  db, err := sql.Open("postgres", dbinfo)
  checkErr(err)
  defer db.Close()
  migrator, err := gomigrate.NewMigrator(db, gomigrate.Postgres{}, "./migrations")
  checkErr(err)
  migrator.Migrate()
}

func checkErr(err error) {
  if err != nil {
    panic(err)
  }
}