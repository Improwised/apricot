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
  "os"
  "math/rand"
  "bytes"
  "crypto/sha1"
  "encoding/hex"
  _"github.com/icza/session"
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
  m.SetHeader("Subject", "Hello!")
  m.SetBody("text/html", " <div style='font-size: 15px'>This is an automated mail from Improwised Technologies for your interview process. <br>Visit the link below to start with your interview process. <br><br> <div style='color: Blue'>localhost:8000/information?key="+ key + "</div> <br><div style='color: Red'>Note: Your link is active for next 7 days only. </div> <br><br> Please ignore if you are done with the process. <br><br><br> Best Regards, <br> Improwised Technologies </div>")
  d := gomail.NewPlainDialer("smtp.gmail.com", 587, "akumbhani666@gmail.com", "9712186012")
  if err := d.DialAndSend(m); err != nil {
  checkErr(err)
  }
}

type noOfChallenges struct{
  Challenges int
}

func random(min, max int) int {
  rand.Seed(time.Now().Unix())
  return rand.Intn(max - min) + min
}

type getDeletedChallenges struct {
  Deleted *time.Time
}

func randomChallengeGenerator() string {
  rows5, err := db.Query("select COUNT(id) from challenges")
  checkErr(err)
  challengeCounter := []noOfChallenges{}
  counter1 := noOfChallenges{}
  for rows5.Next() {
  err = rows5.Scan(&counter1.Challenges)
    challengeCounter = append(challengeCounter, counter1)
    checkErr(err)
  }
  myrand := random(1, (challengeCounter[0].Challenges + 1))
  challengeNo := strconv.Itoa(myrand)
  stmt6, _ := db.Prepare("select deleted from challenges where id = ($1)")
  rows6, err := stmt6.Query(challengeNo)
  checkErr(err)
  dInfo := getDeletedChallenges{}
  for rows6.Next() {
    rows6.Scan(&dInfo.Deleted)
  }
  if dInfo.Deleted == nil {

  }else {
    randomChallengeGenerator()
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
      flag = 1
    }
    //If candidate already registread before
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
      stmt1, err := db.Prepare("INSERT INTO candidates (email) VALUES($1)")
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
  // ===============TODO check hash expired or not...===================
  query, _ := db.Prepare("SELECT expired FROM sessions where hash = ($1)")
  result, _ := query.Query(hash)

  mysession := []sessionInfo{}
  info := sessionInfo{}
  for result.Next() {
    err := result.Scan(&info.expireDate)
    mysession = append(mysession, info)
    checkErr(err)
  }

  remainTime := mysession[0].expireDate.Sub(time.Now())
  if remainTime.Seconds() < 0 {
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
    //  buffer.WriteString("'" + id + "'")
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
  hash := r.FormValue("email")
  http.Redirect(w, r, "/challenge?key="+ hash, 302)
}

type getChallenge struct {
  Description string
}

type getAllChallengesInfo struct {
  GetChallenge []getChallenge
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

  query, _ := db.Prepare("SELECT expired FROM sessions where hash = ($1)")
  result, _ := query.Query(hash)

  mysession := []sessionInfo{}
  info := sessionInfo{}
  for result.Next() {
    err := result.Scan(&info.expireDate)
    mysession = append(mysession, info)
    checkErr(err)
  }

  remainTime := mysession[0].expireDate.Sub(time.Now())
  if remainTime.Seconds() < 0 {
    http.Redirect(w, r, "/expired", 302)
    return
  }

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

  for rows1.Next() {
    err := rows1.Scan(&cInfo.Description)
    challengeInfo = append(challengeInfo, cInfo)
    checkErr(err)
  }

  getAllChallengesInfo.GetChallenge = challengeInfo
  t, _ := template.ParseFiles("./views/challenge.html")
  t.Execute(w, getAllChallengesInfo)
}

func apihandler(c web.C, w http.ResponseWriter, r *http.Request){
  if origin := r.Header.Get("Host"); origin != "" {
    w.Header().Set("Access-Control-Allow-Origin", origin)
  }
  w.Header().Set("Access-Control-Allow-Origin" , "*")
  w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
  t, _ := template.ParseFiles("./views/hrapi.html")
  t.Execute(w, t)
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
    hash := r.FormValue("email")
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
  attempts int
  entryDate time.Time
  modifyDate time.Time
  expireDate time.Time
  challenge int
}

type passHrResponse struct {
  Compilemessage string
  Stdout []string
}

func getHrResponse(c web.C, w http.ResponseWriter, r *http.Request){
  source := r.FormValue("source")
  language := r.FormValue("language")
  hash := r.FormValue("hash")
  id := r.FormValue("id")
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

  var testcases []string
  if id == "runCode" {
    stmt1, _ := db.Prepare("select input, output from challenge_cases where challengeid = (select challengeid from sessions where hash = ($1) ) and defaultcase = true")
    rows1, _ := stmt1.Query(hash)
    cCases := []GetChallenge_cases{}
    c := GetChallenge_cases{}
    for rows1.Next() {
      rows1.Scan(&c.Input, &c.Output)
      cCases = append(cCases, c)
      testcases = append(testcases, c.Input)
    }
  }

  if id == "submitCode" {
    cCases := []GetChallenge_cases{}
    c := GetChallenge_cases{}
    stmt1, _ := db.Prepare("select input, output from challenge_cases where challengeid = (select challengeid from sessions where hash = ($1))")
    rows1, _ := stmt1.Query(hash)
    for rows1.Next() {
      rows1.Scan(&c.Input, &c.Output)
      cCases = append(cCases, c)
      testcases = append(testcases, c.Input)
    }
  }

  bytetestcases, err := json.Marshal(testcases)
  if err != nil {
    fmt.Printf("Error: %s", err)
  }

  var strTestCases string
  strTestCases =  string(bytetestcases)
  api_key := "hackerrank|768030-708|2f417cf30f50ac1385dd76338a5e5c78c7dd87e9"
  var buffer bytes.Buffer
  buffer.WriteString("format=json&wait=true")
  buffer.WriteString("&source=" + source)
  buffer.WriteString("&lang=" + language)
  buffer.WriteString("&api_key="+ api_key)
  buffer.WriteString("&testcases=" + strTestCases )
  req, _ := http.NewRequest("POST", "http://api.hackerrank.com/checker/submission.json", strings.NewReader(buffer.String()))
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
  t, _ := template.ParseFiles("./views/challenge.html")
  t.Execute(w, t)
}

func confirmationPage(c web.C, w http.ResponseWriter, r *http.Request) {
  t, _ := template.ParseFiles("./views/confirmation.html")
  t.Execute(w, t)
}

func expiredPage(c web.C, w http.ResponseWriter, r *http.Request) {
  t, _ := template.ParseFiles("./views/expired.html")
  t.Execute(w, t)
}

func main() {
  db = setupDB()
  defer db.Close()
  goji.Handle("/index", indexHandler)
  goji.Handle("/hrapi", apihandler)
  goji.Handle("/information", informationHandler)
  goji.Handle("/challenges", challengesHandler)
  goji.Handle("/challenge", challengeHandler)
  goji.Handle("/hrresponse", getHrResponse)
  goji.Handle("/confirmation", confirmationPage)
  goji.Handle("/expired", expiredPage)
  http.Handle("/assets/css/", http.StripPrefix("/assets/css/", http.FileServer(http.Dir("assets/css"))))
  http.Handle("/assets/js/", http.StripPrefix("/assets/js/", http.FileServer(http.Dir("assets/js"))))
  http.Handle("/assets/img/", http.StripPrefix("/assets/img/", http.FileServer(http.Dir("assets/img"))))
  http.Handle("/assets/fonts/", http.StripPrefix("/assets/fonts/", http.FileServer(http.Dir("assets/fonts"))))
  goji.Serve()
}
