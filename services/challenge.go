package services

import (
	"fmt"
	"time"
	"bytes"
	"encoding/json"
	"encoding/base64"
)

type getChallengeDescription struct {
	Description string
}

func DisplayChallenges(challengeId int) string {
	var encodedChallenge string
	err := db.QueryRow("select description from challenges where id = $1", challengeId).Scan(&encodedChallenge)
	CheckErr(err)

	decoded := Decoding(encodedChallenge)
	challenge := getChallengeDescription{}
	challenge.Description = decoded
	b, err := json.Marshal(challenge)
	if err != nil {
			fmt.Printf("Error: %s", err)
	}
	return string(b)
}

func AddChallenge(desc string) {
	//Encrypt the description to store in database with special charecters
	description := base64.StdEncoding.EncodeToString([]byte(desc))
	//==================================
	stmt2, err := db.Prepare("insert into challenges (description, created) values($1, NOW())")
	CheckErr(err)
	stmt2.Query(description)
}

func EditChallenge(description string, challengeId int) {
	//Encrypt the description to store in database with special charecters
	encodeChallenge := base64.StdEncoding.EncodeToString([]byte(description))
	//==================================

	stmt1, err := db.Prepare("update challenges set description = ($1) where id = ($2)")
	CheckErr(err)
	stmt1.Query(encodeChallenge, challengeId)
}

func DeleteChallenge(challengeId int) string {
	stmt, err := db.Prepare("select deleted from challenges where id = ($1)")
	CheckErr(err)
	rows, err2 := stmt.Query(challengeId)
	CheckErr(err2)
	q := questionsInformation{}
	var status string = "no"
	for rows.Next() {
		rows.Scan(&q.Deleted)
		if q.Deleted != nil{
			status = "no"
		} else if q.Deleted == nil{
			status = "yes"
		}
		if status == "yes" {
			stmt1, _ := db.Prepare("update challenges set deleted = NOW() where id = ($1)")
			stmt1.Query(challengeId)
		} else if status == "no" {
			stmt1, _ := db.Prepare("update challenges set deleted = NULL where id = ($1)")
			stmt1.Query(challengeId)
		}
	}
	return status
}

// get questions information
type challengeInformation struct {
	Id string
	Description string
	Flag int
	Deleted *time.Time
}

// link with child structure
type getAllChallengesInfo struct {
	QuestionsInfo []challengeInformation
}

func FilterChallenges(query string) getAllChallengesInfo {
	getAllChallengesInfo := getAllChallengesInfo{}
	var buffer bytes.Buffer
	buffer.WriteString(query)
	rows3, err := db.Query(buffer.String())
	CheckErr(err)

	questionsInfo := []challengeInformation{}
	q := challengeInformation{}
	var encodedChallenge string

	for rows3.Next() {
		err := rows3.Scan(&q.Id, &encodedChallenge, &q.Deleted)
		if q.Deleted != nil{
			q.Flag = 1
		} else if q.Deleted == nil{
			q.Flag = 0
		}
		//decode encoded challenge ...
		decoded := Decoding(encodedChallenge)

		// to print only few description of challenge rather than whole challenge on challenge list ..===
		if(len(decoded) > 150){
			//will display first 150 characters of challenge only if challenge have more than 150 char.
			decoded = Decoding(encodedChallenge)[0:150] + "  ...etc"
		}
		//=============================

		q.Description = decoded

		questionsInfo = append(questionsInfo, q)
		CheckErr(err)
	}
	getAllChallengesInfo.QuestionsInfo = questionsInfo
	return getAllChallengesInfo
}

func AttemptWiseSource(candidateId int, attemptNo int) string {
	attemptDetails := ChallengeAttemptsDetails{}
	stmt, err :=db.Prepare("SELECT answer, language FROM challenge_answers WHERE attempts=($1) AND sessionid = (select MAX(id) from sessions where candidateid=($2) AND status= 0 )")
	CheckErr(err)
	rows, err2 := stmt.Query(attemptNo, candidateId)
	CheckErr(err2)
	var encodedSource string

	for rows.Next() {
		err := rows.Scan(&encodedSource, &attemptDetails.Language)
		CheckErr(err)
	}
	decoded := Decoding(encodedSource)
	attemptDetails.Source = decoded

	//========to convert response to JSON ==========
	b, err := json.Marshal(attemptDetails)
	if err != nil {
		fmt.Printf("Error: %s", err)
	}
	//==========================
	return string(b)
}

type ChallengeAttemptsDetails struct{
	Source string
	Language string
}

type getChallenge struct {
	Description string
	Answer string
	RemainTime time.Duration
}

type getAllChallengeInfo struct {
	GetChallenge []getChallenge
	GetSource []getChallenge
	Hash string
}

func ChallengeDetails(hash string, remainTime time.Duration) getAllChallengeInfo{
	var sessionid int
	err := db.QueryRow("select id from sessions where hash='"+ hash +"'").Scan(&sessionid)
	CheckErr(err)

	getAllChallengeInfo := getAllChallengeInfo{}
	challengeInfo := []getChallenge{}
	cInfo := getChallenge{}

	var encodedChallenge string
	err3 := db.QueryRow("select description from challenges where id = (select challengeId from sessions where hash = ($1))", hash).Scan(&encodedChallenge)
	CheckErr(err3)

	decoded := Decoding(encodedChallenge)

	cInfo.Description = decoded
	//======================================================

	challengeInfo = append(challengeInfo, cInfo)

	//getting last answer for perticullar challege
	var encodeSource string
	var query1 = "select answer from challenge_answers"
			query1 += " where sessionid="
			query1 += " (select id from sessions where hash ='"+ hash +"')"
			query1 += " AND id=(select MAX(id) from challenge_answers"
			query1 += " where sessionid="
			query1 += " (select id from sessions"
			query1 += " where hash ='"+ hash +"'))"

	err4 := db.QueryRow(query1).Scan(&encodeSource)
	if(err4 != nil){
		encodeSource = "Ly8gV3JpdGUgWW91ciBDb2RlIEhlcmUgLi4="//ecoded form of "//write your code here .."
	}
	// db.Close()
	sourcecode := []getChallenge{}
	infoChallenge := getChallenge{}
	//===========

	//Decode the encrypted SourceCode from database...
	decoded2 := Decoding(encodeSource)
	//======================================================

	infoChallenge.Answer = decoded2
	infoChallenge.RemainTime = remainTime

	sourcecode = append(sourcecode, infoChallenge)

	getAllChallengeInfo.GetChallenge = challengeInfo
	getAllChallengeInfo.GetSource = sourcecode

	return getAllChallengeInfo
}
//**************************************************************************
