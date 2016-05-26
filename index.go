package main

import (
	"fmt"
	"net/http"
	"database/sql"
	_"github.com/lib/pq"
	"github.com/zenazn/goji"
	"github.com/zenazn/goji/web"
	"html/template"
	"encoding/json"
	"time"
	"encoding/base64"
	"os"
	"math/rand"
	"bytes"
	"crypto/sha1"
	"encoding/hex"
	"gopkg.in/gomail.v2"
	_"reflect"
	"io/ioutil"
	"strings"
	"strconv"
)

type Configuration struct {
	DbName string
	UserName string
}

var db *sql.DB
func setupDB() *sql.DB {
	file, _ := os.Open("./config/configuration.json")
	decoder := json.NewDecoder(file)
	configuration := Configuration{}
	decoder.Decode(&configuration)

	userName := configuration.UserName
	dbName := configuration.DbName

	dbinfo := fmt.Sprintf("user=%s dbname=%s sslmode=disable",
		userName, dbName)
	db, err := sql.Open("postgres", dbinfo)
	checkErr(err)
	return db
}

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}

type GeneralInfo struct {
	Id string
	Name string
	Contact string
	Degree string
	College string
	YearOfCompletion string
}

var hash string
func hashGenerator() {
	random := time.Now().String()
	random += "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ!@#$^&*()_+"
	h := sha1.New()
	h.Write([]byte(random))
	hash = hex.EncodeToString(h.Sum(nil))
}

func mail(key string , mail string) {
	m := gomail.NewMessage()
	m.SetHeader("From", "akumbhani666@gmail.com")
	m.SetHeader("To", mail)
	m.SetHeader("Subject", "Interview Process")
	m.SetBody("text/html", " <div style='font-size: 15px'>This is an automated mail from Improwised Technology for your interview process. <br>Visit the link below to start with your interview process. <br><br> <div>http://localhost:8000/information?key="+ key + "</div> <br><div style='color: Red'>Note: Your link is active for next 7 days only. </div> <br> Please ignore if you are done with the process. <br><br><br> Best Regards, <br> Improwised Technologies </div>")
	d := gomail.NewPlainDialer("smtp.gmail.com", 587, "akumbhani666@gmail.com", "9712186012")
	if err := d.DialAndSend(m); err != nil {
	checkErr(err)
	}
}

type noOfChallenges struct{
	Challenge int
}

func randomNumber(min, max int) int {
	rand.Seed(time.Now().UnixNano())
	randomNumber :=  rand.Intn(max - min) + min
	return randomNumber
}

type getDeletedChallenges struct {
	Deleted *time.Time
}

func randomChallengeGenerator() int {
var challengeNo int
	rows5, err := db.Query("select COUNT(id) from challenges")
	checkErr(err)
	challengeCounter := []noOfChallenges{}
	counter1 := noOfChallenges{}
	for rows5.Next() {
		err = rows5.Scan(&counter1.Challenge)
		challengeCounter = append(challengeCounter, counter1)
		checkErr(err)
	}
	flag := true
	for flag{
		challengeNo = randomNumber(1, (challengeCounter[0].Challenge + 1))
		stmt6, _ := db.Prepare("select deleted from challenges where id = ($1)")
		rows6, _ :=  stmt6.Query(challengeNo)
		dInfo := getDeletedChallenges{}
		for rows6.Next() {
			rows6.Scan(&dInfo.Deleted)
		}
		if dInfo.Deleted != nil {
		} else {
			flag = false
		}
	}
	return challengeNo
}

type Questions struct {
	Qid string
}

