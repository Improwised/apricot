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

type questionsInformation struct {
  Id string
  Description string
  Sequence string
}

var cId string
func questionsHandler(c web.C, w http.ResponseWriter, r *http.Request) {
  if r.FormValue("qId") != "" {
    description := r.FormValue("description")
    sequence := r.FormValue("sequence")
    qId := r.FormValue("qId")
    stmt1, _ := db.Prepare("update questions set description = ($1), sequence = ($2) where id = ($3)")
    stmt1.Query(description, sequence, qId)
    fmt.Println(description)
  }

  if r.FormValue("submit") != "" {
    description := r.FormValue("description")
    sequence := r.FormValue("sequence")
    stmt2, _ := db.Prepare("insert into questions (description, sequence, created) values($1, $2, NOW())")
    stmt2.Query(description, sequence)
  }
  stmt3, _ := db.Prepare("select id, description from questions order by id")
  rows3, _ := stmt3.Query()
  questions := []questionsInformation{}
  q := questionsInformation{}
  for rows3.Next() {
    err := rows3.Scan(&q.Id, &q.Description)
    questions = append(questions, q)
    checkErr(err)
  }

  t, _ := template.ParseFiles("./views/questions.html")
  t.Execute(w, questions)
}

func editHandler(c web.C, w http.ResponseWriter, r *http.Request) {
  qId := r.URL.Query().Get("qid")
  stmt1, _ := db.Prepare("select description, sequence from questions where id = ($1)")
  rows1, _ := stmt1.Query(qId)
  questions := []questionsInformation{}
  q := questionsInformation{}
  for rows1.Next() {
    err := rows1.Scan(&q.Description, &q.Sequence)
    questions = append(questions, q)
    checkErr(err)
  }

  t, _ := template.ParseFiles("./views/edit.html")
  t.Execute(w, questions)

}

func main() {
  db = setupDB()
  defer db.Close()
  goji.Handle("/questions", questionsHandler)
  goji.Handle("/edit", editHandler)
  http.Handle("/assets/css/", http.StripPrefix("/assets/css/", http.FileServer(http.Dir("assets/css"))))
  http.Handle("/assets/js/", http.StripPrefix("/assets/js/", http.FileServer(http.Dir("assets/js"))))
  http.Handle("/assets/img/", http.StripPrefix("/assets/img/", http.FileServer(http.Dir("assets/img"))))
  http.Handle("/assets/fonts/", http.StripPrefix("/assets/fonts/", http.FileServer(http.Dir("assets/fonts"))))
  goji.Serve()
}
