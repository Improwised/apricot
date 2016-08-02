package main

import (
	"time"
	"bytes"
	"strings"
	"strconv"
	"net/http"
	"database/sql"
	"encoding/json"
	"html/template"

	_"github.com/lib/pq"
	"github.com/zenazn/goji"
	"github.com/zenazn/goji/web"
	"github.com/improwised/apricot/services"
)

var db *sql.DB = services.SetupDB()

func indexHandler(c web.C, w http.ResponseWriter, r *http.Request) {
	email := r.FormValue("email")
	var flag int
	//Check whether email empty or not
	if email == "" {
		t, err := template.ParseFiles("./views/client/index.html")
		services.CheckErr(err)
		t.Execute(w, t)
	} else {
		var candidateId int
		err := db.QueryRow("SELECT id FROM candidates WHERE email = ($1)",email).Scan(&candidateId)
		if(err == nil){
			flag = 2
		}
		if(flag == 2) {//check If candidate already registread before and link has been expired....
			var hash string
			flag, hash = services.RegistreadExpiredCandidate(candidateId, email)
			if(flag != 1) {//flag will 1 only if candidate is not expired and link is active currently ...
				services.Mail(hash, email)
			}
		}
		if (flag == 1) {//If candidate already registread before and still active....
			hash := services.RegistreadActiveCandidate(candidateId, email)
			services.Mail(hash, email)
		}
		if (flag == 0) {//For new registration
			hash := services.NewCandidateRegistration(email)
			services.Mail(hash, email)
		}
		http.Redirect(w, r, "/confirmation", 302)
	}
}

type AllDetail struct {
	GeneralInfo services.GeneralInfo
	GetQuestions []services.GetQuestions
	// RemainTime services.TotalTime
	RemainTime int
}

func informationHandler(c web.C, w http.ResponseWriter, r *http.Request) {
	hash := r.FormValue("key")
	remainTime := services.CheckHash(c, w, r, hash)
	allDetails := AllDetail{}

	// allDetails.RemainTime = services.CalculateDayHrMin(remainTime)
	allDetails.RemainTime = int(remainTime.Seconds())


	var candidateid int
	err := db.QueryRow("select candidateId from sessions where hash = ($1)", hash).Scan(&candidateid)
	services.CheckErr(err)

	//get user detail
	allDetails.GeneralInfo = services.PersonalInfo(candidateid)
	// get all questions
	allDetails.GetQuestions, _ = services.QuestionsAnswers(candidateid)

	t, err := template.ParseFiles("./views/client/information.html")
	services.CheckErr(err)
	t.Execute(w, allDetails)
}

//information and answer of question of candidate will be save to database on keyUp ...
func dataUpdate(w http.ResponseWriter, r *http.Request) {
	data := r.FormValue("data")
	id := r.FormValue("id")
	hash := r.FormValue("hash")
	services.AutoSave(data, id, hash)
}

func challengesHandler(c web.C, w http.ResponseWriter, r *http.Request) {
	services.ProcessFormData(w, r)
	hash := r.FormValue("hash")
	http.Redirect(w, r, "/challenge?key="+ hash, 302)
}

type getChallenge struct {
	Description string
	Answer string
	Hash string
	RemainTime int
}

func challengeDetails(hash string, remainTime time.Duration) getChallenge {
	var candidateid int
	err := db.QueryRow("select candidateid from sessions where hash='"+ hash +"'").Scan(&candidateid)
	services.CheckErr(err)

	challengeInfo := getChallenge{}
	challengeInfo.Description = services.ChallengeDescription(candidateid)
	//======================================================

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
	if(err4 != nil ){
		encodeSource = ""
	}
	// db.Close()
	//===========

	//Decode the encrypted SourceCode from database...
	decoded2 := services.Decoding(encodeSource)
	//======================================================

	totalTime := int(remainTime.Seconds())

	challengeInfo.Answer = decoded2
	challengeInfo.RemainTime = totalTime
	challengeInfo.Hash = hash

	return challengeInfo
}

