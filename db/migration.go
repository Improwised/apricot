package main

import (
  "github.com/DavidHuie/gomigrate"
  "database/sql"
  _"github.com/lib/pq"
  "fmt"
  "encoding/json"
  "os"
)

type Configuration struct {
  DbName string
  UserName string
}

func main() {
  file, _ := os.Open("../config/configuration.json")
  decoder := json.NewDecoder(file)
  configuration := Configuration{}
  decoder.Decode(&configuration)

  userName := configuration.UserName
  dbName := configuration.DbName

  dbinfo := fmt.Sprintf("user=%s dbname=%s sslmode=disable",
    userName, dbName)
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
