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

type generalInfo struct {
  Id string
  Description string
}

var cId string
func indexHandler(c web.C, w http.ResponseWriter, r *http.Request) {
  // stmt1, _ := db.Prepare("INSERT INTO candidates (email, created) VALUES($1,NOW())")
  // stmt1.Exec(email)

  stmt2, _ := db.Prepare("select id, description from questions")
  rows2, _ := stmt2.Query()

  information := []generalInfo{}
  info := generalInfo{}
  for rows2.Next() {
    err := rows2.Scan(&info.Id, &info.Description)
    information = append(information, info)
    checkErr(err)
  }

  fmt.Println(information)
  // stmt3, _ := db.Prepare("insert into questions_answers (candidateId, questionsId, answer, created) values ($1, $2, $3, NOW())")
  // rows4, _ := db.Query("select id from questions")
  // counter := 1
  // for rows4.Next() {
  //   stmt3.Query(information[0].Id, counter, "")
  //   counter += 1
  //   fmt.Println(counter)
  // }
  t, _ := template.ParseFiles("./views/questions.html")
  t.Execute(w, information)
}

func main() {
  db = setupDB()
  defer db.Close()
  goji.Handle("/questions", indexHandler)
  http.Handle("/assets/css/", http.StripPrefix("/assets/css/", http.FileServer(http.Dir("assets/css"))))
  http.Handle("/assets/js/", http.StripPrefix("/assets/js/", http.FileServer(http.Dir("assets/js"))))
  http.Handle("/assets/img/", http.StripPrefix("/assets/img/", http.FileServer(http.Dir("assets/img"))))
  http.Handle("/assets/fonts/", http.StripPrefix("/assets/fonts/", http.FileServer(http.Dir("assets/fonts"))))
  goji.Serve()
}