func indexHandler(c web.C, w http.ResponseWriter, r *http.Request) {
	email := r.URL.Query().Get("email")

	//Check whether email empty or not
	if email == "" {
		t, _ := template.ParseFiles("./views/index.html")
		t.Execute(w, t)
	} else {
		var flag int = 0
		stmt, _ := db.Prepare("SELECT id FROM candidates WHERE email = ($1)")
		rows, _ := stmt.Query(email)
		candidateid := []GeneralInfo{}
		info := GeneralInfo{}
		for rows.Next() {
			err := rows.Scan(&info.Id)
			candidateid = append(candidateid, info)
			checkErr(err)
			flag = 3
		}

		//check If candidate already registread before and link has been expired....
		if(flag == 3) {
			stmt1, _ := db.Prepare("SELECT status FROM sessions WHERE candidateid = $1 AND id = (select MAX(id) from sessions where candidateid=$2)")
			rows1, _ := stmt1.Query(candidateid[0].Id, candidateid[0].Id)
			for rows1.Next() {
				status := []sessionInfo{}
				myinfo := sessionInfo{}
				err := rows1.Scan(&myinfo.status)
				checkErr(err)

				status = append(status, myinfo)
				if(status[0].status == 0){
					flag = 2

					hashGenerator() // will generate unique hash...

					//===== will return random challenge among no of challenges =========
					challengeNo := randomChallengeGenerator()
					//======================================================================
					stmt3, _ := db.Prepare("insert into questions_answers (candidateId, questionsId, answer, created) values ($1, $2, $3, NOW())")
					rows4, _ := db.Query("select id from questions where deleted IS NULL")
					questionInfo := Questions{}
					for rows4.Next() {
						rows4.Scan(&questionInfo.Qid)
						stmt3.Query(candidateid[0].Id, questionInfo.Qid, "")
					}
					//=========counting time for 7 days
					now := time.Now()
					sevenDay := time.Hour * 24 * 7
					time := now.Add(sevenDay)
					//=======
					stmt4, _ := db.Prepare("INSERT INTO sessions (hash, candidateId, created, expired, challengeId, status) VALUES($1, $2, NOW(), $3, $4, $5)")
					stmt4.Exec(hash, candidateid[0].Id, time, challengeNo, 1)
					mail(hash,email)
				} else {

					flag = 1
				}

			}
		}

		//If candidate already registread before and still active....
		if flag == 1 {
			query, _ := db.Prepare("SELECT candidateid, hash, created, expired, challengeId FROM sessions WHERE candidateid = (select id from candidates where email = ($1)) AND status = 1")
			row2,_ := query.Query(email)
			mysession := []sessionInfo{}
			info := sessionInfo{}

			for row2.Next() {
				err := row2.Scan(&info.candidateid, &info.hash, &info.entryDate, &info.expireDate, &info.challenge)
				mysession = append(mysession, info)
				checkErr(err)
			}
			//check for hash expired or not...
			remainTime := mysession[0].expireDate.Sub(time.Now())

			if remainTime.Seconds() < 0 {

				hashGenerator()

				now := time.Now()
				sevenDay := time.Hour * 24 * 7
				time := now.Add(sevenDay)
				stmt3 := "UPDATE sessions SET status=0 WHERE candidateId=" + candidateid[0].Id
				db.Query(stmt3)

				query1, _ := db.Prepare("INSERT INTO sessions (hash, candidateId, created, expired, challengeId, status) VALUES($1, $2, NOW(), $3, $4, $5)")
				query1.Exec(hash, candidateid[0].Id, time, mysession[0].challenge, 1)
				mail(hash,email)

			} else {
				hash = ""
				hash += mysession[0].hash
				mail(hash,email)
			}

		} else if flag == 0 {//For new registration

			hashGenerator()
			//===== will return random challenge among no of challenges =========
			challengeNo := randomChallengeGenerator()
			//======================================================================
			stmt1, err := db.Prepare("INSERT INTO candidates (email,created) VALUES($1,NOW())")
			stmt1.Exec(email)
			checkErr(err)
			stmt2, err := db.Prepare("select id from candidates where email = ($1)")
			checkErr(err)
			rows2, err := stmt2.Query(email)
			information := []GeneralInfo{}
			info := GeneralInfo{}
			for rows2.Next() {
				err := rows2.Scan(&info.Id)
				information = append(information, info)
				checkErr(err)
			}

			stmt3, _ := db.Prepare("insert into questions_answers (candidateId, questionsId, answer, created) values ($1, $2, $3, NOW())")
			rows4, _ := db.Query("select id from questions where deleted IS NULL")
			questionInfo := Questions{}
			for rows4.Next() {
				rows4.Scan(&questionInfo.Qid)
				stmt3.Query(information[0].Id, questionInfo.Qid, "")
			}
			//=========counting time for 7 days
			now := time.Now()
					sevenDay := time.Hour * 24 * 7
					time := now.Add(sevenDay)
			//=======
			stmt4, _ := db.Prepare("INSERT INTO sessions (hash, candidateId, created, expired, challengeId, status) VALUES($1, $2, NOW(), $3, $4, $5)")
			stmt4.Exec(hash, information[0].Id, time, challengeNo, 1)
			mail(hash,email)
		}
		http.Redirect(w, r, "/confirmation", 302)
	}
}

