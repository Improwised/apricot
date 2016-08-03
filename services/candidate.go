package services

import (
	"fmt"
	"time"
	"encoding/json"
)

type GeneralInfo struct {
	Id int
	Name string
	Contact string
	Degree string
	College string
	YearOfCompletion string
	Email string
	Created time.Time
	Modified time.Time
	ChallengeAttempts int
	DateOnly string
	QuestionsAttended int
	NoOfTry int
}

func PassCandidatesInformation() []GeneralInfo{
var query = "SELECT c.id, c.name, c.email, c.degree, c.college, c.yearOfCompletion, c.modified, max(c1.attempts)"
		query += " FROM candidates c JOIN sessions s ON c.id = s.candidateid"
		query += " JOIN challenge_answers c1 ON s.id = c1.sessionid"
		query += " where s.status = 0"
		query += " group by c.id"
		query += " order by c.id asc"

	stmt1 := fmt.Sprintf(query)
	rows1, err := db.Query(stmt1)
	CheckErr(err)

	UsersInfo := []GeneralInfo{}
	user := GeneralInfo{}
	for rows1.Next() {
		rows1.Scan(&user.Id, &user.Name, &user.Email, &user.Degree, &user.College, &user.YearOfCompletion, &user.Modified, &user.ChallengeAttempts)

		//extract only date from timestamp========
		t := &user.Modified
		user.DateOnly = t.Format("2006-01-02")
		//================================

		query := "SELECT count(id) FROM questions_answers WHERE length(answer) > 0 AND  candidateid = $1"
		stmt, err := db.Prepare(query)
		CheckErr(err)
		rows2, err2 := stmt.Query(user.Id)
		CheckErr(err2)
		for rows2.Next() {
			rows2.Scan(&user.QuestionsAttended)
		}
		UsersInfo = append(UsersInfo, user)
	}
	return UsersInfo
}

func QuestionsAnswers(id int) ([]GetQuestions, string) {
	var query = "SELECT max(questions.id), questions.description, max(questions_answers.answer), max(questions_answers.Created)"
			query += " FROM questions"
			query += " INNER JOIN questions_answers"
			query += " ON questions.id = questions_answers.questionsid"
			query += " where candidateid = ($1) "
			query += " group BY questions.ID"
			query += " ORDER BY questions.sequence"

	stmt, err := db.Prepare(query)
	CheckErr(err)
	rows, err2 := stmt.Query(id)
	CheckErr(err2)

	questionsInfo := []GetQuestions{}
	qinfo := GetQuestions{}

	for rows.Next() {
		err := rows.Scan(&qinfo.QuestionId, &qinfo.Questions, &qinfo.Ans, &qinfo.Created)
		// extract only date - time from timestamp========
		t :=&qinfo.Created
		qinfo.DateTimeOnly = t.Format("02/06/2016 03:04:05 PM")
		//================================
		questionsInfo = append(questionsInfo, qinfo)
		CheckErr(err)
	}
	var candidateName string
	err3 := db.QueryRow("SELECT DISTINCT name FROM candidates WHERE id = ($1)", id).Scan(&candidateName)
	if(err3 != nil){
		candidateName = ""
	}

	return questionsInfo, candidateName
}

func ChallengeDescription(id int) string{
	var encodedChallenge string

	query := "select description from challenges where id = (select challengeid from sessions where candidateid = ($1) AND id =(select MAX(id) from sessions where candidateid = ($2)))"
	stmt2, err := db.Prepare(query)
	CheckErr(err)
	rows2, err2 := stmt2.Query(id, id)
	CheckErr(err2)
	for rows2.Next() {
		rows2.Scan(&encodedChallenge)
	}

	//decode the challnge ...=====
	decoded := Decoding(encodedChallenge)
	// ==============================

	return decoded
}

type ChallengeDefi struct{
	CandidateId int
	ChallengeDesc string
	CandidateName string
}

type GetAnswers struct {
	Answer string
	Attempt string
	LattestAnswer string
	Lang string
}

type TestCases struct {
	AllTestCases int
	PassedTestCases int
}

type AllDetail struct {
	GetChallenge ChallengeDefi
	GetAnswers []GetAnswers
	LattestAnswer GetAnswers
	TestCases TestCases
	Lang string
	CandidateName string
}

