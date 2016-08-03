package services

import (
	"fmt"
	"time"
	"encoding/json"
)

type GetQuestions struct {
	Questions string
	Ans  string
	Created time.Time
	DateTimeOnly string
	QuestionId int
	CandidateId int
	CandidateName string
}

func DisplayQuestions(questionId int) string {
	stmt1, err := db.Prepare("select id, description, sequence from questions where id = ($1)")
	CheckErr(err)
	rows1, err2 := stmt1.Query(questionId)
	CheckErr(err2)
	questions := []questionsInformation{}
	q := questionsInformation{}
	for rows1.Next() {
		err := rows1.Scan(&q.Id, &q.Description, &q.Sequence)
		questions = append(questions, q)
		CheckErr(err)
	}
	b, err := json.Marshal(questions)
	if err != nil {
		fmt.Printf("Error: %s", err)
		return "";
	}
	return string(b)
}

func AddQuestion(description string, sequence int) {
	stmt2, err := db.Prepare("insert into questions (description, sequence, created) values($1, $2, NOW())")
	CheckErr(err)
	stmt2.Query(description, sequence)
}

func EditQuestion(description string, sequence int, questionId int) {
	stmt1, err := db.Prepare("update questions set description = ($1), sequence = ($2) where id = ($3)")
	CheckErr(err)
	_, err2 := stmt1.Query(description, sequence, questionId)
	CheckErr(err2)
}

func DeleteQuestion(questionId int) string{
	stmt, err := db.Prepare("select deleted from questions where id = ($1)")
	CheckErr(err)
	rows, err2 := stmt.Query(questionId)
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
			stmt1, err3 := db.Prepare("update questions set deleted = NOW() where id = ($1)")
			CheckErr(err3)
			stmt1.Query(questionId)
		} else if status == "no" {
			stmt1, err4 := db.Prepare("update questions set deleted = NULL where id = ($1)")
			CheckErr(err4)
			stmt1.Query(questionId)
		}
	}
	return status
}

// get questions information
type questionsInformation struct {
	Id int
	Description string
	Sequence int
	Flag int
	Deleted *time.Time
}

// link with child structure
type getAllQuestionsInfo struct {
	QuestionsInfo []questionsInformation
}

func FilterQuestions(query string) getAllQuestionsInfo {
	rows3, err := db.Query(query)
	CheckErr(err)
	getAllQuestionsInfo := getAllQuestionsInfo{}
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
		CheckErr(err)
	}
	getAllQuestionsInfo.QuestionsInfo = questionsInfo

	return getAllQuestionsInfo
}
