package pages

import (
	"time"

	"cherp.chat/birdwatcher/cherp_api"
	"fyne.io/fyne/v2/container"
)

var Pages []*container.TabItem
var Tabs *container.AppTabs
var TickRate = 1 //Backround tasks are ticked every second.
var CurrentTick = 0
var QuitChannel chan struct{}

func GeneratePages() {
	GenerateAccountPage()
	GenerateChatsPage()
	GenerateTagsPage()

	Tabs = container.NewAppTabs(
		accountPage,
		chatsPage,
		tagsPage,
	)

	if !cherp_api.SessionValid() {
		Tabs.DisableIndex(1)
		Tabs.DisableIndex(2)
	}

	Tabs.SetTabLocation(container.TabLocationLeading)

	ticker := time.NewTicker(time.Duration(TickRate) * time.Second)
	QuitChannel = make(chan struct{})
	go func() {
		for {
			select {
			case <-ticker.C:
				if cherp_api.LoggedIn {
					TickChats()
					TickTags()
					CurrentTick++
				} else {
					CurrentTick = 0
				}
			case <-QuitChannel:
				ticker.Stop()
				return
			}
		}
	}()
}