type GetQuestions struct {
	Questions string
	Id string
	Ans  string
}

type GetAnswers struct {
	Answer string
	Qid string
}

type AllDetail struct {
	GeneralInfo []GeneralInfo
	GetQuestions []GetQuestions
}

func informationHandler(c web.C, w http.ResponseWriter, r *http.Request) {
	hash := r.URL.Query().Get("key")

	if(len(hash)!=40 || hash == ""){//chek whether hash modified by candidate..
		http.Redirect(w, r, "/index", 302)
	}

	// ===============TODO check hash expired or not...===================
	query, _ := db.Prepare("SELECT expired,status FROM sessions where hash = ($1)")
	result, _ := query.Query(hash)

	mysession := []sessionInfo{}
	info := sessionInfo{}
	for result.Next() {
		err := result.Scan(&info.expireDate, &info.status)
		mysession = append(mysession, info)
		checkErr(err)
	}
	remainTime := mysession[0].expireDate.Sub(time.Now())
	status := mysession[0].status
	if remainTime.Seconds() < 0 || status != 1 {
	 http.Redirect(w, r, "/expired", 302)
		return
	}
	allDetails := AllDetail{}
	// get all questions
	rows, _ := db.Query("select id, description from questions where deleted IS NULL order by sequence")
	questionsInfo := []GetQuestions{}
	qinfo := GetQuestions{}
	for rows.Next() {
		err := rows.Scan(&qinfo.Id, &qinfo.Questions)
		questionsInfo = append(questionsInfo, qinfo)
		checkErr(err)
	}
	stmt2, _ := db.Prepare("SELECT answer, questionsid FROM questions_answers where candidateId = (select candidateId from sessions where hash = ($1))")
	rows2, _ := stmt2.Query(hash)

	for rows2.Next() {
		getanswer := GetAnswers{}
		err := rows2.Scan(&getanswer.Answer, &getanswer.Qid)
		for index, element := range questionsInfo {
			if(element.Id == getanswer.Qid) {
				questionsInfo[index].Ans = getanswer.Answer
			}
		}
		checkErr(err)
	}

	//get user detail
	User := []GeneralInfo{}
	user := GeneralInfo{}
	stmt3, _ := db.Prepare("select id, name, contact, degree, college,yearOfCompletion from candidates where id = (select candidateId from sessions where hash = ($1))")
	row3, _ := stmt3.Query(hash)
	for row3.Next() {
		row3.Scan(&user.Id, &user.Name, &user.Contact, &user.Degree, &user.College, &user.YearOfCompletion)
		User = append(User, user)
	}
	dataUpdate(w, r, hash)
	t, _ := template.ParseFiles("./views/information.html")

	allDetails.GeneralInfo = User
	allDetails.GetQuestions = questionsInfo

	t.Execute(w, allDetails)
}

func dataUpdate(w http.ResponseWriter, r *http.Request, hash string) {
	db = setupDB()
	data := r.URL.Query().Get("data")
	id := r.URL.Query().Get("id")
	var table="questions_answers";
	if id == "name" || id == "contact" || id == "degree" || id == "college" || id == "yearOfCompletion"{
		table = "candidates"
	}
	var buffer bytes.Buffer
	buffer.WriteString("UPDATE ")
	buffer.WriteString(table)

	if(table == "questions_answers"){
		buffer.WriteString(" set answer=")
		buffer.WriteString("'" + data + "'")
		buffer.WriteString(",modified=NOW() where questionsid="+ id)
		buffer.WriteString(" AND")
		buffer.WriteString(" candidateid=(select candidateId from sessions where hash=")
		buffer.WriteString("'" + hash + "'")
		buffer.WriteString(")")
	}

	if(table == "candidates"){
		buffer.WriteString(" SET ")
		buffer.WriteString(id)
		buffer.WriteString("=")
		buffer.WriteString("'" + data + "'")
		buffer.WriteString(",modified=NOW() where id=(select candidateId from sessions where hash =")
		buffer.WriteString("'" + hash + "')")
	}
	db.Query(buffer.String())
}

func challengesHandler(c web.C, w http.ResponseWriter, r *http.Request) {
	processFormData(w, r)
	hash := r.FormValue("hash")

	http.Redirect(w, r, "/challenge?key="+ hash, 302)
}

