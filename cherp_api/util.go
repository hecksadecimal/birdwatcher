package cherp_api

import (
	"fyne.io/fyne/v2"
)

func CopyMap(m map[string]interface{}) map[string]interface{} {
	cp := make(map[string]interface{})
	for k, v := range m {
		vm, ok := v.(map[string]interface{})
		if ok {
			cp[k] = CopyMap(vm)
		} else {
			cp[k] = v
		}
	}

	return cp
}

func SendNotification(title string, content string) {
	notif := fyne.NewNotification(title, content)
	fyne.CurrentApp().SendNotification(notif)
}
