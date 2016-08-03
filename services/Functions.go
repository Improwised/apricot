package services

import (
	"os"
	"time"
	"bytes"
	"net/http"
	"math/rand"
	"crypto/sha1"
	"encoding/hex"
	"database/sql"
	"encoding/json"
	"encoding/base64"

	_"github.com/lib/pq"
	"gopkg.in/gomail.v2"
	"github.com/zenazn/goji/web"
	"github.com/improwised/apricot/dbconfig"
)

func CheckErr(err error) {
	if err != nil {
		panic(err)
	}
}

type Configuration struct {
	DbName string
	UserName string
	EmailId string
	EmailPassword string
}

func SetupDB() *sql.DB {
	var path string
	Go_Env := (os.Getenv("GO_ENV2"))

	if (Go_Env == ""){
		Go_Env = "settings"
		path = "./dbconfig/test-files/"+ Go_Env +".json"
	} else if(Go_Env == "production"){
		path = "../dbconfig/test-files/"+ Go_Env +".json"
	} else if(Go_Env == "testing") {
		path = "../dbconfig/test-files/"+ Go_Env +".json"
	}
	connectionString := dbconfig.PostgresConnectionString(path, "disable") // second parameter for sslmode
	db1, err := sql.Open("postgres", connectionString)
	CheckErr(err)
	return db1
}

var db *sql.DB = SetupDB()

func Mail(key string, mail string) {
	file, err := os.Open("./config/configuration.json")
	CheckErr(err)
	decoder := json.NewDecoder(file)
	configuration := Configuration{}
	decoder.Decode(&configuration)

	Email_id := configuration.EmailId
	Email_Pass := configuration.EmailPassword

	m := gomail.NewMessage()
	m.SetHeader("From", Email_id)
	m.SetHeader("To", mail)
	m.SetHeader("Subject", " Interview Process - Improwised Technologies")
	var mailBody = " <div style='font-size: 15px'>"
			mailBody += "This is an automated mail from Improwised Technology for your interview process. <br>"
			mailBody +=	"Visit the link below to start with your interview process. <br><br>"
			mailBody +=	"<div>http://localhost:8000/information?key="+ key + "</div>"
			mailBody +=	"<br><div style='color: Red'>"
			mailBody +=	"Note: Your link is active for next 7 days only. </div><br> "
			mailBody +=	"Please ignore if you are done with the process. <br><br><br>"
			mailBody +=	"Best Regards, <br> Improwised Technologies </div>"

	m.SetBody("text/html", mailBody)
	d := gomail.NewPlainDialer("smtp.gmail.com", 587, Email_id, Email_Pass)
	if err := d.DialAndSend(m);
	err != nil {
		CheckErr(err)
	}
}

//will decode base 64 encoded text ...
func Decoding(encodedText string) string {
	//Decode the encrypted text from database...
	decodedChallenge, err := base64.StdEncoding.DecodeString(encodedText)
	CheckErr(err)
	//==================================================

	//convert decrypted text to string from byte and store it into structure====================
	var M = map[string]*struct{ challenge string } {
		"foo": {"Challenge"},
	}
	M["foo"].challenge = string(decodedChallenge)

	return M["foo"].challenge
	//======================================================
}

//for comparision of strings ...
func CompareString(a, b string) int {
	if a == b {
		return 0
	}
	return +1
}

//this function will save final the information of candidate from information page when user submit...
func ProcessFormData(w http.ResponseWriter, r *http.Request) int {
	r.ParseForm()
	var passed int = 0
	if r.Method == "POST" {
		name := r.FormValue("name")
		contact := r.FormValue("contact")
		degree := r.FormValue("degree")
		college := r.FormValue("college")
		yearOfCompletion := r.FormValue("yearOfCompletion")
		hash := r.FormValue("hash")
		for key, values := range r.Form["message"] {   // range over map
			stmt, _ := db.Prepare("update questions_answers set answer=($1), modified=NOW() where candidateid=(select candidateid from sessions where hash=($2)) AND questionsid=($3)")
			stmt.Query(values, hash ,key+1)
		}
		var buffer bytes.Buffer
		buffer.WriteString("UPDATE candidates SET name='" + name +"', contact='" + contact + "', degree='" + degree +"', college='" + college + "', yearOfCompletion='" + yearOfCompletion + "', modified=NOW() where id = (select candidateid from sessions where hash = '" + hash + "')")
		_,err := db.Query(buffer.String())
		if err != nil {
			panic(err)
		} else {
			passed = 1
		}
	}
	return passed
}