type getChallenge struct {
	Description string
	Answer string
}

type getAllChallengesInfo struct {
	GetChallenge []getChallenge
	GetSource []getChallenge
	Hash string
}

type getLanguages struct {
	Languages []string
}

type getHash struct {
	Hash string
}

//save programme and no of attempts into database while compiling
func challengeHandler(c web.C, w http.ResponseWriter, r *http.Request) {

	// ===============TODO check hash expired or not...===================
	hash := r.URL.Query().Get("key")

	if(len(hash)!=40 || hash == ""){//chek whether hash modified by candidate..
		http.Redirect(w, r, "/index", 302)
	}

	//will check whether hash expired or not...===========
	query, _ := db.Prepare("SELECT expired,status FROM sessions where hash = ($1)")
	result, _ := query.Query(hash)

	mysession := []sessionInfo{}
	info := sessionInfo{}
	for result.Next() {
		err := result.Scan(&info.expireDate,&info.status)
		mysession = append(mysession, info)

		checkErr(err)
	}
	remainTime := mysession[0].expireDate.Sub(time.Now())
	status := mysession[0].status
	if remainTime.Seconds() < 0 || status != 1 {
		http.Redirect(w, r, "/expired", 302)
		return
	}
	//==============================================================

	source := r.FormValue("source")
	key := r.FormValue("hash")
	var buffer, buffer2 bytes.Buffer

	if source != "" {
		buffer.WriteString("select id from sessions where hash=")
		buffer.WriteString("'" + key + "'")
		query := buffer.String()

		rows5, _ := db.Query(query)
		mysession := []sessionInfo{}
		info := sessionInfo{}
		for rows5.Next() {
			err := rows5.Scan(&info.sessionid)
			mysession = append(mysession, info)
			checkErr(err)
		}
		sessionid := mysession[0].sessionid
		buffer2.WriteString("select MAX(attempts) from challenge_answers where sessionid=")
		buffer2.WriteString("'" + sessionid + "'")
		query2 := buffer2.String()
		rows6, err := db.Query(query2)

		attempts := 0
		mysession2 := []sessionInfo{}
		info2 := sessionInfo{}
		var query3 string
		for rows6.Next() {
			err := rows6.Scan(&info2.attempts)
			mysession2 = append(mysession2, info2)

			if err != nil{
				attempts = 0;
				query3 = "INSERT INTO challenge_answers (sessionId, answer, attempts, created) VALUES($1, $2, $3, NOW())"
			} else  {
				attempts = mysession2[0].attempts
				query3 = "INSERT INTO challenge_answers (sessionId, answer, attempts, modified) VALUES($1, $2, $3, NOW())"
			}
		}
		stmt1, err := db.Prepare(query3)
		stmt1.Exec(sessionid, source, attempts + 1)
		checkErr(err)
	}
	getAllChallengesInfo := getAllChallengesInfo{}
	challengeInfo := []getChallenge{}
	cInfo := getChallenge{}

	stmt1, _ := db.Prepare("select description from challenges where id = (select challengeId from sessions where hash = ($1))")
	rows1,_ := stmt1.Query(hash)

	var encodedChallenge string
	for rows1.Next() {
		err := rows1.Scan(&encodedChallenge)
		checkErr(err)
	}
	//Decode the encrypted challenge from database...
	decodedChallenge, err := base64.StdEncoding.DecodeString(encodedChallenge)
	if err != nil {
		fmt.Println("decode error:", err)
		return
	}
	//==================================================

	//convert decrypted challenge to string from byte and store it into structure====================
	var m = map[string]*struct{ challenge string }{
	"foo": {"Challenge"},
	}

	m["foo"].challenge = string(decodedChallenge)

	cInfo.Description = m["foo"].challenge
	//======================================================

	challengeInfo = append(challengeInfo, cInfo)

	//getting last answer for perticullar challege
	var query1 = "select answer from challenge_answers"
			query1 += " where sessionid="
			query1 += " (select id from sessions where hash ='"+ hash +"')"
			query1 += " AND id=(select MAX(id) from challenge_answers"
			query1 += " where sessionid="
			query1 += " (select id from sessions"
			query1 += " where hash ='"+ hash +"'))"

	result1, _ := db.Query(query1)
	sourcecode := []getChallenge{}
	infoChallenge := getChallenge{}
	for result1.Next() {
		err := result1.Scan(&infoChallenge.Answer)
		sourcecode = append(sourcecode, infoChallenge)
		checkErr(err)
	}
	//===========

	getAllChallengesInfo.GetChallenge = challengeInfo
	getAllChallengesInfo.GetSource = sourcecode
	t, _ := template.ParseFiles("./views/challenge.html")
	t.Execute(w, getAllChallengesInfo)
}

