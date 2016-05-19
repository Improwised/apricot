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
	// "io/ioutil"
	"time"
	"strings"
	"bytes"
	// "strconv"
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
	stmt1 := fmt.Sprintf("SELECT c.id,c.name, c.email, c.degree, c.college, c.yearOfCompletion, c.modified, max(c1.attempts) FROM candidates c JOIN sessions s ON c.id = s.candidateid JOIN challenge_answers c1 ON s.id = c1.sessionid where s.status=0 group by c.id order by c.id asc ")
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

		stmt2 := fmt.Sprintf("SELECT count(id) FROM questions_answers WHERE length(answer) > 0 AND  candidateid="+user.Id)//+ user.Id)
		rows2, _ := db.Query(stmt2)
		for rows2.Next() {
			rows2.Scan(&user.QuestionsAttended)
		}
		UsersInfo = append(UsersInfo, user)
	}
	t, _ := template.ParseFiles("./views/candidates.html")
	t.Execute(w, UsersInfo)
}

type ChallengeInfo struct {
	Id int
}

func newChallengeHandler(c web.C, w http.ResponseWriter, r *http.Request) {
	description := r.FormValue("desc")

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

var id string
func personalInformationHandler(c web.C, w http.ResponseWriter, r *http.Request) {
	id = r.FormValue("id")

	QuestionsAttended := r.FormValue("queAttempt")
	ChallengeAttempts := r.FormValue("challengeAttmpt")

	stmt2 := fmt.Sprintf("SELECT name, email, contact, degree, college, yearofcompletion from candidates where id ="+id)
	rows3, _ := db.Query(stmt2)

	UsersInfo := []GeneralInfo{}
	user := GeneralInfo{}

	for rows3.Next() {
		rows3.Scan(&user.Name, &user.Email, &user.Contact, &user.Degree, &user.College, &user.YearOfCompletion)
		user.ChallengeAttempts = ChallengeAttempts
		user.QuestionsAttended = QuestionsAttended
		UsersInfo = append(UsersInfo, user)
	}
	t, _ := template.ParseFiles("./views/personalInformation.html")
	t.Execute(w, UsersInfo)
}

type GetQuestions struct {
	Questions string
	Id string
	Ans  string
}

func questionDetailsHandler(c web.C, w http.ResponseWriter, r *http.Request) {
	rows, _ := db.Query("SELECT questions.description,questions_answers.answer FROM questions INNER JOIN questions_answers ON questions.id=questions_answers.questionsid where candidateid = "+ id +" ORDER BY questions.sequence")
	questionsInfo := []GetQuestions{}
	qinfo := GetQuestions{}
	for rows.Next() {
		err := rows.Scan(&qinfo.Questions, &qinfo.Ans)
		questionsInfo = append(questionsInfo, qinfo)
		checkErr(err)
	}
	t, _ := template.ParseFiles("./views/questionDetails.html")
	t.Execute(w, questionsInfo)
}

type GetChallenge struct {
	Challenge string
}

type GetAnswers struct {
	Answer string
}

type AllDetail struct {
	GetChallenge []GetChallenge
	GetAnswers []GetAnswers
}

func challengeDetailsHandlers(c web.C, w http.ResponseWriter, r *http.Request) {

	stmt1, _ := db.Prepare("select description from challenges where id = ($1)")
	rows1, _ := stmt1.Query(id)
	challenge := []GetChallenge{}
	q := GetChallenge{}
	for rows1.Next() {
		err := rows1.Scan(&q.Challenge)
		challenge = append(challenge, q)
		checkErr(err)
	}
	stmt2, _ := db.Prepare("select answer from challenge_answers where sessionid = (select id from sessions where candidateid=($1)) order by attempts")
	rows2, _ := stmt2.Query(id)
	answer := []GetAnswers{}
	A := GetAnswers{}
	for rows2.Next() {
		err := rows2.Scan(&A.Answer)
		answer = append(answer, A)
		checkErr(err)
	}
	allDetails := AllDetail{}
	allDetails.GetChallenge = challenge
	allDetails.GetAnswers = answer

	t, _ := template.ParseFiles("./views/challengeDetails.html")
	t.Execute(w, allDetails)
}

var challengeId string
func addTestCase(c web.C, w http.ResponseWriter, r *http.Request) {

	if r.FormValue("qId") != "" {
		input := r.FormValue("input")
		output := r.FormValue("output")

		stmt1, _ := db.Prepare("insert into challenge_cases(challengeid, input, output, created) values ($1, $2, $3,NOW());")
		stmt1.Query(challengeId, input, output)
		http.Redirect(w, r, "programmingtest", 301)
	} else {
		qId := r.URL.Query().Get("qid")
		challengeId = qId

		stmt1, _ := db.Prepare("select description from challenges where id = ($1)")
		rows1, _ := stmt1.Query(qId)
		questions := []questionsInformation{}
		q := questionsInformation{}
		for rows1.Next() {
			err := rows1.Scan(&q.Description)
			questions = append(questions, q)
			checkErr(err)
		}
		t, _ := template.ParseFiles("./views/addTestCases.html")
		t.Execute(w, questions)
	}
}

func searchHandler(c web.C, w http.ResponseWriter, r *http.Request) {
	//data comes from admin side...
	name := r.FormValue("name")
	degree := r.FormValue("degree")
	college := r.FormValue("college")
	year := r.FormValue("year")

	// default query for search...
	var query ="SELECT c.id,c.name, c.email, c.degree, c.college, c.yearOfCompletion, c.modified, max(c1.attempts)"
	query += " FROM candidates c"
	query += " JOIN sessions s ON c.id = s.candidateid"
	query += " JOIN challenge_answers c1 ON s.id = c1.sessionid"
	query += " where s.status=0 "

	var stmt1 string

	// =======================making query for search =================================

if(year =="All"){//will search for all the year passing out candidates..
	if(name ==""){
		if(degree == "" && college == ""){//search for all the field..
			stmt1 = fmt.Sprintf(query+" group by c.id order by c.id asc ")
			} else if(degree == ""){//will search for college only..
				stmt1 = fmt.Sprintf(query+" AND (c.college ILIKE '%%%s%%')  group by c.id order by c.id asc ",college)
				}else if(college == ""){//will search for degree only..
					stmt1 = fmt.Sprintf(query+" AND (c.degree ILIKE '%%%s%%') group by c.id order by c.id asc ",degree)
					}

	} else if(degree == ""){
		 if(degree == "" && college == ""){//will search for name only..
				stmt1 = fmt.Sprintf(query+" AND ((c.name ILIKE '%%%s%%') OR (c.email LIKE '%%%s%%')) group by c.id order by c.id asc ",name,name)
				} else if(degree == ""){// will search for both name and college fields...
					stmt1 = fmt.Sprintf(query+" AND (((c.name ILIKE '%%%s%%') OR (c.email ILIKE '%%%s%%')) AND (c.college ILIKE '%%%s%%')) group by c.id order by c.id asc ",name,name,college)
						}

	} else if(college == ""){//will search for name and degree both field....
		stmt1 = fmt.Sprintf(query+" AND ((c.name ILIKE '%%%s%%') OR (c.email ILIKE '%%%s%%') AND (c.degree ILIKE '%%%s%%')) group by c.id order by c.id asc ",name,name,degree)
		} else {//will search for all the fields..
			stmt1 = fmt.Sprintf(query+" AND (((c.name ILIKE '%%%s%%') OR (c.email ILIKE '%%%s%%')) AND (c.college ILIKE '%%%s%%') AND (c.degree ILIKE '%%%s%%')) group by c.id order by c.id asc ",name,name,college,degree)
			}

} else {//will search for specific year passing out candidates..
	if(name ==""){
		if(degree == "" && college == ""){//search for all the field with specific year..
			stmt1 = fmt.Sprintf(query+" AND (c.yearOfCompletion::text LIKE '%%%s%%')group by c.id order by c.id asc ",year)
			} else if(degree == ""){//will search for college only with specific year..
				stmt1 = fmt.Sprintf(query+" AND ((c.college ILIKE '%%%s%%') AND (c.yearOfCompletion::text LIKE '%%%s%%')) group by c.id order by c.id asc ",college,year)
				}else if(college == ""){//will search for degree only with specific year..
					stmt1 = fmt.Sprintf(query+" AND ((c.degree ILIKE '%%%s%%') AND (c.yearOfCompletion::text LIKE '%%%s%%')) group by c.id order by c.id asc ",degree,year)
					}

	} else if(degree == ""){
		 if(degree == "" && college == ""){//will search for name only with specific year..
				stmt1 = fmt.Sprintf(query+" AND (((c.name ILIKE '%%%s%%') OR (c.email LIKE '%%%s%%')) AND (c.yearOfCompletion::text LIKE '%%%s%%')) group by c.id order by c.id asc ",name,name,year)
				} else if(degree == ""){// will search for both name and college fields with specific year...
					stmt1 = fmt.Sprintf(query+" AND (((c.name ILIKE '%%%s%%') OR (c.email ILIKE '%%%s%%')) AND (c.college ILIKE '%%%s%%') AND (c.yearOfCompletion::text LIKE '%%%s%%')) group by c.id order by c.id asc ",name,name,college,year)
						}

	} else if(college == ""){//will search for name and degree both field with specific year....
		stmt1 = fmt.Sprintf(query+" AND (((c.name ILIKE '%%%s%%') OR (c.email ILIKE '%%%s%%') AND (c.degree ILIKE '%%%s%%')) AND ((c.yearOfCompletion::text LIKE '%%%s%%')) group by c.id order by c.id asc ",name,name,degree,year)
		} else {//will search for all the fields with specific year..
			stmt1 = fmt.Sprintf(query+" AND (((c.name ILIKE '%%%s%%') OR (c.email ILIKE '%%%s%%')) AND (c.college ILIKE '%%%s%%') AND (c.degree ILIKE '%%%s%%') AND (c.yearOfCompletion::text LIKE '%%%s%%')) group by c.id order by c.id asc ",name,name,college,degree,year)
			}
}
	//==============================================================================================================================================

	rows1, err := db.Query(stmt1)
	if(err != nil){
		panic (err)
	}
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

		//=========will count no of attended questions========
		stmt2 := fmt.Sprintf("SELECT count(id) FROM questions_answers WHERE length(answer) > 0 AND  candidateid="+user.Id)
		rows2, _ := db.Query(stmt2)
		for rows2.Next() {
			rows2.Scan(&user.QuestionsAttended)
		}
		UsersInfo = append(UsersInfo, user)
	}
	//================================================

	//========to convert response to JSON ==========
	b, err := json.Marshal(UsersInfo)
	if err != nil {
			fmt.Printf("Error: %s", err)
			return;
	}
	//==========================

	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(b))//set response...
}

