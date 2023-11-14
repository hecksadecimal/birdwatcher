package cherp_api

import (
	"bufio"
	"encoding/json"
	"log"
	"net/url"

	cookiejar "github.com/juju/persistent-cookiejar"
)

var loginUrl = ApiUrl + "user/login"

var LoggedIn = false
var Username = ""

func Logout() {
	defer SaveJar()
	client.Jar.(*cookiejar.Jar).RemoveAll()
	LoggedIn = false
	Username = ""
}

func Login(user string, pass string) (success bool) {
	defer SaveJar()
	Logout()

	csrf := GetCsrf("")
	resp, err := client.PostForm(loginUrl, url.Values{"username": {user}, "password": {pass}, "csrfname": {csrf["csrfname"].(string)}, "csrf": {csrf["csrf"].(string)}})
	if err != nil {
		log.Fatal(err)
		return false
	}
	defer resp.Body.Close()

	var jsonData map[string]interface{}
	scanner := bufio.NewScanner(resp.Body)
	for i := 0; scanner.Scan() && i < 5; i++ {
		data := scanner.Text()
		json.Unmarshal([]byte(data), &jsonData)
	}
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
		return false
	}
	if jsonData["status"].(string) != "success" {
		return false
	}

	for _, cookie := range resp.Cookies() {
		if cookie.Name == "cherpusername" {
			Username = cookie.Value
		}
	}

	LoggedIn = true
	return true
}

func SessionValid() bool {
	if LoggedIn {
		return true
	}

	defer SaveJar()

	resp, err := client.Get(unreadChatsUrl)
	if err != nil {
		log.Println(err)
		return false
	}

	for _, cookie := range resp.Cookies() {
		if cookie.Name == "cherpusername" {
			Username = cookie.Value
		}
	}

	var jsonData map[string]interface{}
	scanner := bufio.NewScanner(resp.Body)
	for i := 0; scanner.Scan() && i < 5; i++ {
		data := scanner.Text()
		json.Unmarshal([]byte(data), &jsonData)
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
		return false
	}

	if jsonData["status"].(string) == "failure" {
		return false
	}

	LoggedIn = true
	return true
}