func processFormData(w http.ResponseWriter, r *http.Request)  {
	db = setupDB()
	r.ParseForm()

	if r.Method == "POST" {
		name := r.FormValue("name")
		contact := r.FormValue("contact")
		degree := r.FormValue("degree")
		college := r.FormValue("college")
		yearOfCompletion := r.FormValue("yearOfCompletion")
		hash := r.FormValue("hash")
		for key, values := range r.Form["message"] {   // range over map
			stmt, _ := db.Prepare("update questions_answers set answer=($1),modified=NOW() where candidateid=(select candidateid from sessions where hash=($2)) AND questionsid=($3)")
			stmt.Query(values, hash ,key+1)
		}
		var buffer bytes.Buffer
		buffer.WriteString("UPDATE candidates SET name='" + name +"',contact='" + contact + "',degree='" + degree +"',college='" + college + "',yearOfCompletion='" + yearOfCompletion + "',modified=NOW() where id=(select candidateid from sessions where hash ='" + hash + "')")
		db.Query(buffer.String())
	}
}

type GetChallenge_cases struct {
	Input string
	Output string
}

type HrCodeCheckerResponse struct {
	Result struct {
		CallbackURL string
		CensoredCompileMessage string
		CensoredStderr []string
		CodecheckerHash string
		Compilemessage string
		CreatedAt string
		DiffStatus []int
		ErrorCode int
		Hash string
		Memory []int
		Message []string
		ResponseS3Path string
		Result int
		Server string
		Signal []int
		Stderr []bool
		Stdout []string
		Time []int
	}
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

type AllLanguages struct {
	Languages struct {
		Names struct {
			C string
			Cpp string
			Java string
			Csharp string
			Php string
			Ruby string
			Python string
			Perl string
			Haskell string
			Clojure string
			Scala string
			Bash string
			Lua string
			Erlang string
			Javascript string
			Go string
			D string
			Ocaml string
			Pascal string
			Sbcl string
			Python3 string
			Groovy string
			Objectivec string
			Fsharp string
			Cobol string
			Visualbasic string
			Lolcode string
			Smalltalk string
			Tcl string
			Whitespace string
			Tsql string
			Java8 string
			Db2 string
			Octave string
			R string
			Xquery string
			Racket string
			Rust string
			Fortran string
			Swift string
			Oracle string
			Mysql string
		}
		Codes struct {
			C int
			Cpp int
			Java int
			Python int
			Perl int
			Php int
			Ruby int
			Mysql int
			Oracle int
			Haskell int
			Clojure int
			Bash int
			Scala int
			Erlang int
			Lua int
			Javascript int
			Go int
			D int
			Ocaml int
			R int
			Pascal int
			Sbcl int
			Python3 int
			Groovy int
			Objectivec int
			Fsharp int
			Cobol int
			Visualbasic int
			Lolcode int
			Smalltalk int
			Tcl int
			Whitespace int
			Tsql int
			Java8 int
			Db2 int
			Octave int
			Xquery int
			Racket int
			Rust int
			Swift int
			Fortran int
		}
	}
}

func languagesHandler(c web.C, w http.ResponseWriter, r *http.Request){
		var buffer2 bytes.Buffer

		buffer2.WriteString("format=json&wait=true")
		req, _ := http.NewRequest("GET", "http://api.hackerrank.com/checker/languages.json", strings.NewReader(buffer2.String()))
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded; param=value")
		client := http.Client{}
		resp, _ := client.Do(req)
		var strBody string
		body, _ := ioutil.ReadAll(resp.Body)
		strBody = string(body)
		allLang := AllLanguages{}
		json.Unmarshal([]byte(strBody), &allLang)

		bytes, _ := json.Marshal(allLang)
		w.Write([]byte(bytes))
}

type passHrResponse struct {
	Compilemessage string
	Stdout []string
}

func CompareString(a, b string) int {
	if a == b {
		return 0
	}
	return +1
}

func getHrResponse(c web.C, w http.ResponseWriter, r *http.Request){

	source := r.FormValue("source")
	language := r.FormValue("language")
	hash := r.FormValue("hash")
	id := r.FormValue("id")
	var sessionid string
	//will save the source code to databse when user run the code
	var buffer2 bytes.Buffer
	if source != ""{
		var buffer bytes.Buffer
		buffer.WriteString("select id from sessions where hash=")
		buffer.WriteString("'" + hash + "'")
		query := buffer.String()
		rows5, _ := db.Query(query)
		mysession := []sessionInfo{}
		info := sessionInfo{}
		for rows5.Next() {
			err := rows5.Scan(&info.sessionid)
			mysession = append(mysession, info)
			checkErr(err)
		}
		sessionid = mysession[0].sessionid

		buffer2.WriteString("select MAX(attempts) from challenge_answers where sessionid=")
		buffer2.WriteString("'" + sessionid + "'")
		query2 := buffer2.String()

		rows6, err := db.Query(query2)
		attempts := 0
		mysession2 := []sessionInfo{}
		info2 := sessionInfo{}
		var query3 string
		for rows6.Next() {
			err := rows6.Scan(&info2.attempts)
			mysession2 = append(mysession2, info2)

			if err != nil{
				attempts = 0;
				query3 = "INSERT INTO challenge_answers (sessionId, answer, attempts, created) VALUES($1, $2, $3, NOW())"
			} else  {
				attempts = mysession2[0].attempts
				query3 = "INSERT INTO challenge_answers (sessionId, answer, attempts, modified) VALUES($1, $2, $3, NOW())"
			}
		}

		stmt1, err := db.Prepare(query3)
		stmt1.Exec(sessionid, source, attempts + 1)
		checkErr(err)

		//chek whether user run the code or submit the code
		var testcases []string
		var outputDatabase []string
		var clientResponse []string
		//when user run the code only default testcase will chek for compilation
		if id == "runCode" {
			stmt1, _ := db.Prepare("select input, output from challenge_cases where challengeid = (select challengeid from sessions where hash = ($1) ) and defaultcase =true")
			rows1, _ := stmt1.Query(hash)
			cCases := []GetChallenge_cases{}
			c := GetChallenge_cases{}
			for rows1.Next() {
				rows1.Scan(&c.Input, &c.Output)
				cCases = append(cCases, c)
				testcases = append(testcases, c.Input)
				outputDatabase = append(outputDatabase, c.Output)
			}
			inputDatabase := cCases[0].Input//input from database
			clientResponse = append(clientResponse, inputDatabase)
			outputDatabases := cCases[0].Output//output from database
			clientResponse = append(clientResponse, outputDatabases)
		}

		//when user submit the code all the testcases will chek for compilation and response from hackerrank will be stored in database
		if id == "submitCode" {
			cCases := []GetChallenge_cases{}
			c := GetChallenge_cases{}
			stmt1, _ := db.Prepare("select input, output from challenge_cases where challengeid = (select challengeid from sessions where hash = ($1))")
			rows1, _ := stmt1.Query(hash)
			for rows1.Next() {
				rows1.Scan(&c.Input, &c.Output)
				cCases = append(cCases, c)
				testcases = append(testcases, c.Input)
				outputDatabase = append(outputDatabase, c.Output)
			}

		}
		bytetestcases, err := json.Marshal(testcases)
		if err != nil {
			fmt.Printf("Error: %s", err)
		}

		var strTestCases string
		strTestCases =  string(bytetestcases)
		api_key := "hackerrank|768030-708|2f417cf30f50ac1385dd76338a5e5c78c7dd87e9"
		var buffer2 bytes.Buffer
		buffer2.WriteString("format=json&wait=true")
		buffer2.WriteString("&source=" + source)
		buffer2.WriteString("&lang=" + language)
		buffer2.WriteString("&api_key="+ api_key)
		buffer2.WriteString("&testcases=" + strTestCases )
		req, _ := http.NewRequest("POST", "http://api.hackerrank.com/checker/submission.json", strings.NewReader(buffer2.String()))
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded; param=value")
		client := http.Client{}
		resp, _ := client.Do(req)
		var strBody string
		body, _ := ioutil.ReadAll(resp.Body)
		strBody = string(body)
		HrInfo := HrCodeCheckerResponse{}
		json.Unmarshal([]byte(strBody), &HrInfo)
		HrMesageToDisplay := passHrResponse{}
		HrMesageToDisplay.Compilemessage = HrInfo.Result.Compilemessage
		HrMesageToDisplay.Stdout = HrInfo.Result.Stdout

		if(HrMesageToDisplay.Stdout == nil){
			clientResponse = append(clientResponse, " ")
		} else {
			clientResponse = append(clientResponse, HrInfo.Result.Stdout[0])
		}
		clientResponse = append(clientResponse, HrMesageToDisplay.Compilemessage)
		 //check for all the testcases from database..
		var outputResponse []string
		var length = len(outputDatabase)
		var status []int

		if(HrMesageToDisplay.Stdout == nil){
			for i := 0; i < length; i++ {
				HrMesageToDisplay.Stdout = append(HrMesageToDisplay.Stdout ,"0")
			}
		}

		outputResponse = HrMesageToDisplay.Stdout
		var count int

		for i := 0; i < length; i++ {
			outputDatabase[i] = strings.TrimSpace(outputDatabase[i]);
			outputResponse[i] = strings.TrimSpace(outputResponse[i]);
			//
			if(CompareString(outputDatabase[i], outputResponse[i]) == 0){
				count = 1
			} else{
				count = 0
			}
			status = append(status, count)
		}
		clientResponse = append(clientResponse,strconv.Itoa(status[0]))
		//converting to JSON=======
		bytes, err := json.Marshal(clientResponse)
		if err != nil {
			return
		}
		// =====================
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(bytes))

		// will store final response from api to database..
		if id == "submitCode" {
			// will store source code in database once========================
			query4, _ := db.Prepare("insert into challenge_attempts(sessionid,input,output,solution,status,created) values($1,$2,$3,$4,$5,NOW())")
			query4.Exec(sessionid,testcases[0],outputResponse[0],source,status[0])
			// =========================================

			//will store all the test cases iƒ databse for a source code======================
			for i:=1;i<length;i++{
				query5, _ := db.Prepare("insert into challenge_attempts(sessionid,input,output,status) values($1,$2,$3,$4)")
				query5.Exec(sessionid,testcases[i],outputResponse[i],status[i])
			}
			// =========================================session will be expire....
			query6 := "update sessions set status='0' where hash='" + hash + "'"
			db.Query(query6)
		}
	}
}

