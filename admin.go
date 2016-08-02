package main

import (
	"strconv"
	"net/http"
	"database/sql"
	"html/template"

	_"github.com/lib/pq"
	"github.com/zenazn/goji"
	"github.com/zenazn/goji/web"
	"github.com/improwised/apricot/services"
)

// db connection
var db *sql.DB

// get question information from database using qId
func getQuestionInfoHandler(c web.C, w http.ResponseWriter, r *http.Request) {
	qId, err := strconv.Atoi( r.FormValue("id"))
	services.CheckErr(err)
	b := services.DisplayQuestions(qId)
	//==========================
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(b))//set response...
}

// display only active questions in view...
func allQuestionsHandler(c web.C, w http.ResponseWriter, r *http.Request) {
	getAllQuestionsInfo := services.FilterQuestions("select id, description, deleted, sequence from questions  order by sequence")
	t, err := template.ParseFiles("./views/admin/questions.html")
	services.CheckErr(err)
	t.Execute(w, getAllQuestionsInfo)
}

// display all questions in view...
func questionsHandler(c web.C, w http.ResponseWriter, r *http.Request) {
	getAllQuestionsInfo := services.FilterQuestions("select id, description, deleted, sequence from questions where deleted is null order by sequence")
	t, err := template.ParseFiles("./views/admin/questions.html")
	services.CheckErr(err)
	t.Execute(w, getAllQuestionsInfo)
}

// perform edit functionality
func editQuesionHandler(c web.C, w http.ResponseWriter, r *http.Request) {
	description := r.FormValue("description")
	sequence, err := strconv.Atoi(r.FormValue("sequence"))
	services.CheckErr(err)
	qId, err2 := strconv.Atoi(r.FormValue("qId"))
	services.CheckErr(err2)
	services.EditQuestion(description, sequence, qId)
	http.Redirect(w, r, "questions", 301)
}

//  delete questions functionality
func deleteQuestionHandler(c web.C, w http.ResponseWriter, r *http.Request) {
	qId, err := strconv.Atoi(r.FormValue("qid"))
	services.CheckErr(err)

	status := services.DeleteQuestion(qId)
	w.Write([]byte(status))
}

// add questions functionality
func addQuestionsHandler(c web.C, w http.ResponseWriter, r *http.Request) {
	description := r.FormValue("description")
	sequence, err := strconv.Atoi(r.FormValue("sequence"))
	services.CheckErr(err)

	services.AddQuestion(description, sequence)
	http.Redirect(w, r, "questions", 301)
}

func allChallengesHandler(c web.C, w http.ResponseWriter, r *http.Request) {
	getAllQuestionsInfo := services.FilterChallenges("select id, description, deleted from challenges order by id")
	t, err := template.ParseFiles("./views/admin/challenges.html")
	services.CheckErr(err)
	t.Execute(w, getAllQuestionsInfo)
}

// retrive challenges from database.
func challengesHandler(c web.C, w http.ResponseWriter, r *http.Request) {
	getAllQuestionsInfo := services.FilterChallenges("select id, description, deleted from challenges where deleted is null order by id")
	t, err := template.ParseFiles("./views/admin/challenges.html")
	services.CheckErr(err)
	t.Execute(w, getAllQuestionsInfo)
}

// mark challenge as a deleted.
func deleteChanllengesHandler(c web.C, w http.ResponseWriter, r *http.Request) {
	cId, err := strconv.Atoi(r.FormValue("qid"))
	services.CheckErr(err)
	status := services.DeleteChallenge(cId)
	w.Write([]byte(status))
}

// edit perticualr challenge.
func editChallengeHandler(c web.C, w http.ResponseWriter, r *http.Request) {
	challengeId, err := strconv.Atoi(r.FormValue("challengeId"))
	services.CheckErr(err)
	description := r.FormValue("challengeDescription")
	services.EditChallenge(description, challengeId)
	http.Redirect(w, r, "challenges", 301)
}

// display candidates information
func candidateHandler(c web.C, w http.ResponseWriter, r *http.Request) {
	UsersInfo := services.PassCandidatesInformation()
	t, err := template.ParseFiles("./views/admin/candidates.html")
	services.CheckErr(err)
	t.Execute(w, UsersInfo)
}

