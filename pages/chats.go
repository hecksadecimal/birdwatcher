package pages

import (
	"fmt"
	"net/url"

	"cherp.chat/birdwatcher/cherp_api"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/widget"
)

var chatsPage *container.TabItem
var chatsTickRate = 15
var chatsUnread binding.String
var unreadChatsContainer *fyne.Container

func unreadChatsPluralize(n int) string {
	if n == 1 {
		return "chat"
	} else {
		return "chats"
	}
}

func setUnreadsLabel() {
	chatsUnread.Set(fmt.Sprintf("You have %d unread %s", cherp_api.Unreads, unreadChatsPluralize(cherp_api.Unreads)))
}

func updateCards() {
	unreadChatsContainer.RemoveAll()
	if cherp_api.NotifObject == nil {
		return
	}
	chats := cherp_api.NotifObject["chats"].([]interface{})
	for _, element := range chats {
		chatURL := element.(map[string]interface{})["chatURL"].(string)
		chatTitle := ""
		chatDescription := ""

		if element.(map[string]interface{})["chatTitle"] != nil {
			chatTitle = element.(map[string]interface{})["chatTitle"].(string)
		}

		if element.(map[string]interface{})["chatDescription"] != nil {
			chatDescription = element.(map[string]interface{})["chatDescription"].(string)
		}

		if chatTitle == "" {
			chatTitle = chatURL
		}

		hyperlink, err := url.Parse(cherp_api.BaseUrl + "chats/" + chatURL)
		if err == nil {
			unreadChatsContainer.Add(
				widget.NewCard(
					chatTitle,
					chatDescription,
					widget.NewHyperlink("Open Chat", hyperlink),
				),
			)
		}
	}
	unreadChatsContainer.Refresh()
}

func GenerateChatsPage() {
	chatsUnread = binding.NewString()
	setUnreadsLabel()
	unreadChatsContainer = container.NewVBox()
	updateCards()

	chatsContent := container.NewVBox(
		widget.NewLabelWithData(chatsUnread),
		widget.NewSeparator(),
		unreadChatsContainer,
	)
	chatsPage = container.NewTabItem("Chats",
		chatsContent,
	)
}

func TickChats() {
	if CurrentTick%chatsTickRate == 0 && cherp_api.LoggedIn {
		cherp_api.GetNotifications()
		setUnreadsLabel()
		updateCards()
	}
}