func ChallengeAnswerDetails(id int) AllDetail {
	var candidateName string
	err5 := db.QueryRow("SELECT DISTINCT name FROM candidates WHERE id = ($1)", id).Scan(&candidateName)
	CheckErr(err5)

	challengeInfo := ChallengeDefi{}
	challengeInfo.CandidateName = candidateName
	challengeInfo.CandidateId = id
	challengeInfo.ChallengeDesc = ChallengeDescription(id)

	stmt2, err := db.Prepare("select answer, attempts, language from challenge_answers where sessionid = (select max(id) from sessions where candidateid=($1) AND status = 0) order by attempts")
	CheckErr(err)
	rows2, err2 := stmt2.Query(id)
	CheckErr(err2)
	answer := []GetAnswers{}

	A := GetAnswers{}
	B := GetAnswers{}

	var encodedSource string
	for rows2.Next() {
		err := rows2.Scan(&encodedSource, &A.Attempt, &A.Lang)
		answer = append(answer, A)
		CheckErr(err)
	}
	//Decode the encrypted source from database...
	decoded2 := Decoding(encodedSource)
	//==================================================

	B.Answer = decoded2
	testcases := TestCases{}

	err3 := db.QueryRow("select count(input) from challenge_attempts where sessionid =(select max(id) from sessions where candidateId = ($1) AND status = 0)",id).Scan(&testcases.AllTestCases)
	CheckErr(err3)

	err4 := db.QueryRow("select count(input) from challenge_attempts where sessionid =(select max(id) from sessions where candidateId = ($1) AND status = 0) AND status = 1", id).Scan(&testcases.PassedTestCases)
	CheckErr(err4)

	allDetails := AllDetail{}
	allDetails.GetChallenge = challengeInfo
	allDetails.GetAnswers = answer
	allDetails.LattestAnswer = B
	allDetails.TestCases = testcases
	allDetails.Lang = A.Lang

	return allDetails
}

// candiate informations for admin view ...
func CandidateInfo(Id int) GeneralInfo{
	var query = "SELECT c.name, c.email, c.Contact, c.degree, c.college, c.yearOfCompletion, max(c1.attempts)"
			query += " FROM candidates c JOIN sessions s ON c.id = s.candidateid"
			query += " JOIN challenge_answers c1 ON s.id = c1.sessionid"
			query += " where s.status = 0 AND c.Id = ($1)"
			query += " group by c.id"
			query += " order by c.id asc"

	stmt, err := db.Prepare(query)
	CheckErr(err)
	rows, err2 := stmt.Query(Id)
	CheckErr(err2)

	user := GeneralInfo{}
	for rows.Next() {
		rows.Scan(&user.Name, &user.Email, &user.Contact, &user.Degree, &user.College, &user.YearOfCompletion, &user.ChallengeAttempts)
		db.QueryRow("SELECT count(id) from sessions where status = 0 AND candidateid = ($1)", Id).Scan(&user.NoOfTry)
	}
	var noOfQueAttended int
	err3 := db.QueryRow("SELECT count(DISTINCT questionsid) FROM questions_answers WHERE length(answer) > 0 AND  candidateid = ($1)", Id).Scan(&noOfQueAttended)
	CheckErr(err3)

	user.QuestionsAttended = noOfQueAttended
	user.Id = Id

	return user
}

// candidate information for client view ..
func PersonalInfo(Id int) GeneralInfo{
	var query = "SELECT name, email, Contact, degree, college, yearOfCompletion FROM candidates where id =($1)"
	stmt, err := db.Prepare(query)
	CheckErr(err)
	rows, err2 := stmt.Query(Id)
	CheckErr(err2)

	user := GeneralInfo{}
	for rows.Next() {
		rows.Scan(&user.Name, &user.Email, &user.Contact, &user.Degree, &user.College, &user.YearOfCompletion)
	}
	return user
}

