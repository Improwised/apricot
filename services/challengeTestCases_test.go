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
  _, err := exec.Command("../shell_scripts/challengeTestCases_test_before.sh").Output()
	services.CheckErr(err)
  // Run Tests
  exitCode := m.Run()

  // Exit
  os.Exit(exitCode)
}

func TestAddTestCases(t *testing.T) {
	services.AddTestCases("1", "2, 3", "4")
	services.AddTestCases("1", "4, 4", "8")

	var actualData string
	err2 := db.QueryRow("select max(input) from challenge_cases where id = 1").Scan(&actualData)
	services.CheckErr(err2)

	expectedData := "2, 3"
	if actualData != expectedData {
		t.Errorf("Test(AddTestCases) failed, expected input: '%s', got:  '%s'", expectedData, actualData)
	}
}

func TestEditTestCase(t *testing.T) {
	services.EditTestCase("1", "1", "3, 3", "6")
	var actualData string
	err2 := db.QueryRow("select max(input) from challenge_cases where id = 1").Scan(&actualData)
	services.CheckErr(err2)

	expectedData := "3, 3"
	if actualData != expectedData {
		t.Errorf("Test(EditTestCase) failed, expected input: '%s', got:  '%s'", expectedData, actualData)
	}
}

func TestSetDefaultTestCase(t *testing.T) {
	services.SetDefaultTestCase(2, 1)
	var actualData bool
	err2 := db.QueryRow("select defaultcase from challenge_cases where id = 2").Scan(&actualData)
	services.CheckErr(err2)

	expected := true
	if actualData != expected {
		t.Errorf("Test(SetDefaultTestCase) failed, testcase not seted default")
	}
}

func TestDisplayTestCaseForEdit(t *testing.T){
	services.DisplayTestCaseForEdit(1, 1)
	var actualData string
	err2 := db.QueryRow("select input from challenge_cases where id = 1 AND challengeid = 1").Scan(&actualData)
	services.CheckErr(err2)

	expected := "3, 3"
	if actualData != expected {
		t.Errorf("Test(DisplayTestCaseForEdit) failed, testcase not seted default")
	}
}

func TestDisplayTestCases(t *testing.T) {
	services.DisplayTestCases("1")
	testCases := services.ChallengeCases{}

	stmt, err := db.Prepare("select input from challenge_cases where challengeid = 1")
	services.CheckErr(err)
	rows, err2 := stmt.Query()
	services.CheckErr(err2)
	for rows.Next(){
		rows.Scan(&testCases.Input)
	}
	expected := "4, 4"
	if(testCases.Input != expected){
		t.Errorf("Test(DisplayTestCases) failed")
	}
}

func TestDeleteTestCase(t *testing.T) {
	services.DeleteTestCase(1, 1)
	stmt, err := db.Prepare("select input from challenge_cases where challengeid = 1 AND id = 1")
	services.CheckErr(err)
	rows, err2 := stmt.Query()
	services.CheckErr(err2)
	var input string
	for rows.Next(){
		rows.Scan(&input)
	}
	if(input != ""){
		t.Errorf("Test(DeleteTestCase) failed")
	}
}