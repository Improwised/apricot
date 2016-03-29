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
  "time"
  "os"
  "bytes"
  "crypto/sha1"
  "encoding/hex"
  "github.com/icza/session"
  "gopkg.in/gomail.v2"
  _"reflect"
  // "strconv"
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

type GeneralInfo struct {
  Id string
  Name string
  Contact string
  Degree string
  College string
  YearOfCompletion string
}

type sessionInfo struct {
  candidateid string
  hash string
}

var cId string
func indexHandler(c web.C, w http.ResponseWriter, r *http.Request) {
  email := r.URL.Query().Get("email")
  //Check whether email empty or not
  if email == "" {
    t, _ := template.ParseFiles("./views/index.html")
    t.Execute(w, t)
  } else {

    var hash string

    stmt, _ := db.Prepare("SELECT id FROM candidates WHERE email = ($1)")
    rows, _ := stmt.Query(email)
    //If candidate already registread before
    if rows.Next() != false {
      query, _ := db.Prepare("SELECT candidateid,hash FROM sessions WHERE candidateid = (select id from candidates where email = ($1))")
      row2,_ := query.Query(email)

       mysession := []sessionInfo{}
       info := sessionInfo{}
      for row2.Next() {
        err := row2.Scan(&info.candidateid,&info.hash)
        mysession = append(mysession, info)
        checkErr(err)
      }
      hash += mysession[0].hash

    } else {//For new registration

      random := time.Now().String()
      random += "abcdefghijklmnopqrstuvwxyzabcdefghijklmnopqrstuvwxyz!@#$^&*()_+"
      h := sha1.New()
      h.Write([]byte(random))
      hash = hex.EncodeToString(h.Sum(nil))

      sesssionHash := session.NewSession()
      session.Add(sesssionHash, w)
      sess := sesssionHash.Id()
      fmt.Println(sess)
      stmt1, err := db.Prepare("INSERT INTO candidates (email,created) VALUES($1,NOW())")
      stmt1.Exec(email)
      checkErr(err)
      stmt2, err := db.Prepare("select id from candidates where email = ($1)")
      checkErr(err)
      rows2, err := stmt2.Query(email)
      information := []GeneralInfo{}
      info := GeneralInfo{}
      for rows2.Next() {
        err := rows2.Scan(&info.Id)
        information = append(information, info)
        checkErr(err)
      }

      stmt3, _ := db.Prepare("insert into questions_answers (candidateId, questionsId, answer, created) values ($1, $2, $3, NOW())")
      rows4, _ := db.Query("select id from questions")
      counter := 1
      for rows4.Next() {
        stmt3.Query(information[0].Id, counter, "")
        counter += 1
      }
      stmt4, _ := db.Prepare("INSERT INTO sessions (hash, candidateId, created) VALUES($1, $2, NOW())")
      stmt4.Exec(hash,information[0].Id)

    }
    m := gomail.NewMessage()
    m.SetHeader("From", "akumbhani666@gmail.com")
    m.SetHeader("To", email)
    m.SetHeader("Subject", "Hello!")
    m.SetBody("text/html", "localhost:8000/information?key=" + hash)
    d := gomail.NewPlainDialer("smtp.gmail.com", 587, "akumbhani666@gmail.com", "9712186012")
    if err := d.DialAndSend(m); err != nil {
        panic(err)

  }
   t, _ := template.ParseFiles("./views/index.html")
    t.Execute(w, t)
  }
}

type GetQuestions struct {
  Questions string
  Id string
  Ans  string
}

type GetAnswers struct {
  Answer string
  Qid string
}

type AllDetail struct {
  GeneralInfo []GeneralInfo
  GetQuestions []GetQuestions
}