// add new programming challenges.
func newChallengeHandler(c web.C, w http.ResponseWriter, r *http.Request) {
	desc := r.FormValue("description")
	services.AddChallenge(desc)
	http.Redirect(w, r, "challenges", 301)
}

func addChallengeHandler(c web.C, w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "challenges", 301)
}

func testcaseHandler(c web.C, w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "challenges", 301)
}

//will display personal information of candidates..
func personalInformationHandler(c web.C, w http.ResponseWriter, r *http.Request) {
	candidateId, err := strconv.Atoi(r.URL.Query().Get("id"))
	services.CheckErr(err)

	UsersInfo := services.CandidateInfo(candidateId)//personal information of candidate ..
	t, err := template.ParseFiles("./views/admin/personalInformation.html")
	services.CheckErr(err)
	t.Execute(w, UsersInfo)
}

type AllQueDetails struct{
	GetQuestions []services.GetQuestions
	CandidateId int
	CandidateName string
}

//will display questions and answer given by candidates..
func questionDetailsHandler(c web.C, w http.ResponseWriter, r *http.Request) {
	candidateId, err := strconv.Atoi(r.URL.Query().Get("id"))
	services.CheckErr(err)
	Details, CandidateName := services.QuestionsAnswers(candidateId)

	AllQueDetails := AllQueDetails{}
	AllQueDetails.GetQuestions = Details
	AllQueDetails.CandidateId = candidateId
	AllQueDetails.CandidateName = CandidateName

	t, err := template.ParseFiles("./views/admin/questionDetails.html")
	services.CheckErr(err)
	t.Execute(w, AllQueDetails)
}

//will display answer of challenge...
func challengeDetailsHandlers(c web.C, w http.ResponseWriter, r *http.Request) {
	candidateId, err := strconv.Atoi(r.URL.Query().Get("id"))
	services.CheckErr(err)
	challengeAnswer := services.ChallengeAnswerDetails(candidateId)
	t, err := template.ParseFiles("./views/admin/challengeDetails.html")
	services.CheckErr(err)
	t.Execute(w, challengeAnswer)
}

// add testcases for perticular challenge.
var challengeId string
func addTestCase(c web.C, w http.ResponseWriter, r *http.Request) {
	var qId string = r.URL.Query().Get("qid")

	if qId != "" {
		challengeId = qId
		allDetails := services.DisplayTestCases(challengeId)
		t, err := template.ParseFiles("./views/admin/addTestCases.html")
		services.CheckErr(err)
		t.Execute(w, allDetails)
	} else {
		input := r.FormValue("input")
		output := r.FormValue("output")
		services.AddTestCases(challengeId, input, output)
		http.Redirect(w, r, "addTestCases?qid=" + challengeId , 301)
	}
}

// retrive testcase input output
func getTestCaseHandler(c web.C, w http.ResponseWriter, r *http.Request) {
	testCaseId, err := strconv.Atoi(r.FormValue("testCaseId"))
	services.CheckErr(err)
	challengeId, err2 := strconv.Atoi(r.FormValue("challengeId"))
	services.CheckErr(err2)
	b := services.DisplayTestCaseForEdit(challengeId, testCaseId)
	//==========================
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(b))//set response...
}

// edit test case .
func editTestCaseHandler(c web.C, w http.ResponseWriter, r *http.Request) {
	testCaseId := r.FormValue("testCaseId")
	challengeId := r.FormValue("challengeId")
	input := r.FormValue("input")
	output := r.FormValue("output")
	services.EditTestCase(testCaseId, challengeId, input, output)
	http.Redirect(w, r,  "addTestCases?qid=" + challengeId , 301)
}

//Searching will perform on candidates information...
func searchHandler(c web.C, w http.ResponseWriter, r *http.Request) {
	//data comes from UI ...
	name := r.FormValue("name")
	degree := r.FormValue("degree")
	college := r.FormValue("college")
	year := r.FormValue("year")

	//pass data to function to make query for search ...
	stmt1 := services.SearchQueryBuilder(year, name ,degree, college)
	//pass prepared query to function to get searched result ...
	b := services.ReturnSearchResult(stmt1)
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(b))//set response...
}