func main() {
	db = setupDB()
	defer db.Close()

	goji.Handle("/", candidateHandler)
	goji.Handle("/search", searchHandler)
	goji.Handle("/candidates", candidateHandler)
	goji.Get("/questions", questionsHandler)
	goji.Handle("/addQuestions", addQuestionsHandler)
	goji.Handle("/editquestion", editQuesionHandler)
	goji.Handle("/deleteQuestion", deleteQuestionHandler)
	goji.Handle("/deleteChallenges", deleteChanllengesHandler)
	goji.Handle("/personalInformation", personalInformationHandler)
	goji.Handle("/questionDetails", questionDetailsHandler)
	goji.Handle("/challengeDetails", challengeDetailsHandlers)
	goji.Handle("/addchallenge", addChallengeHandler)
	goji.Get("/testcase", testcaseHandler)
	goji.Handle("/addTestCases", addTestCase)
	goji.Handle("/programmingtest", challengesHandler)
	goji.Handle("/editchallenge", editChallengeHandler)
	goji.Handle("/newChallenge", newChallengeHandler)

	http.Handle("/assets/css/", http.StripPrefix("/assets/css/", http.FileServer(http.Dir("assets/css"))))
	http.Handle("/assets/jquery/", http.StripPrefix("/assets/jquery/", http.FileServer(http.Dir("assets/jquery"))))
	http.Handle("/assets/js/", http.StripPrefix("/assets/js/", http.FileServer(http.Dir("assets/js"))))
	http.Handle("/assets/img/", http.StripPrefix("/assets/img/", http.FileServer(http.Dir("assets/img"))))
	http.Handle("/assets/fonts/", http.StripPrefix("/assets/fonts/", http.FileServer(http.Dir("assets/fonts"))))
	goji.Serve()
}
