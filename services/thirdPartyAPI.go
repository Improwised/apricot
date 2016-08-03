package services

import (
	"io"
	"fmt"
	"bytes"
	"net/url"
	"strings"
	"net/http"
	"io/ioutil"
	"encoding/json"
)

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

func ThirdPartyAPI(strTestCases string, Source string, language string) io.ReadCloser {
	api_key := "hackerrank|768030-708|2f417cf30f50ac1385dd76338a5e5c78c7dd87e9"
	baseUrl := ""

	params := url.Values{}
	params.Add("wait", "true")
	params.Add("format", "json")
	params.Add("lang", language)
	params.Add("api_key", api_key)
	params.Add("testcases", strTestCases)
	params.Add("source", Source)

	finalUrl := baseUrl + params.Encode()

	req, err := http.NewRequest("POST", "http://api.hackerrank.com/checker/submission.json", strings.NewReader(finalUrl))
	CheckErr(err)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded; param=value")
	client := http.Client{}//Request to API ...

	resp, err2 := client.Do(req)//response from API ...
	CheckErr(err2)
	return resp.Body
}

func ApiLanguages() string{
	var buffer2 bytes.Buffer

	buffer2.WriteString("format=json&wait=true")
	req, err := http.NewRequest("GET", "http://api.hackerrank.com/checker/languages.json", strings.NewReader(buffer2.String()))
	CheckErr(err)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded; param=value")
	client := http.Client{}
	resp, err2 := client.Do(req)
	CheckErr(err2)
	var strBody string
	body, err3 := ioutil.ReadAll(resp.Body)
	CheckErr(err3)
	strBody = string(body)
	allLang := AllLanguages{}
	json.Unmarshal([]byte(strBody), &allLang)
	bytes, err4 := json.Marshal(allLang)
	CheckErr(err4)

	return string(bytes)
}

type GetChallenge_cases struct {
	Input string
	Output string
}

func ThirdPartyAPIResponse(id string, hash string, Source string, language string) (string, []string, []string, []string){
	//when user run the code only default testcase will chek for compilation
	var testcases []string
	var outputDatabase []string
	var clientResponse []string

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
		stmt1, err := db.Prepare("select input, output from challenge_cases where challengeid = (select challengeid from sessions where hash = ($1))")
		CheckErr(err)
		rows1, err2 := stmt1.Query(hash)
		CheckErr(err2)
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

	//testcase for challenge passes to compile ..========
	var strTestCases string
	strTestCases =  string(bytetestcases)
	//============================

	// api will be call ..======
	hrResponse := ThirdPartyAPI(strTestCases, Source, language)
	//===================

	// response from api ===============
	var strBody string
	body, _ := ioutil.ReadAll(hrResponse)
	strBody = string(body)
	//=================================

	return strBody, outputDatabase, testcases, clientResponse
}