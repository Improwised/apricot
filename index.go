package main

import (
  "fmt"
  "net/http"
  "database/sql"
  _"github.com/lib/pq"
  "github.com/zenazn/goji"
  "github.com/zenazn/goji/web"
  "html/template"
  "encoding/json"
  "os"
)

type Configuration struct {
  DbName string
  UserName string
}

var db *sql.DB
func setupDB() *sql.DB {
  file, _ := os.Open("./config/configuration.json")
  decoder := json.NewDecoder(file)
  configuration := Configuration{}
  decoder.Decode(&configuration)

  userName := configuration.UserName
  dbName := configuration.DbName

  dbinfo := fmt.Sprintf("user=%s dbname=%s sslmode=disable",
    userName, dbName)
  db, err := sql.Open("postgres", dbinfo)
  checkErr(err)
  return db
}

func checkErr(err error) {
  if err != nil {
    panic(err)
  }
}

func indexHandler(c web.C, w http.ResponseWriter, r *http.Request) {
  email := r.URL.Query().Get("Email")
  if email == "" {
    t, _ := template.ParseFiles("./views/index.html")
    t.Execute(w, t)
  }
  stmt, _ := db.Prepare("SELECT * FROM candidates WHERE email = ($1)")
  rows, _ := stmt.Query(email)
  if rows.Next() != false {
    fmt.Println("Email already registered")
  }
}

func main() {
  db = setupDB()
  defer db.Close()
  goji.Handle("/index", indexHandler)
  http.Handle("/assets/css/", http.StripPrefix("/assets/css/", http.FileServer(http.Dir("assets/css"))))
  http.Handle("/assets/img/", http.StripPrefix("/assets/img/", http.FileServer(http.Dir("assets/img"))))
  http.Handle("/assets/fonts/", http.StripPrefix("/assets/fonts/", http.FileServer(http.Dir("assets/fonts"))))
  goji.Serve()
}
