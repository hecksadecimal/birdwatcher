package cherp_api

import (
	"bufio"
	"encoding/json"
	"log"
)

var taglistUrl = ApiUrl + "user/taglist/get"
var searchUrl = ApiUrl + "directory/smarttagged/"

func GetTagList() (data map[string]interface{}) {
	resp, err := client.Get(taglistUrl)
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

func GetSearchResults(query string) (data map[string]interface{}) {
	resp, err := client.Get(searchUrl + query)
	if err != nil {
		log.Println(err)
		return
	}
	var jsonData map[string]interface{}
	scanner := bufio.NewScanner(resp.Body)
	buf := []byte{}
	scanner.Buffer(buf, 2048*1024)
	for i := 0; scanner.Scan() && i < 5; i++ {
		data := scanner.Text()
		json.Unmarshal([]byte(data), &jsonData)
	}
	if err := scanner.Err(); err != nil {
		panic(err)
	}

	return jsonData
}
