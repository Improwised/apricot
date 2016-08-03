package services

import (
		"os"
		"time"
		"os/exec"
		"testing"
		"database/sql"

		."github.com/improwised/apricot/services"
)

var db *sql.DB = services.SetupDB()
func TestMain(m *testing.M) {
	// Setup
	_, err := exec.Command("../shell_scripts/challenge_test_before.sh").Output()
	services.CheckErr(err)
	// Run Tests
	exitCode := m.Run()

	// Exit
	os.Exit(exitCode)
}

func TestAddChallenge(t *testing.T) {
	services.AddChallenge("Write a Programme to print sum of 2 no...")
	var actualChallenge string
	err := db.QueryRow("select description from challenges where id = 2").Scan(&actualChallenge)
	services.CheckErr(err)

	expectedChallenge := "V3JpdGUgYSBQcm9ncmFtbWUgdG8gcHJpbnQgc3VtIG9mIDIgbm8uLi4="
	if(actualChallenge != expectedChallenge){
		t.Errorf("Test(AddChallenge) Failed, expected challenge : '%s', Got : '%s' ",expectedChallenge, actualChallenge)
	}
}

func TestEditChallenge(t *testing.T) {
	services.EditChallenge("Write a Programme to print sum of n no...", 1)
	var actualChallenge string
	err := db.QueryRow("select description from challenges where id = 1").Scan(&actualChallenge)
	services.CheckErr(err)

	expectedChallenge := "V3JpdGUgYSBQcm9ncmFtbWUgdG8gcHJpbnQgc3VtIG9mIG4gbm8uLi4="
	if(actualChallenge != expectedChallenge){
		t.Errorf("Test(EditChallenge) Failed")
	}
}

func TestDisplayChallenges(t *testing.T) {
	services.DisplayChallenges(1)

	var actualChallenge string
	err := db.QueryRow("select description from challenges where id = 1").Scan(&actualChallenge)
	services.CheckErr(err)

	expectedChallenge := "V3JpdGUgYSBQcm9ncmFtbWUgdG8gcHJpbnQgc3VtIG9mIG4gbm8uLi4="
	if(actualChallenge != expectedChallenge){
		t.Errorf("Test(DisplayChallenges) Failed", actualChallenge)
	}
}

func TestDeleteChallenge(t *testing.T) {
	actual := services.DeleteChallenge(1)

	expected := "yes"
	if(actual != expected){
		t.Errorf("Test(DeleteChallenge) Failed, expected : challenge deleted, Found : not deleted")
	}
}

func TestFilterChallenges(t *testing.T) {
	actual := services.FilterChallenges("select id, description, deleted from challenges order by id")

	if(len(actual.QuestionsInfo) == 0) {
		t.Errorf("Test(FilterChallenges) Failed, expected some filtered challenges, Got nothing filtered")
	}
}

func TestChallengeDetails(t *testing.T) {
	//=========counting time for 7 days
	now := time.Now()
	sevenDay := time.Second * 585060 // Here 585060 is total seconds of 7 days ...
	time1 := now.Add(time.Duration(sevenDay))
	//=======
	stmt4, err6 := db.Prepare("UPDATE sessions set created = NOW(), expired = ($1), status = 1 where id = 1")
	services.CheckErr(err6)
	stmt4.Exec(time1)

	var hash string
	err := db.QueryRow("select hash from sessions where candidateid = 1").Scan(&hash)
	services.CheckErr(err)

	var expiredTime time.Time
	err2 := db.QueryRow("SELECT expired FROM sessions where hash = ($1)", hash).Scan(&expiredTime)
	services.CheckErr(err2)

	remainTime := expiredTime.Sub(time.Now())
	actual := services.ChallengeDetails(hash, remainTime)

	expectedDesc := "Write a Programme to print sum of n no..."
	if(expectedDesc != actual.GetChallenge[0].Description) {
		t.Errorf("Test(ChallengeDetails) Failed, expected : '%s', Got : '%s'", expectedDesc, actual.GetChallenge[0].Description)
	}

	expectedSource := "print 5"
	if(expectedSource != actual.GetSource[0].Answer) {
		t.Errorf("Test(ChallengeDetails) Failed, expected : '%s', Got : '%s'", expectedSource, actual.GetChallenge[0].Answer)
	}
}

func TestAttemptWiseSource(t *testing.T) {
	stmt4, err6 := db.Prepare("UPDATE sessions set status = 0 where id = ($1)")
	services.CheckErr(err6)
	stmt4.Exec(1)

	actual := services.AttemptWiseSource(1, 1)

	expectedSource := "print 5"
	if(actual[11:18] != expectedSource) {
		t.Errorf("Test(AttemptWiseSource) Failed",expectedSource, actual[11:18])
	}

	expectedLanguage := "python"
	if(actual[32:38] != expectedLanguage) {
		t.Errorf("Test(AttemptWiseSource) Failed",expectedLanguage, actual[32:38])
	}
}