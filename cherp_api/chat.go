package cherp_api

import (
	"bufio"
	"encoding/json"
	"log"
	"os"

	"github.com/google/go-cmp/cmp"
)

var unreadChatsUrl = ApiUrl + "chat/list/unread"
var Unreads = 0
var NotifObject map[string]interface{}

func InitNotifs() {
	if _, err := os.Stat("./notifstate.json"); err == nil {
		Load("./notifstate.json", &NotifObject)
	}
	if NotifObject != nil {
		Unreads = len((NotifObject)["chats"].([]interface{}))
		if Unreads > 0 {
			TriggerUnreadNotification()
		}
	}
}

func TriggerUnreadNotification() {
	SendNotification("Cherp Chats", "You have new replies!")
}

func GetNotifications() {
	defer SaveJar()

	o := make(map[string]interface{})
	o = CopyMap(NotifObject)
	resp, err := client.Get(unreadChatsUrl)
	if err != nil {
		log.Println(err)
		return
	}

	scanner := bufio.NewScanner(resp.Body)
	for i := 0; scanner.Scan() && i < 5; i++ {
		json.Unmarshal([]byte(scanner.Text()), &NotifObject)
		// This'll let us restore our state if the program is closed and run again.
		Save("./notifstate.json", &NotifObject)
	}
	if err := scanner.Err(); err != nil {
		panic(err)
	}
	// On our first run, our notifstate would not have been saved already.
	// We'll initialize the value here and also check to see if there's something unread.
	if len(o) == 0 {
		o = CopyMap(NotifObject)
		Unreads = len((NotifObject)["chats"].([]interface{}))
		if Unreads > 0 {
			log.Println("New notification detected.")
			TriggerUnreadNotification()
		}
	}
	if !cmp.Equal(o, NotifObject) {
		// A simple comparison, if there's more unread chats than the last time we checked, trigger a notification.
		// TODO: Compare the two lists more thoroughly. It's possible to check a chat and get a new reply in another within the interval.
		Unreads = len((NotifObject)["chats"].([]interface{}))
		if Unreads > len(o["chats"].([]interface{})) {
			log.Println("New notification detected.")
			TriggerUnreadNotification()
		} else if Unreads < len(o["chats"].([]interface{})) {
			log.Println("You've checked on an unread reply.")
		}
	}
}
