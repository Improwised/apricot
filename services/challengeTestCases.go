
package services

import (
	"fmt"
	"encoding/json"
)

type ChallengeCases struct{
	Id int
	Input string
	Output string
	Default bool
	Challenge string
	Flag int
}

type ChallengeDesc struct{
	Challenge string
}

type AllDetails struct{
	ChallengeCases []ChallengeCases
	ChallengeDesc ChallengeDesc
}

func DisplayTestCases(challengeId string) AllDetails {
	Challenge := ChallengeDesc{}
	var encodedChallenge string
	err := db.QueryRow("select description from challenges where id = $1", challengeId).Scan(&encodedChallenge)
	CheckErr(err)

	//Decode the encrypted challenge from database...
	decoded := Decoding(encodedChallenge)

	Challenge.Challenge = decoded

	stmt1, err := db.Prepare("select id, input, output, defaultcase from challenge_cases where challengeid = ($1) order by id")
	CheckErr(err)
	rows1, err2 := stmt1.Query(challengeId)
	CheckErr(err2)

	challengeCases := []ChallengeCases{}
	q := ChallengeCases{}

	for rows1.Next() {
		err := rows1.Scan(&q.Id, &q.Input, &q.Output, &q.Default)
		CheckErr(err)
		if q.Default == true{
			q.Flag = 0
		} else if q.Default == false{
			q.Flag = 1
		}

		challengeCases = append(challengeCases, q)
	}
	allDetails := AllDetails{}
	allDetails.ChallengeCases = challengeCases
	allDetails.ChallengeDesc = Challenge

	return allDetails
}

func DisplayTestCaseForEdit(challengeId int, testCaseId int) string {
	stmt, err := db.Prepare("select input, output from challenge_cases where challengeId = ($1) AND id = ($2)")
	CheckErr(err)
	rows, err2 := stmt.Query(challengeId, testCaseId)
	CheckErr(err2)
	challengeCases := []ChallengeCases{}
	q := ChallengeCases{}
	for rows.Next() {
		err := rows.Scan(&q.Input, &q.Output)
		CheckErr(err)
		challengeCases = append(challengeCases, q)
	}
	b, err := json.Marshal(challengeCases)
	if err != nil {
		fmt.Printf("Error: %s", err)
	}
	return string(b)
}

func AddTestCases(challengeId string, input string, output string) {
	var defaultCase bool
	var tempVar string

	//check if very first test case of challenge ? if it is, than make it default testcase for challenge ..----
	err := db.QueryRow("select input from challenge_cases where challengeid = ($1)", challengeId).Scan(&tempVar)
	if(err == nil){
		defaultCase = false
	} else {
		defaultCase = true
	}
	// -----------------------------------------------
	stmt1, err := db.Prepare("insert into challenge_cases(challengeid, input, output, defaultCase, created) values ($1, $2, $3, $4, NOW());")
	CheckErr(err)
	stmt1.Query(challengeId, input, output, defaultCase)
}

func EditTestCase(testcaseid string, challengeid string, input string, output string) {
	stmt,err :=db.Prepare("update challenge_cases set input = ($1), output = ($2) where id = ($3) AND challengeId = ($4)")
	CheckErr(err)
	_,err2 := stmt.Query(input, output, testcaseid, challengeid)
	CheckErr(err2)
}

func DeleteTestCase(challengeId int, testCaseId int) string {
	var defaultStatus string
	err2 := db.QueryRow("select defaultcase from challenge_cases where id = ($1)", testCaseId).Scan(&defaultStatus)
	if(err2 != nil){
		CheckErr(err2)
	}

	if(defaultStatus == "false"){
		stmt, err :=db.Prepare("DELETE from challenge_cases WHERE challengeid = ($1) AND id = ($2)")
		CheckErr(err)
		_, err2 := stmt.Query(challengeId, testCaseId)
		CheckErr(err2)
	}

	//========to convert response to JSON ==========
	b, err := json.Marshal(defaultStatus)
	if err != nil {
		CheckErr(err)
	}
	return string(b)
}

func SetDefaultTestCase(testCaseId int, challengeId int) string {
	db.Query("update challenge_cases set defaultcase = false where challengeId = ($1)", challengeId)

	stmt1, err :=db.Prepare("UPDATE challenge_cases SET defaultcase = true WHERE id = ($1) AND challengeid = ($2)")
	CheckErr(err)
	_, errr := stmt1.Query(testCaseId, challengeId)
	CheckErr(errr)

	rows, err := db.Query("select id, defaultcase from challenge_cases where challengeid = ($1)", challengeId)
	CheckErr(err)
	challengeCases := []ChallengeCases{}
	q := ChallengeCases{}
	for rows.Next() {
		err := rows.Scan(&q.Id, &q.Default)
		CheckErr(err)
	}
	challengeCases = append(challengeCases, q)

	//========to convert response to JSON ==========
	b, err := json.Marshal(challengeCases)
	if err != nil {
		fmt.Printf("Error: %s", err)
	}
	//==========================
	return string(b)
}