//function will make query for search according to parameters ....
func SearchQueryBuilder(year string, name string, degree string, college string) string {
	// default query for search...
	var query ="SELECT c.id,c.name, c.email, c.degree, c.college, c.yearOfCompletion, c.modified, max(c1.attempts)"
		query += " FROM candidates c"
		query += " JOIN sessions s ON c.id = s.candidateid"
		query += " JOIN challenge_answers c1 ON s.id = c1.sessionid"
		query += " where s.status=0 "

	var stmt1 string
	// ======================= making query for search =================================

	if(year =="All"){//will search for all the year passing out candidates..
		if(name ==""){
			if(degree == "" && college == ""){//search for all the field..
				stmt1 = fmt.Sprintf(query+" group by c.id order by c.id asc ")
				} else if(degree == ""){//will search for college only..
					stmt1 = fmt.Sprintf(query+" AND (c.college ILIKE '%%%s%%')  group by c.id order by c.id asc ",college)
					}else if(college == ""){//will search for degree only..
						stmt1 = fmt.Sprintf(query+" AND (c.degree ILIKE '%%%s%%') group by c.id order by c.id asc ",degree)
						}else{
							stmt1 = fmt.Sprintf(query+" AND ((c.degree ILIKE '%%%s%%') AND (c.college ILIKE '%%%s%%') ) group by c.id order by c.id asc ",degree,college)
						}

		} else if(degree == ""){
			 if(degree == "" && college == ""){//will search for name only..
					stmt1 = fmt.Sprintf(query+" AND ((c.name ILIKE '%%%s%%') OR (c.email LIKE '%%%s%%')) group by c.id order by c.id asc ",name,name)
					} else if(degree == ""){// will search for both name and college fields...
						stmt1 = fmt.Sprintf(query+" AND (((c.name ILIKE '%%%s%%') OR (c.email ILIKE '%%%s%%')) AND (c.college ILIKE '%%%s%%')) group by c.id order by c.id asc ",name,name,college)
							}

		} else if(college == ""){//will search for name and degree both field....
			stmt1 = fmt.Sprintf(query+" AND (((c.name ILIKE '%%%s%%') OR (c.email ILIKE '%%%s%%')) AND (c.degree ILIKE '%%%s%%')) group by c.id order by c.id asc ",name,name,degree)
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
						}else{//will search for all the fields excepting name/email..
							stmt1 = fmt.Sprintf(query+" AND ((c.degree ILIKE '%%%s%%') AND (c.college ILIKE '%%%s%%') AND (c.yearOfCompletion::text LIKE '%%%s%%')) group by c.id order by c.id asc ",degree,college,year)
						}

		} else if(degree == ""){
			if(degree == "" && college == ""){//will search for name only with specific year..
				stmt1 = fmt.Sprintf(query+" AND (((c.name ILIKE '%%%s%%') OR (c.email LIKE '%%%s%%')) AND (c.yearOfCompletion::text LIKE '%%%s%%')) group by c.id order by c.id asc ",name,name,year)
				} else if(degree == ""){// will search for both name and college fields with specific year...
					stmt1 = fmt.Sprintf(query+" AND (((c.name ILIKE '%%%s%%') OR (c.email ILIKE '%%%s%%')) AND (c.college ILIKE '%%%s%%') AND (c.yearOfCompletion::text LIKE '%%%s%%')) group by c.id order by c.id asc ",name,name,college,year)
						}

		} else if(college == ""){//will search for name and degree both field with specific year....
			stmt1 = fmt.Sprintf(query+" AND (((c.name ILIKE '%%%s%%') OR (c.email ILIKE '%%%s%%')) AND (c.degree ILIKE '%%%s%%') AND (c.yearOfCompletion::text LIKE '%%%s%%')) group by c.id order by c.id asc ",name,name,degree,year)
			} else {//will search for all the fields with specific year..
				stmt1 = fmt.Sprintf(query+" AND (((c.name ILIKE '%%%s%%') OR (c.email ILIKE '%%%s%%')) AND (c.college ILIKE '%%%s%%') AND (c.degree ILIKE '%%%s%%') AND (c.yearOfCompletion::text LIKE '%%%s%%')) group by c.id order by c.id asc ",name,name,college,degree,year)
				}
	}
	//==============================================================================================================================================
	return stmt1
}

func ReturnSearchResult(stmt1 string) string {
	rows1, err := db.Query(stmt1)
	if(err != nil){
		panic (err)
	}
	UsersInfo := []GeneralInfo{}
	user := GeneralInfo{}
	for rows1.Next() {

		rows1.Scan(&user.Id, &user.Name, &user.Email, &user.Degree, &user.College, &user.YearOfCompletion, &user.Modified, &user.ChallengeAttempts)

		//extract only date from timestamp========
		t := &user.Modified
		user.DateOnly = t.Format("2006-01-02")
		//================================

		//=========will count no of attended questions========
		query := "SELECT count(id) FROM questions_answers WHERE length(answer) > 0 AND  candidateid = $1"
		stmt2, err3 := db.Prepare(query)
		CheckErr(err3)
		rows2, err2 := stmt2.Query(user.Id)
		CheckErr(err2)
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
	}
	//==========================
	return string(b)
}

type Questions struct {
	Qid string
}

type sessionInfo struct {
	candidateid string
	sessionid string
	hash string
	status int
	attempts int
	entryDate time.Time
	modifyDate time.Time
	expireDate time.Time
	challenge int
}

