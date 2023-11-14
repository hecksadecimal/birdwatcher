package cherp_api

import (
	"log"
	"net/http"

	cookiejar "github.com/juju/persistent-cookiejar"
	"golang.org/x/net/publicsuffix"
)

var client http.Client
var initialized bool
var BaseUrl = "https://cherp.chat/"
var ApiUrl = BaseUrl + "api/"

// Almost everything will need access to the same HTTP client for session persistance.
func InitializeClient() {
	if initialized {
		return
	}

	jar, err := cookiejar.New(&cookiejar.Options{PublicSuffixList: publicsuffix.List, Filename: "cookies.txt"})
	if err != nil {
		log.Fatal(err)
	}

	client = http.Client{
		Jar: jar,
	}

	initialized = true
	return
}