//get the challenge description and source code and display it ...
func challengeHandler(c web.C, w http.ResponseWriter, r *http.Request) {
	hash := r.URL.Query().Get("key")

	// ===============TODO: check hash expired or not...===================
	remainTime := services.CheckHash(c, w, r, hash)
	//=================================================

	//get the prvious attended source code and diplay it in editor ..==
	getAllChallengesInfo := challengeDetails(hash, remainTime)
	//===================================

	t, err := template.ParseFiles("./views/client/challenge.html")
	services.CheckErr(err)
	t.Execute(w, getAllChallengesInfo)
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

func languagesHandler(c web.C, w http.ResponseWriter, r *http.Request){
	bytes := services.ApiLanguages()
	w.Write([]byte(bytes))
}

type passHrResponse struct {
	Compilemessage string
	Stdout []string
}

func getHrResponse(c web.C, w http.ResponseWriter, r *http.Request){
	Source := r.FormValue("source")
	language := r.FormValue("language")
	hash := r.FormValue("hash")
	id := r.FormValue("id")
	aceLang := r.FormValue("aceLang")

	//will save the source code to database when user run the code ...
	if Source != ""{
		var clientResponse []string
		//=======
		 source, sessionid := services.SaveChallengeAnswer(hash, Source, aceLang)
		//=======

		//chek whether user run the code or submit the code
		strBody, outputDatabase, testcases, clientResponse  := services.ThirdPartyAPIResponse(id, hash, Source, language)
	 //=============================

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

			//converting both output to bytes to make comparision simple...
			outputDatabaseToBytes := []byte(outputDatabase[i])
			outputApiToBytes := []byte(outputResponse[i])

			//remove unnecessary \r ,which comes from database which is barrior in comparision ...
			outputDatabaseBytes := bytes.Replace(outputDatabaseToBytes, []byte("\r"), []byte(""), -1)

			//finally compare databse output with API output ...
			if(services.CompareString(string(outputApiToBytes), string(outputDatabaseBytes)) == 0){
				count = 1 //if both out put same..
			} else{
				count = 0 // if both outputs are not same ...
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
			query4, err := db.Prepare("insert into challenge_attempts(sessionid,input,output,solution,status,created) values($1,$2,$3,$4,$5,NOW())")
			services.CheckErr(err)
			query4.Exec(sessionid,testcases[0],outputResponse[0],source,status[0])
			// =========================================

			//will store all the test cases iƒ databse for a source code======================
			for i := 1; i < length; i++ {
				query5, err := db.Prepare("insert into challenge_attempts(sessionid,input,output,status) values($1,$2,$3,$4)")
				services.CheckErr(err)
				query5.Exec(sessionid,testcases[i],outputResponse[i],status[i])
			}
			// =========================================session will be expire....
			query6 := "update sessions set status='0' where hash='" + hash + "'"
			db.Query(query6)
		}
	}
}

func confirmationPage(c web.C, w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles("./views/client/confirmation.html")
	services.CheckErr(err)
	t.Execute(w, t)
}

func expiredPage(c web.C, w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles("./views/client/expired.html")
	services.CheckErr(err)
	t.Execute(w, t)
}

func thankYouHandler(c web.C, w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles("./views/client/thankYouPage.html")
	services.CheckErr(err)
	t.Execute(w, t)
}

func main() {

	goji.Handle("/", indexHandler)

	goji.Get("/confirmation", confirmationPage)
	goji.Post("/challenges", challengesHandler)
	goji.Get("/challenge", challengeHandler)
	goji.Get("/expired", expiredPage)
	goji.Post("/getLanguages", languagesHandler)
	goji.Post("/hrresponse", getHrResponse)
	goji.Get("/information", informationHandler)
	goji.Handle("/index", indexHandler)
	goji.Post("/saveData", dataUpdate)
	goji.Get("/thankYouPage", thankYouHandler)

	http.Handle("/assets/css/", http.StripPrefix("/assets/css/", http.FileServer(http.Dir("assets/css"))))
	http.Handle("/assets/js/", http.StripPrefix("/assets/js/", http.FileServer(http.Dir("assets/js"))))
	http.Handle("/assets/img/", http.StripPrefix("/assets/img/", http.FileServer(http.Dir("assets/img"))))
	http.Handle("/assets/fonts/", http.StripPrefix("/assets/fonts/", http.FileServer(http.Dir("assets/fonts"))))
	goji.Serve()
}