func confirmationPage(c web.C, w http.ResponseWriter, r *http.Request) {
	t, _ := template.ParseFiles("./views/confirmation.html")
	t.Execute(w, t)
}

func expiredPage(c web.C, w http.ResponseWriter, r *http.Request) {
	t, _ := template.ParseFiles("./views/expired.html")
	t.Execute(w, t)
}

func thankYouHandler(c web.C, w http.ResponseWriter, r *http.Request) {
	t, _ := template.ParseFiles("./views/thankYouPage.html")
	t.Execute(w, t)
}

func main() {
	db = setupDB()
	defer db.Close()

	goji.Get("/index", indexHandler)
	goji.Get("/", indexHandler)
	goji.Handle("/information", informationHandler)
	goji.Handle("/challenges", challengesHandler)
	goji.Handle("/challenge", challengeHandler)
	goji.Post("/hrresponse", getHrResponse)
	goji.Post("/getLanguages", languagesHandler)
	goji.Handle("/confirmation", confirmationPage)
	goji.Handle("/thankYouPage", thankYouHandler)
	goji.Handle("/expired", expiredPage)

	http.Handle("/assets/css/", http.StripPrefix("/assets/css/", http.FileServer(http.Dir("assets/css"))))
	http.Handle("/assets/js/", http.StripPrefix("/assets/js/", http.FileServer(http.Dir("assets/js"))))
	http.Handle("/assets/img/", http.StripPrefix("/assets/img/", http.FileServer(http.Dir("assets/img"))))
	http.Handle("/assets/fonts/", http.StripPrefix("/assets/fonts/", http.FileServer(http.Dir("assets/fonts"))))
	http.Handle("/vendor/github.com/ajaxorg/ace/lib/ace/mode/", http.StripPrefix("/vendor/github.com/ajaxorg/ace/lib/ace/mode/", http.FileServer(http.Dir("vendor/github.com/ajaxorg/ace/lib/ace/mode/"))))
	goji.Serve()
}
