package services

import (
	"os"
	"testing"
	"os/exec"
	"database/sql"

	"github.com/improwised/apricot/services"
)
var db *sql.DB = services.SetupDB()
func TestMain(m *testing.M) {
	// Setup
	_, err := exec.Command("../shell_scripts/candidate_test_before.sh").Output()
	services.CheckErr(err)

	// Run Tests
	exitCode := m.Run()

	// Exit
	os.Exit(exitCode)
}

func TestNewCandidateRegistration(t *testing.T) {
	expected := 40
	actual := services.NewCandidateRegistration("ashvin+test2@improwised.com")
	if len(actual) != expected {
		t.Errorf("Test(NewCandidateRegistration) failed, expected: '%d', got:  '%d'", expected, actual)
	}
	var hash string
	err := db.QueryRow("SELECT hash from sessions where candidateid = (select id from candidates where email = 'ashvin+test2@improwised.com')").Scan(&hash)
	services.CheckErr(err)
	services.AutoSave("ASHVIN", "name", hash)
	services.AutoSave("+91 9712186012", "contact", hash)
	services.AutoSave("BE", "degree", hash)
	services.AutoSave("RKU", "college", hash)
	services.AutoSave("2015", "yearOfCompletion", hash)
	services.AutoSave("Answer Of Question", "1", hash)

	var actualData string
	err2 := db.QueryRow("select name from candidates where email = 'ashvin+test2@improwised.com'").Scan(&actualData)
	services.CheckErr(err2)

	expectedData := "ASHVIN"
	if actualData != expectedData {
		t.Errorf("Test(NewCandidateRegistration) failed, expected Name: '%s', got:  '%s'", expectedData, actualData)
	}
}

func TestRegistreadActiveCandidate(t *testing.T) {
	expected := 40
	actual := services.RegistreadActiveCandidate(1, "ashvin+test2@improwised.com")
	if len(actual) != expected {
		t.Errorf("Test(RegistreadActiveCandidate) failed, expected: '%d', got:  '%d'", expected, len(actual))
	}
}

func TestRegistreadExpiredCandidate(t *testing.T) {
	expected := 40
	flag, actual := services.RegistreadExpiredCandidate(2, "ashvin+test@improwised.com")
	if len(actual) != expected && flag > 3 {
		t.Errorf("Test(RegistreadExpiredCandidate) failed, expected: '%d', got:  '%d'", expected, actual)
	}
}

func TestPassCandidatesInformation(t *testing.T) {
	actual := services.PassCandidatesInformation()

	expectedId := 1
	if(actual[0].Id != expectedId){
		t.Errorf("Test(PassCandidatesInformation) failed, expected Id : '%d', Got :'%d'", expectedId, actual[0].Id)
	}
	expectedName := "Ashvin"
	if(actual[0].Name != expectedName){
		t.Errorf("Test(PassCandidatesInformation) failed, expected Name : '%s', Got :'%s'", expectedName, actual[0].Name)
	}
	expectedEmail := "ashvin+test@improwised.com"
	if(actual[0].Email != expectedEmail){
		t.Errorf("Test(PassCandidatesInformation) failed, expected Email : '%s', Got :'%s'", expectedEmail, actual[0].Email)
	}
	expectedCollege := "RK University"
	if(actual[0].College != expectedCollege){
		t.Errorf("Test(PassCandidatesInformation) failed, expected College : '%s', Got :'%s'", expectedCollege, actual[0].College)
	}
	expectedDegree := "B.Tech"
	if(actual[0].Degree != expectedDegree){
		t.Errorf("Test(PassCandidatesInformation) failed, expected Degree : '%s', Got :'%s'", expectedDegree, actual[0].Degree)
	}
	expectedYear := "2016"
	if(actual[0].YearOfCompletion != expectedYear){
		t.Errorf("Test(PassCandidatesInformation) failed, expected Year Of Completion : '%s', Got :'%s'", expectedYear, actual[0].YearOfCompletion)
	}
	expectedChallengeAttempts := 0
	if(actual[0].ChallengeAttempts != expectedChallengeAttempts){
		t.Errorf("Test(PassCandidatesInformation) failed, expected Challenge Attempts : '%d', Got :'%d'", expectedChallengeAttempts, actual[0].ChallengeAttempts)
	}
	expectedQuestionAttempts := 1
	if(actual[0].QuestionsAttended != expectedQuestionAttempts){
		t.Errorf("Test(PassCandidatesInformation) failed, expected Questions Attended : '%d', Got :'%d'", expectedQuestionAttempts, actual[0].QuestionsAttended)
	}
}

func TestQuestionsAnswers(t *testing.T) {
	actual, _ := services.QuestionsAnswers(1)
	if(actual == nil){
		t.Errorf("Test(ChallengeDescription) failed, No data found in Database")
	}
	var actualAns string
	err := db.QueryRow("select answer from questions_answers where id = 1").Scan(&actualAns)
	services.CheckErr(err)

	expectedAns := "Answer of First Question"
	if(actualAns != expectedAns){
		t.Errorf("Test(ChallengeDescription) failed, expected Ans : '%s', Got : '%s'",expectedAns, actualAns)
	}
}

