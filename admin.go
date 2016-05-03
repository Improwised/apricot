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
	"strings"
	"bytes"
		// "reflect"
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
	var buffer bytes.Buffer
	buffer.WriteString("select id, description, deleted, sequence from questions order by sequence")
	rows3, _ := db.Query(buffer.String())

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
func editQuesionHandler(c web.C, w http.ResponseWriter, r *http.Request) {
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
		t, _ := template.ParseFiles("./views/editquestion.html")
		t.Execute(w, questions)
	}
}

//  delete questions functionality
func deleteQuestionHandler(c web.C, w http.ResponseWriter, r *http.Request) {
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
func addQuestionsHandler(c web.C, w http.ResponseWriter, r *http.Request) {
	description := r.FormValue("description")
	sequence := r.FormValue("sequence")
	stmt2, _ := db.Prepare("insert into questions (description, sequence, created) values($1, $2, NOW())")
	stmt2.Query(description, sequence)
	http.Redirect(w, r, "questions", 301)
}

func challengesHandler(c web.C, w http.ResponseWriter, r *http.Request) {
	getAllQuestionsInfo := getAllQuestionsInfo{}
	var buffer bytes.Buffer
	buffer.WriteString("select id, description, deleted from challenges")
	rows3, _ := db.Query(buffer.String())

	questionsInfo := []questionsInformation{}
	q := questionsInformation{}
	for rows3.Next() {
		err := rows3.Scan(&q.Id, &q.Description, &q.Deleted)
		if q.Deleted != nil{
			q.Flag = 1
		} else if q.Deleted == nil{
			q.Flag = 0
		}
		questionsInfo = append(questionsInfo, q)
		checkErr(err)
	}

	getAllQuestionsInfo.QuestionsInfo = questionsInfo
	t, _ := template.ParseFiles("./views/programmingtest.html")
	t.Execute(w, getAllQuestionsInfo)
}

func deleteChanllengesHandler(c web.C, w http.ResponseWriter, r *http.Request) {
	qId := r.URL.Query().Get("qid")
	status := r.URL.Query().Get("deleted")
	if status == "no" {
		stmt1, _ := db.Prepare("update challenges set deleted = NOW() where id = ($1)")
		stmt1.Query(qId)
	} else if status == "yes" {
		stmt1, _ := db.Prepare("update challenges set deleted = NULL where id = ($1)")
		stmt1.Query(qId)
	}
}

func editChallengeHandler(c web.C, w http.ResponseWriter, r *http.Request) {
	if r.FormValue("qId") != "" {
		description := r.FormValue("description")
		qId := r.FormValue("qId")
		stmt1, _ := db.Prepare("update challenges set description = ($1) where id = ($2)")
		stmt1.Query(description, qId)
		http.Redirect(w, r, "programmingtest", 301)
	} else {
		qId := r.URL.Query().Get("qid")
		stmt1, _ := db.Prepare("select description from challenges where id = ($1)")
		rows1, _ := stmt1.Query(qId)
		questions := []questionsInformation{}
		q := questionsInformation{}
		for rows1.Next() {
			err := rows1.Scan(&q.Description)
			questions = append(questions, q)
			checkErr(err)
		}
		t, _ := template.ParseFiles("./views/editchallenge.html")
		t.Execute(w, questions)
	}
}

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
	ChallengeAttempts string
	DateOnly string
	QuestionsAttended string
}

// display candidates information
func candidateHandler(c web.C, w http.ResponseWriter, r *http.Request) {
	stmt1 := fmt.Sprintf("SELECT c.id,c.name, c.email, c.degree, c.college, c.yearOfCompletion, c.modified, max(c1.attempts) FROM candidates c JOIN sessions s ON c.id = s.candidateid JOIN challenge_answers c1 ON s.id = c1.sessionid where s.status=0 group by c.id")
	rows1, _ := db.Query(stmt1)

	UsersInfo := []GeneralInfo{}
	user := GeneralInfo{}
	for rows1.Next() {

		rows1.Scan(&user.Id, &user.Name, &user.Email, &user.Degree, &user.College, &user.YearOfCompletion, &user.Modified, &user.ChallengeAttempts)

		//extract only date from timestamp========
	 	str :=&user.Modified
		str1 := str.String()
		s := strings.Split(str1," ")
		user.DateOnly = s[0]
		//================================

		stmt2 := fmt.Sprintf("SELECT count(id) FROM questions_answers WHERE length(answer) < 1 AND  candidateid="+user.Id)//+ user.Id)
		rows2, _ := db.Query(stmt2)
		for rows2.Next() {
			rows2.Scan(&user.QuestionsAttended)
		}
		UsersInfo = append(UsersInfo, user)
	}
	fmt.Println(UsersInfo)
	t, _ := template.ParseFiles("./views/candidates.html")
	t.Execute(w, UsersInfo)
}

type ChallengeInfo struct {
	Id int
}

func newChallengeHandler(c web.C, w http.ResponseWriter, r *http.Request) {
	description := r.URL.Query().Get("desc")

	stmt2, _ := db.Prepare("insert into challenges (description, created) values($1, NOW())")
	stmt2.Query(description)

	http.Redirect(w, r, "programmingtest", 301)
}

func addChallengeHandler(c web.C, w http.ResponseWriter, r *http.Request) {

	http.Redirect(w, r, "programmingtest", 301)
}

func testcaseHandler(c web.C, w http.ResponseWriter, r *http.Request) {

	http.Redirect(w, r, "programmingtest", 301)
}

func main() {
	db = setupDB()
	defer db.Close()

	goji.Handle("/", candidateHandler)
	goji.Handle("/candidates", candidateHandler)
	goji.Get("/questions", questionsHandler)
	goji.Handle("/addQuestions", addQuestionsHandler)
	goji.Handle("/editquestion", editQuesionHandler)
	goji.Handle("/deleteQuestion", deleteQuestionHandler)
	goji.Handle("/deleteChallenges", deleteChanllengesHandler)
	goji.Handle("/addchallenge", addChallengeHandler)
	goji.get("/testcase", testcaseHandler)
	goji.Handle("/programmingtest", challengesHandler)
	goji.Handle("/editchallenge", editChallengeHandler)
	goji.Handle("/newChallenge", newChallengeHandler)
	http.Handle("/assets/css/", http.StripPrefix("/assets/css/", http.FileServer(http.Dir("assets/css"))))
	http.Handle("/assets/js/", http.StripPrefix("/assets/js/", http.FileServer(http.Dir("assets/js"))))
	http.Handle("/assets/img/", http.StripPrefix("/assets/img/", http.FileServer(http.Dir("assets/img"))))
	http.Handle("/assets/fonts/", http.StripPrefix("/assets/fonts/", http.FileServer(http.Dir("assets/fonts"))))
	goji.Serve()
}