//will return the perticullar attempted challenge source code
func challengeAttemptHandler(c web.C, w http.ResponseWriter, r *http.Request) {
	candidateId, err := strconv.Atoi(r.FormValue("candidateID"))
	services.CheckErr(err)
	attemptNo, err2 := strconv.Atoi(r.FormValue("attemptNo"))
	services.CheckErr(err2)
	b := services.AttemptWiseSource(candidateId, attemptNo)
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(b))//set response...
}

func deleteTestcaseHandler(c web.C, w http.ResponseWriter, r *http.Request) {
	challengeId, err := strconv.Atoi(r.FormValue("challengeId"))
	services.CheckErr(err)
	testCaseId, err2 := strconv.Atoi(r.FormValue("testCaseId"))
	services.CheckErr(err2)
	b := services.DeleteTestCase(challengeId, testCaseId)
	//==========================
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(b))//set response...
}

func DefaultTestcaseHandler(c web.C, w http.ResponseWriter, r *http.Request) {
	challengeId, err := strconv.Atoi(r.FormValue("challengeId"))
	services.CheckErr(err)
	testCaseId, err2 := strconv.Atoi(r.FormValue("testCaseId"))
	services.CheckErr(err2)
	b := services.SetDefaultTestCase(testCaseId, challengeId)
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(b))//set response...
}

// get challenge information using challenge id.
func getChallengeInfoHandler(c web.C, w http.ResponseWriter, r *http.Request) {
	challengeId, err := strconv.Atoi(r.FormValue("challengeId"))
	services.CheckErr(err)
	b := services.DisplayChallenges(challengeId)
	//==========================
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(b))//set response...
}

func main() {
	db = services.SetupDB()
	defer db.Close()

	goji.Get("/", candidateHandler)
	goji.Get("/allQuestions", allQuestionsHandler)
	goji.Get("/addchallenge", addChallengeHandler)
	goji.Get("/allChallenges", allChallengesHandler)
	goji.Post("/addQuestions", addQuestionsHandler)

	goji.Get("/candidates", candidateHandler)
	goji.Get("/candidate/challengeDetails", challengeDetailsHandlers)
	goji.Get("/candidate/personalInformation", personalInformationHandler)
	goji.Get("/candidate/questionDetails", questionDetailsHandler)
	goji.Post("/candidate/challengeAttempt", challengeAttemptHandler)

	goji.Get("/challenges", challengesHandler)
	goji.Handle("/challenges/addTestCases", addTestCase)
	goji.Post("/challenges/deleteTestCase", deleteTestcaseHandler)
	goji.Post("/challenges/editTestCase", editTestCaseHandler)
	goji.Post("/challenges/getTestCase", getTestCaseHandler)
	goji.Post("/challenges/setDefaultTestcase", DefaultTestcaseHandler)

	goji.Post("/deleteQuestion", deleteQuestionHandler)
	goji.Post("/deleteChallenge", deleteChanllengesHandler)

	goji.Post("/editchallenge", editChallengeHandler)
	goji.Post("/editquestion", editQuesionHandler)

	goji.Post("/search", searchHandler)

	goji.Post("/getQuestionInfo", getQuestionInfoHandler);
	goji.Post("/getChallengeInfo", getChallengeInfoHandler)
	goji.Post("/newChallenge", newChallengeHandler)

	goji.Get("/questions", questionsHandler)
	goji.Get("/testcase", testcaseHandler)

	http.Handle("/assets/css/", http.StripPrefix("/assets/css/", http.FileServer(http.Dir("assets/css"))))
	http.Handle("/assets/jquery/", http.StripPrefix("/assets/jquery/", http.FileServer(http.Dir("assets/jquery"))))
	http.Handle("/assets/js/", http.StripPrefix("/assets/js/", http.FileServer(http.Dir("assets/js"))))
	http.Handle("/assets/img/", http.StripPrefix("/assets/img/", http.FileServer(http.Dir("assets/img"))))
	http.Handle("/assets/fonts/", http.StripPrefix("/assets/fonts/", http.FileServer(http.Dir("assets/fonts"))))
	goji.Serve()
}
