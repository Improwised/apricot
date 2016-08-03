package main

import (
  "github.com/DavidHuie/gomigrate"
  "database/sql"
  _"github.com/lib/pq"
  "os"
  "github.com/improwised/apricot/dbconfig"
)

func main() {
  var connectionString string
  var db *sql.DB

  if(os.Getenv("GO_ENV2") == ""){
    os.Setenv("GO_ENV2", "production")
  }

  if(os.Getenv("GO_ENV2") == "production"){
      connectionString = dbconfig.PostgresConnectionString("../dbconfig/test-files/production.json", "disable")
      db, _ = sql.Open("postgres", connectionString)

      migrator, err2 := gomigrate.NewMigrator(db, gomigrate.Postgres{}, "./migrations")
      if err2 != nil {
        panic(err2)
      }
      migrator.Migrate()
  } else if((os.Getenv("GO_ENV2")) == "development") {
      connectionString = dbconfig.PostgresConnectionString("./dbconfig/test-files/settings.json", "disable")
      db, _ = sql.Open("postgres", connectionString)
  } else {
      connectionString = dbconfig.PostgresConnectionString("../dbconfig/test-files/testing.json", "disable")
      db, _ = sql.Open("postgres", connectionString)

      migrator, err2 := gomigrate.NewMigrator(db, gomigrate.Postgres{}, "../db/migrations")
      if err2 != nil {
        panic(err2)
      }
      migrator.Migrate()
   }
//************ uncomment following code to show connected user's info  **********
  // settings := dbconfig.Settings(connectionString)
  // connection := []string{
  //   "host=", settings["host"], " ",
  //   "password=", settings["password"], " ",
  //   "user=", settings["username"], " ",
  //   "dbname=", settings["database"], " ",
  //   }
  // fmt.Println(strings.Join(connection, ""))
//**********************************************************************
    defer db.Close()
}