func CheckHash(c web.C, w http.ResponseWriter, r *http.Request, hash string) time.Duration {
	if(len(hash)!=40 || hash == ""){//chek whether hash modified by candidate..
		http.Redirect(w, r, "/index", 302)
	}
	// ===============TODO check hash expired or not...===================
	query, err := db.Prepare("SELECT expired, status FROM sessions where hash = ($1)")
	CheckErr(err)
	result, err2 := query.Query(hash)
	CheckErr(err2)

	info := sessionInfo {}
	for result.Next() {
		err := result.Scan(&info.expireDate, &info.status)
		CheckErr(err)
	}
	remainTime := info.expireDate.Sub(time.Now())
	status := info.status

	if remainTime.Seconds() < 0 || status != 1 {
		http.Redirect(w, r, "/expired", 302)
	}
	return remainTime//will return remain time to expire hash ..
}

func SaveChallengeAnswer(hash string, Source string, aceLang string) (string, int) {
	var sessionid int
	err2 := db.QueryRow("select id from sessions where hash='"+ hash +"'").Scan(&sessionid)
	CheckErr(err2)

	attempts := 0
	var maxAttempts int
	err3 := db.QueryRow("select MAX(attempts) from challenge_answers where sessionid=($1)", sessionid).Scan(&maxAttempts)

	var query3 string
	if err3 != nil{
		attempts = 0;
		query3 = "INSERT INTO challenge_answers (sessionId, answer, attempts, language, created) VALUES($1, $2, $3, $4, NOW())"
	} else  {
		attempts = maxAttempts
		query3 = "INSERT INTO challenge_answers (sessionId, answer, attempts, language, modified) VALUES($1, $2, $3, $4, NOW())"
	}
	stmt3, _ := db.Prepare(query3)
	source := base64.StdEncoding.EncodeToString([]byte(Source))
	stmt3.Exec(sessionid, source, attempts + 1, aceLang)

	return source, sessionid
}

func HashGenerator() string {
	var hash string
	random := time.Now().String()

	random += "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ!@#$^&*()_+"
	h := sha1.New()
	h.Write([]byte(random))
	hash = hex.EncodeToString(h.Sum(nil))
	return hash
}

func RandomNumber(min, max int) int {
	rand.Seed(time.Now().UnixNano())
	randomNumber :=  rand.Intn(max - min) + min
	return randomNumber
}

func RandomChallengeGenerator() int {
	var challengeNo int
	var challengeCounter int

	err := db.QueryRow("select COUNT(id) from challenges").Scan(&challengeCounter)
	CheckErr(err)

	flag := true
	for flag{
		challengeNo = RandomNumber(1, (challengeCounter + 1))
		var deleted string
		err := db.QueryRow("select deleted from challenges where id = ($1)",challengeNo).Scan(&deleted)
		if(err != nil){
			flag = false
		}
	}
	return challengeNo
}

func AutoSave(data string, id string, hash string) {
	db := SetupDB()
	var table="questions_answers";
	if id == "name" || id == "contact" || id == "degree" || id == "college" || id == "yearOfCompletion"{
		table = "candidates"
	}
	var buffer bytes.Buffer
	buffer.WriteString("UPDATE ")
	buffer.WriteString(table)

	if(table == "questions_answers") {
		buffer.WriteString(" set answer=")
		buffer.WriteString("'" + data + "'")
		buffer.WriteString(",modified = NOW() where questionsid = "+ id)
		buffer.WriteString(" AND")
		buffer.WriteString(" candidateid = (select candidateId from sessions where hash =")
		buffer.WriteString("'" + hash + "'")
		buffer.WriteString(")")
	}
	if(table == "candidates") {
		buffer.WriteString(" SET ")
		buffer.WriteString(id)
		buffer.WriteString("=")
		buffer.WriteString("'" + data + "'")
		buffer.WriteString(",modified=NOW() where id=(select candidateId from sessions where hash =")
		buffer.WriteString("'" + hash + "')")
	}
	db.Query(buffer.String())
	db.Close()
}
