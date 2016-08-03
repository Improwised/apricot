package services

import (
		"os"
		"reflect"
		"os/exec"
		"testing"
		"database/sql"

		"github.com/improwised/apricot/services"
)

var db *sql.DB = services.SetupDB()
func TestMain(m *testing.M) {
	// Setup
	_, err := exec.Command("../shell_scripts/Functions_test_before.sh").Output()
	services.CheckErr(err)

	// Run Tests
	exitCode := m.Run()

	// Exit
	os.Exit(exitCode)
}

// setup db connection ...
func TestSetupDB(t *testing.T) {
	actual := services.SetupDB()
	if actual == nil {
		t.Errorf("Test(SetupDB) failed")
	}
}

func TestDecoding(t *testing.T) {
	expected := "test decoding"
	actual := services.Decoding("dGVzdCBkZWNvZGluZw==")
	if actual != expected {
		t.Errorf("Test(Decoding) failed, expected: '%s', got:  '%s'", expected, actual)
	}
}

func TestCompareString(t *testing.T) {
	expected := 0
	actual := services.CompareString("this is testing string comparision", "this is testing string comparision")
	if actual != expected {
		t.Errorf("Test(CompareString) failed, expected: '%d', got:  '%d'", expected, actual)
	}
}

func TestRandomChallengeGenerator(t *testing.T) {
	expected := "Some No. (ex. 1,2 etc.)"
	actual := services.RandomChallengeGenerator()
	if (actual <= 0)  {
		t.Errorf("Test5(RandomChallengeGenerator) failed, expected: '%s', got:  '%s'", expected, actual)
	}
}

func TestRandomNumber(t *testing.T) {
	var a int
	actual := services.RandomNumber(1, 10)
	if (actual > 10 || actual < 1 || reflect.TypeOf(actual) != reflect.TypeOf(a)){
		 t.Errorf("Test(RandomNumber) failed")
	}
}

func TestHashGenerator(t *testing.T){
	expected := 40
	actual := services.HashGenerator()
	if(len(actual) != 40){
		t.Errorf("Test(HashGenerator) failed, expected length: '%d', got:  '%d'", expected, len(actual))
	}
}

// func TestCalculateDayHrMin(t *testing.T) {
//   actual := services.CalculateDayHrMin(169)
//   expected := 0
//   if(actual.Days != expected){
//     t.Errorf("Test(CalculateDayHrMin) failed, expected day: '%d', got:  '%d'", expected, actual.Days)
//   }
// }

func TestSaveChallengeAnswer(t *testing.T){
	var hash string
	err := db.QueryRow("select hash from sessions where candidateid = 1").Scan(&hash)
	services.CheckErr(err)
	actualSource, actualid := services.SaveChallengeAnswer(hash, "print 5", "python")
	expected := 1
	if(actualid != expected){
		t.Errorf("Test(SaveChallengeAnswer) failed, expected id: '%d', got:  '%d'", expected, actualid)
	}
	expectedSource := "cHJpbnQgNQ=="
	if(actualSource != expectedSource){
		t.Errorf("Test(SaveChallengeAnswer) failed, expected Source: '%s', got:  '%s'", expected, actualid)
	}
}

func TestAutoSave(t *testing.T) {
	var hash string
	err := db.QueryRow("select hash from sessions where candidateid = 1").Scan(&hash)
	services.CheckErr(err)

	services.AutoSave("Ashvin", "name", hash)
	services.AutoSave("+91 9909970574", "contact", hash)
	services.AutoSave("B.Tech", "degree", hash)
	services.AutoSave("RK University", "college", hash)
	services.AutoSave("2016", "yearOfCompletion", hash)
	services.AutoSave("Answer Of Question", "1", hash)

	candidateInfo := services.GeneralInfo{}
	stmt, err := db.Prepare("select name, contact, degree, college, yearOfCompletion from candidates where id = 1")
	services.CheckErr(err)
	rows, err2 := stmt.Query()
	services.CheckErr(err2)

	for rows.Next(){
		rows.Scan(&candidateInfo.Name, &candidateInfo.Contact, &candidateInfo.Degree, &candidateInfo.College, &candidateInfo.YearOfCompletion)
	}
	if(candidateInfo.Name != "Ashvin"){
		t.Errorf("Test(AutoSave) failed, expected : Ashvin Got : '%s'", candidateInfo.Name)
	}
	if(candidateInfo.Contact != "+91 9909970574"){
		t.Errorf("Test(AutoSave) failed, expected : +91 9909970574 Got : '%s'", candidateInfo.Contact)
	}
	if(candidateInfo.Degree != "B.Tech"){
		t.Errorf("Test(AutoSave) failed, expected : B.Tech Got : '%s'", candidateInfo.Degree)
	}
	if(candidateInfo.College != "RK University"){
		t.Errorf("Test(AutoSave) failed, expected : RK University Got : '%s'", candidateInfo.College)
	}
	if(candidateInfo.YearOfCompletion != "2016"){
		t.Errorf("Test(AutoSave) failed, expected : 2016 : '%s'", candidateInfo.YearOfCompletion)
	}
}