func TestChallengeDescription(t *testing.T) {
	actual := services.ChallengeDescription(1)
	if(actual == ""){
			t.Errorf("Test(ChallengeDescription) failed")
	}

	var ChallengeDescription string
	err := db.QueryRow("select description from challenges where id = 1").Scan(&ChallengeDescription)
	services.CheckErr(err)

	expectedDesc := "Rmlyc3QgVGVzdGluZyBDaGFsbGVuZ2U="
	if(ChallengeDescription != expectedDesc){
			t.Errorf("Test(ChallengeDescription) failed expected : '%s' Got : '%s'", expectedDesc, ChallengeDescription)
	}
}

func TestChallengeAnswerDetails(t *testing.T) {
	actual := services.ChallengeAnswerDetails(1)

	expectedChallenge := "First Testing Challenge"
	if(actual.GetChallenge.ChallengeDesc != expectedChallenge){
		t.Errorf("Test(ChallengeAnswerDetails) failed, Expected definition : '%s', Got : '%s'", expectedChallenge, actual.GetChallenge.ChallengeDesc)
	}

	expectedAnswer := "print 5"
	if(actual.LattestAnswer.Answer != expectedAnswer){
		t.Errorf("Test(ChallengeAnswerDetails) failed, Expected Source Code : '%s', Got : '%s'", expectedAnswer, actual.LattestAnswer.Answer)
	}
}

func TestCandidateInfo(t *testing.T) {
	actual := services.CandidateInfo(1)

	expectedName := "Ashvin"
	if(actual.Name != expectedName){
		t.Errorf("Test(CandidateInfo) failed, expected Name : '%s', Got :'%s'", expectedName, actual.Name)
	}
	expectedEmail := "ashvin+test@improwised.com"
	if(actual.Email != expectedEmail){
		t.Errorf("Test(CandidateInfo) failed, expected Email : '%s', Got :'%s'", expectedEmail, actual.Email)
	}
	expectedContact := "+91 9909970574"
	if(actual.Contact != expectedContact){
		t.Errorf("Test(CandidateInfo) failed, expected Contact : '%s', Got :'%s'", expectedName, actual.Contact)
	}
	expectedCollege := "RK University"
	if(actual.College != expectedCollege){
		t.Errorf("Test(CandidateInfo) failed, expected College : '%s', Got :'%s'", expectedCollege, actual.College)
	}
	expectedDegree := "B.Tech"
	if(actual.Degree != expectedDegree){
		t.Errorf("Test(CandidateInfo) failed, expected Degree : '%s', Got :'%s'", expectedDegree, actual.Degree)
	}
	expectedYear := "2016"
	if(actual.YearOfCompletion != expectedYear){
		t.Errorf("Test(CandidateInfo) failed, expected Year Of Completion : '%s', Got :'%s'", expectedYear, actual.YearOfCompletion)
	}
	expectedAttempts := 1
	if(actual.ChallengeAttempts != expectedAttempts){
		t.Errorf("Test(CandidateInfo) failed, expected Challenge Attempts : '%d', Got :'%d'", expectedAttempts, actual.ChallengeAttempts)
	}
}

func TestSearchQueryBuilder(t *testing.T) {
	expected := "SELECT c.id,c.name, c.email, c.degree, c.college, c.yearOfCompletion, c.modified, max(c1.attempts) FROM candidates c JOIN sessions s ON c.id = s.candidateid JOIN challenge_answers c1 ON s.id = c1.sessionid where s.status=0  AND (((c.name ILIKE '%a%') OR (c.email ILIKE '%a%')) AND (c.college ILIKE '%c%') AND (c.degree ILIKE '%b%')) group by c.id order by c.id asc "
	actual := services.SearchQueryBuilder("All", "a", "b", "c")
	if actual != expected {
		t.Errorf("Test(SearchQueryBuilder) failed, expected: '%s', got:  '%s'", expected, actual)
	}

	SearchResult := services.ReturnSearchResult(actual)
	if SearchResult == "" {
		t.Errorf("Test(ReturnSearchResult) failed")
	}
}

func TestPersonalInfo(t *testing.T) {
	actual := services.PersonalInfo(1)
	if(actual.Name == "" || actual.Email == "" || actual.College == "" || actual.Degree == ""|| actual.YearOfCompletion == ""){
		t.Errorf("Test(PersonalInfo) failed, no data found in database")
	}
	if(actual.Name != "Ashvin"){
		t.Errorf("Test(PersonalInfo) failed, expected Name : Ashvin Got : '%s'", actual.Name)
	}
	if(actual.Email != "ashvin+test@improwised.com"){
		t.Errorf("Test(PersonalInfo) failed, expected Email : ashvin+test@improwised.com Got :'%s'", actual.Email)
	}
	if(actual.Degree != "B.Tech"){
		t.Errorf("Test(PersonalInfo) failed, expected Degree : B.Tech Got :'%s'", actual.Degree)
	}
	if(actual.YearOfCompletion != "2016"){
		t.Errorf("Test(PersonalInfo) failed, expected Year Of Completion : 2016 Got:'%s'", actual.YearOfCompletion)
	}
}



