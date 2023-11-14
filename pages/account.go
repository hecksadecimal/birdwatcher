package pages

import (
	"fmt"

	"cherp.chat/birdwatcher/cherp_api"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/widget"
)

var accountPage *container.TabItem

func GenerateAccountPage() {
	username := binding.NewString()
	username.Set("")
	usernameTextEntry := widget.NewEntry()

	password := binding.NewString()
	password.Set("")
	passwordTextEntry := widget.NewPasswordEntry()

	loginResultValue := binding.NewString()
	loginResultValue.Set("")
	loginResult := widget.NewLabelWithData(loginResultValue)

	loginPrompt := binding.NewString()
	loginPrompt.Set("Welcome, please log in!")

	loginForm := &widget.Form{
		Items: []*widget.FormItem{
			{Text: "Username", Widget: usernameTextEntry},
			{Text: "Password", Widget: passwordTextEntry},
		},
		OnSubmit: func() {
			// Handle login submission
			// You can add your own logic here to authenticate the user
			// For this example, we will simply display a success message
			username.Set(usernameTextEntry.Text)
			password.Set(passwordTextEntry.Text)
			if cherp_api.Login(usernameTextEntry.Text, passwordTextEntry.Text) {
				usernameTextEntry.SetText("")
				passwordTextEntry.SetText("")
				loginResultValue.Set("Successfully logged in!")
				loginPrompt.Set(fmt.Sprintf("Currently logged in as %s", cherp_api.Username))

				Tabs.EnableIndex(1)
				Tabs.EnableIndex(2)
			} else {
				loginResultValue.Set("Login failure.")
			}
		},
		SubmitText: "Log In",
		OnCancel: func() {
			cherp_api.Logout()
			loginResultValue.Set("Logged out.")
			loginPrompt.Set("Welcome, please log in!")

			Tabs.DisableIndex(1)
			Tabs.DisableIndex(2)
		},
		CancelText: "Log Out",
	}

	if cherp_api.SessionValid() {
		loginPrompt.Set(fmt.Sprintf("Currently logged in as %s", cherp_api.Username))
		cherp_api.InitNotifs()
	}

	loginPromptLabel := widget.NewLabelWithData(loginPrompt)

	accountContent := container.NewVBox(
		loginPromptLabel,
		loginForm,
		loginResult,
	)

	accountPage = container.NewTabItem("Account",
		accountContent,
	)
}
