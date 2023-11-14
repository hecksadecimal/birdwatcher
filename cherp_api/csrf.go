package cherp_api

import (
	"bufio"
	"encoding/json"
	"log"
	"net/url"
)

// Retrieves the current CSRF values for a given page containing a form.
func GetCsrf(formUrl string) (data map[string]interface{}) {
	defer SaveJar()

	u, err := url.Parse(BaseUrl + formUrl)
	if err != nil {
		log.Fatal(err)
		return
	}

	_, pageErr := client.Get(u.String())
	if pageErr != nil {
		log.Fatal(err)
		return
	}

	resp, err := client.Get(ApiUrl + "csrf")
	if err != nil {
		log.Println(err)
		return
	}
	var jsonData map[string]interface{}
	scanner := bufio.NewScanner(resp.Body)
	for i := 0; scanner.Scan() && i < 5; i++ {
		data := scanner.Text()
		json.Unmarshal([]byte(data), &jsonData)
	}
	if err := scanner.Err(); err != nil {
		panic(err)
	}

	return jsonData
}
