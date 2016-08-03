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
	_, err := exec.Command("../shell_scripts/question_test_before.sh").Output()
	services.CheckErr(err)
	// Run Tests
	exitCode := m.Run()

	// Exit
	os.Exit(exitCode)
}

// Add new question ...
func TestAddQuestion(t *testing.T){
	services.AddQuestion("This is Testing Question ...", 1)

	var actualQuestion string
	err := db.QueryRow("select description from questions where id = 1").Scan(&actualQuestion)
	services.CheckErr(err)
		expectedQuestion := "This is Testing Question ..."
	if(actualQuestion != expectedQuestion){
		t.Errorf("Test(AddQuestion) Failed, Not inserted data properly")
	}
}

//edit question ...
func TestEditQuestion(t *testing.T) {
	services.EditQuestion("This is Updated Questions with sequence ...", 1, 1)
	var actualQuestion string
	err := db.QueryRow("select description from questions where id = 1").Scan(&actualQuestion)
	services.CheckErr(err)
	expectedQuestion := "This is Updated Questions with sequence ..."
	if(actualQuestion != expectedQuestion){
		t.Errorf("Test(AddQuestion) Failed, Not inserted data properly")
	}
}

func TestDisplayQuestions(t *testing.T) {
	actual := services.DisplayQuestions(1)
	if actual == "" {
		t.Errorf("Test(DisplayQuestions) failed, expected: something, got:  '%s'", actual)
	}
	qInfo := services.GetQuestions{}
	stmt, err := db.Prepare("select id from questions")
	services.CheckErr(err)
	rows, err2 := stmt.Query()
	services.CheckErr(err2)
	for rows.Next(){
		rows.Scan(&qInfo.QuestionId)
	}
	expected := 1
	if(qInfo.QuestionId != expected){
		t.Errorf("Test(DisplayQuestions) failed, expected Question id : '%d', Got: '%d' ", expected, qInfo.QuestionId)
	}
}

func TestDeleteQuestion(t *testing.T) {
	actual := services.DeleteQuestion(1)
	expected := "yes"
	if(actual != expected){
		t.Errorf("Test(DeleteQuestion) Failed, expected : Question deleted, Found : not deleted")
	}
}

func TestFilterQuestions(t *testing.T) {
	actual := services.FilterQuestions("select id, description, deleted, sequence from questions where deleted is null order by sequence")

	if(len(actual.QuestionsInfo) != 0) {
		t.Errorf("Test(FilterQuestions) Failed, expected nothing filtered, Got some filtered Questions")
	}
}