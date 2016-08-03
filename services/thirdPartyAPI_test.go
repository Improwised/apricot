package services

import (
	"os"
	"testing"
	"io/ioutil"

	"github.com/improwised/apricot/services"
)

func TestMain(m *testing.M) {
	// Run Tests
	exitCode := m.Run()
	// Exit
	os.Exit(exitCode)
}

func TestthirdPartyAPI(t *testing.T) {
	actual := services.ThirdPartyAPI("2", "print 3", "python")

	// response from api ===============
	var strBody string
	body, _ := ioutil.ReadAll(actual)
	strBody = string(body)
	// =================================
	if strBody == ""{
		t.Errorf("Test(TestthirdPartyAPI) failed")
	}
}

func TestApiLanguages(t *testing.T){
	actual := services.ApiLanguages()
	if (actual == ""){
		t.Errorf("Test(TestApiLanguages) failed")
	}
}