var hash string
func informationHandler(c web.C, w http.ResponseWriter, r *http.Request) {
  hash = r.URL.Query().Get("key")
  processFormData(w, r)

  allDetails := AllDetail{}

  // get all questions
  rows, _ := db.Query("select id, description from questions order by sequence")
  questionsInfo := []GetQuestions{}
  qinfo := GetQuestions{}
  for rows.Next() {
    err := rows.Scan(&qinfo.Id, &qinfo.Questions)
    questionsInfo = append(questionsInfo, qinfo)
    checkErr(err)
  }
  stmt2, _ := db.Prepare("SELECT answer, questionsid FROM questions_answers where candidateId = (select candidateId from sessions where hash = ($1))")
  rows2, _ := stmt2.Query(hash)

  for rows2.Next() {
    getanswer := GetAnswers{}
    err := rows2.Scan(&getanswer.Answer, &getanswer.Qid)
    for index, element := range questionsInfo {
      if(element.Id == getanswer.Qid) {
        questionsInfo[index].Ans = getanswer.Answer
      }
    }
    checkErr(err)
  }

  //get user detail
  User := []GeneralInfo{}
  user := GeneralInfo{}
  stmt3, _ := db.Prepare("select id, name, contact, degree, college,yearOfCompletion from candidates where id = (select candidateId from sessions where hash = ($1))")
  row3, _ := stmt3.Query(hash)
  for row3.Next() {
    row3.Scan(&user.Id, &user.Name, &user.Contact, &user.Degree, &user.College, &user.YearOfCompletion)
    User = append(User, user)
  }

  t, _ := template.ParseFiles("./views/information.html")

  allDetails.GeneralInfo = User
  allDetails.GetQuestions = questionsInfo

  t.Execute(w, allDetails)
  dataUpdate(w, r, hash)

}

func dataUpdate(w http.ResponseWriter, r *http.Request, email string) {
  db = setupDB()
  data := r.URL.Query().Get("data")
  id := r.URL.Query().Get("id")
  var table="questions_answers";
  if id == "email" || id == "name" || id == "contact" || id == "degree" || id == "college" || id == "yearOfCompletion"{
    table = "candidates"
  }

  var buffer bytes.Buffer
  buffer.WriteString("UPDATE ")
  buffer.WriteString(table)

  if(table == "questions_answers"){
    buffer.WriteString(" set answer=")
    buffer.WriteString("'" + data + "'")
    buffer.WriteString(",modified=NOW() where questionsid="+ id)
    //  buffer.WriteString("'" + id + "'")
    buffer.WriteString(" AND")
    buffer.WriteString(" candidateid=(select candidateId from sessions where hash=")
    buffer.WriteString("'" + hash + "'")
    buffer.WriteString(")")
  }

  if(table == "candidates"){
    buffer.WriteString(" SET ")
    buffer.WriteString(id)
    buffer.WriteString("=")
    buffer.WriteString("'" + data + "'")
    buffer.WriteString(",modified=NOW() where id=(select candidateId from sessions where hash =")
    buffer.WriteString("'" + hash + "')")
  }
  db.Query(buffer.String())
 }

func processFormData(w http.ResponseWriter, r *http.Request)  {
  db = setupDB()
  r.ParseForm()

  if r.Method == "POST" {
    name := r.FormValue("name")
    contact := r.FormValue("contact")
    degree := r.FormValue("degree")
    college := r.FormValue("college")
    yearOfCompletion := r.FormValue("yearOfCompletion")
    hash := r.FormValue("email")

    for key, values := range r.Form["message"] {   // range over map
      stmt, _ := db.Prepare("update questions_answers set answer=($1),modified=NOW() where candidateid=(select candidateid from sessions where hash=($2)) AND questionsid=($3)")
      stmt.Query(values, hash ,key+1)
    }
    var buffer bytes.Buffer
    buffer.WriteString("UPDATE candidates SET name='" + name +"',contact='" + contact + "',degree='" + degree +"',college='" + college + "',yearOfCompletion='" + yearOfCompletion + "',modified=NOW() where id=(select candidateid from sessions where hash ='" + hash + "')")
    db.Query(buffer.String())
  }
}

func main() {
  db = setupDB()
  defer db.Close()
  goji.Handle("/index", indexHandler)
  goji.Handle("/information", informationHandler)
  http.Handle("/assets/css/", http.StripPrefix("/assets/css/", http.FileServer(http.Dir("assets/css"))))
  http.Handle("/assets/js/", http.StripPrefix("/assets/js/", http.FileServer(http.Dir("assets/js"))))
  http.Handle("/assets/img/", http.StripPrefix("/assets/img/", http.FileServer(http.Dir("assets/img"))))
  http.Handle("/assets/fonts/", http.StripPrefix("/assets/fonts/", http.FileServer(http.Dir("assets/fonts"))))
  goji.Serve()
}
