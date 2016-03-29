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
	"time"
)

type Configuration struct {
	DbName string
	UserName string
}

// db connections
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

// get questions information
type questionsInformation struct {
	Id string
	Description string
	Sequence string
	Flag int
	Deleted *time.Time
}

// link with child structure
type getAllQuestionsInfo struct {
	QuestionsInfo []questionsInformation
}

// display questions in view
func questionsHandler(c web.C, w http.ResponseWriter, r *http.Request) {
	getAllQuestionsInfo := getAllQuestionsInfo{}
	stmt3, _ := db.Prepare("select id, description, deleted, sequence from questions order by id")
	rows3, _ := stmt3.Query()

	questionsInfo := []questionsInformation{}
	q := questionsInformation{}
	for rows3.Next() {
		err := rows3.Scan(&q.Id, &q.Description, &q.Deleted, &q.Sequence)
		if q.Deleted != nil{
			q.Flag = 1
		} else if q.Deleted == nil{
			q.Flag = 0
		}
		questionsInfo = append(questionsInfo, q)
		checkErr(err)
	}

	getAllQuestionsInfo.QuestionsInfo = questionsInfo
	t, _ := template.ParseFiles("./views/questions.html")
	t.Execute(w, getAllQuestionsInfo)
}

// perform edit functionality
func editHandler(c web.C, w http.ResponseWriter, r *http.Request) {
	if r.FormValue("qId") != "" {
		description := r.FormValue("description")
		sequence := r.FormValue("sequence")
		qId := r.FormValue("qId")
		stmt1, _ := db.Prepare("update questions set description = ($1), sequence = ($2) where id = ($3)")
		stmt1.Query(description, sequence, qId)
		http.Redirect(w, r, "questions", 301)
	} else {
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
}

//  delete questions functionality
func deleteHandler(c web.C, w http.ResponseWriter, r *http.Request) {
	qId := r.URL.Query().Get("qid")
	status := r.URL.Query().Get("deleted")
	if status == "no" {
		stmt1, _ := db.Prepare("update questions set deleted = NOW() where id = ($1)")
		stmt1.Query(qId)
	} else if status == "yes" {
		stmt1, _ := db.Prepare("update questions set deleted = NULL where id = ($1)")
		stmt1.Query(qId)
	}
}

// add questions functionality
func addHandler(c web.C, w http.ResponseWriter, r *http.Request) {
	description := r.FormValue("description")
	sequence := r.FormValue("sequence")
	stmt2, _ := db.Prepare("insert into questions (description, sequence, created) values($1, $2, NOW())")
	stmt2.Query(description, sequence)
	http.Redirect(w, r, "questions", 301)
}

// get general information of users
type GeneralInfo struct {
	Id string
	Name string
	Contact string
	Degree string
	College string
	YearOfCompletion string
	Email string
	Created time.Time
	Modified time.Time
}

func candidateHandler(c web.C, w http.ResponseWriter, r *http.Request) {
	stmt1 := fmt.Sprintf("select id, name, email, contact, degree, college, yearOfCompletion, created, modified from candidates")
	rows1, _ := db.Query(stmt1)

	UsersInfo := []GeneralInfo{}
	user := GeneralInfo{}
	for rows1.Next() {
		rows1.Scan(&user.Id, &user.Name, &user.Email, &user.Contact, &user.Degree, &user.College, &user.YearOfCompletion ,&user.Created, &user.Modified)
		UsersInfo = append(UsersInfo, user)
	}

	t, _ := template.ParseFiles("./views/candidate.html")
	t.Execute(w, UsersInfo)
}

func main() {
	db = setupDB()
	defer db.Close()
	goji.Handle("/questions", questionsHandler)
	goji.Handle("/candidates", candidateHandler)
	goji.Handle("/add", addHandler)
	goji.Handle("/edit", editHandler)
	goji.Handle("/delete", deleteHandler)
	http.Handle("/assets/css/", http.StripPrefix("/assets/css/", http.FileServer(http.Dir("assets/css"))))
	http.Handle("/assets/js/", http.StripPrefix("/assets/js/", http.FileServer(http.Dir("assets/js"))))
	http.Handle("/assets/img/", http.StripPrefix("/assets/img/", http.FileServer(http.Dir("assets/img"))))
	http.Handle("/assets/fonts/", http.StripPrefix("/assets/fonts/", http.FileServer(http.Dir("assets/fonts"))))
	goji.Serve()
}
