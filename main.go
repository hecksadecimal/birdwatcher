package main

import (
	"time"

	"cherp.chat/birdwatcher/cherp_api"
	"cherp.chat/birdwatcher/pages"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/systray"
)

var tabs *container.AppTabs

func UpdateSystrayIcon() {
	a := fyne.CurrentApp()
	if cherp_api.Unreads > 0 {
		if desk, ok := a.(desktop.App); ok {
			desk.SetSystemTrayIcon(resourceIconNotifPng)
		}
	} else {
		if desk, ok := a.(desktop.App); ok {
			desk.SetSystemTrayIcon(resourceIconPng)
		}
	}
}

func main() {
	a := app.NewWithID("chat.cherp.birdwatcher")
	w := a.NewWindow("Birdwatcher")
	w.Resize(fyne.NewSize(400, 600))

	cherp_api.InitializeClient()

	if desk, ok := a.(desktop.App); ok {
		m := fyne.NewMenu("Birdwatcher",
			fyne.NewMenuItem("Show", func() {
				w.Show()
			}))
		systray.SetTitle("Birdwatcher")
		desk.SetSystemTrayMenu(m)
		desk.SetSystemTrayIcon(resourceIconPng)
	}

	pages.GeneratePages()

	w.SetContent(pages.Tabs)

	w.SetCloseIntercept(func() {
		w.Hide()
	})

	ticker := time.NewTicker(time.Duration(pages.TickRate) * time.Second)
	go func() {
		for {
			select {
			case <-ticker.C:
				if cherp_api.LoggedIn {
					UpdateSystrayIcon()
				}
			case <-pages.QuitChannel:
				ticker.Stop()
				return
			}
		}
	}()

	w.ShowAndRun()
}