//if candidate already registread and link expired...
func RegistreadExpiredCandidate(candidateId int, email string) (int, string){
	var status int
	var flag int = 2
	err := db.QueryRow("SELECT status FROM sessions WHERE candidateid = $1 AND id = (select MAX(id) from sessions where candidateid = ($2))", candidateId, candidateId).Scan(&status)
	if err != nil{
		panic(err)
	}
	var hash string
	if(status == 0){
		hash = HashGenerator() // will generate unique hash...

		//===== will return random challenge among no of challenges =========
		challengeNo := RandomChallengeGenerator()
		//======================================================================
		stmt3, err2 := db.Prepare("insert into questions_answers (candidateId, questionsId, answer, created) values ($1, $2, $3, NOW())")
		CheckErr(err2)
		rows4, err3 := db.Query("select id from questions where deleted IS NULL")
		CheckErr(err3)
		questionInfo := Questions{}
		for rows4.Next() {
			rows4.Scan(&questionInfo.Qid)
			stmt3.Query(candidateId, questionInfo.Qid, "")
		}
		//=========counting time for 7 days
		now := time.Now()
		sevenDay := time.Second * 585060 // Here 585060 is total seconds of 7 days ...
		time := now.Add(time.Duration(sevenDay))
		//=======
		stmt4, err4 := db.Prepare("INSERT INTO sessions (hash, candidateId, created, expired, challengeId, status) VALUES($1, $2, NOW(), $3, $4, $5)")
		CheckErr(err4)
		stmt4.Exec(hash, candidateId, time, challengeNo, 1)
	}else {
		flag = 1
	}
	return flag, hash
}

//if candidate registred and link is not expired yet...
func RegistreadActiveCandidate(candidateId int, email string) string {
	query, err := db.Prepare("SELECT candidateid, hash, created, expired, challengeId FROM sessions WHERE candidateid = (select id from candidates where email = ($1)) AND status = 1")
	CheckErr(err)
	row2, err2 := query.Query(email)
	CheckErr(err2)
	info := sessionInfo{}
	var hash string

	for row2.Next() {
		err := row2.Scan(&info.candidateid, &info.hash, &info.entryDate, &info.expireDate, &info.challenge)
		CheckErr(err)
	}
	//check for hash expired or not...
	remainTime :=info.expireDate.Sub(time.Now())

	if remainTime.Seconds() < 0 {

		hash = HashGenerator()

		now := time.Now()
		sevenDay := time.Second * 585060 // Here 585060 is total seconds of 7 days ...
		time := now.Add(time.Duration(sevenDay))

		stmt3,err := db.Prepare("UPDATE sessions SET status = 0 WHERE candidateId= ($1)")
		stmt3.Query(candidateId)
		CheckErr(err)

		query1, err4 := db.Prepare("INSERT INTO sessions (hash, candidateId, created, expired, challengeId, status) VALUES($1, $2, NOW(), $3, $4, $5)")
		CheckErr(err4)
		query1.Exec(hash, candidateId, time, info.challenge, 1)
	} else {
		hash = ""
		hash += info.hash
	}
	return hash
}

//for first time candidate registration ...
func NewCandidateRegistration(email string) string {
	hash := HashGenerator()
	//===== will return random challenge among no of challenges =========
	challengeNo := RandomChallengeGenerator()
	//======================================================================
	stmt1, err := db.Prepare("INSERT INTO candidates (email,created) VALUES($1,NOW())")
	stmt1.Exec(email)
	CheckErr(err)

	var candidateid int
	err2 := db.QueryRow("SELECT id FROM candidates WHERE email = ($1)",email).Scan(&candidateid)
	CheckErr(err2)

	stmt3, err4 := db.Prepare("insert into questions_answers (candidateId, questionsId, answer, created) values ($1, $2, $3, NOW())")
	CheckErr(err4)
	rows4, err5 := db.Query("select id from questions where deleted IS NULL")
	CheckErr(err5)
	questionInfo := Questions{}
	for rows4.Next() {
		rows4.Scan(&questionInfo.Qid)
		stmt3.Query(candidateid, questionInfo.Qid, "")
	}
	//=========counting time for 7 days
	now := time.Now()
	sevenDay := time.Second * 585060 // Here 585060 is total seconds of 7 days ...
	time := now.Add(time.Duration(sevenDay))
	//=======
	stmt4, err6 := db.Prepare("INSERT INTO sessions (hash, candidateId, created, expired, challengeId, status) VALUES($1, $2, NOW(), $3, $4, $5)")
	CheckErr(err6)
	stmt4.Exec(hash, candidateid, time, challengeNo, 1)
	return hash
}
