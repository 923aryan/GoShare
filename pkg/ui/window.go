package ui

import (
	"TerminalChat/pkg/models"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
)

func GetWidget() fyne.Widget {
	return widget.NewLabel("This is screen 1")
}

func ChatScreen() *fyne.Container {
	chatLog := widget.NewEntry()
	chatLog.MultiLine = true

	scrollableChatLog := container.NewScroll(chatLog)
	scrollableChatLog.SetMinSize(fyne.NewSize(300, 400))

	userInput := widget.NewEntry()

	sendButton := widget.NewButton("Send", func() {
		chatLog.SetText(chatLog.Text + userInput.Text + "\n")
		models.Message = userInput.Text + "\n"
		models.SendPressed <- true
		userInput.SetText("")
	})

	content := container.New(layout.NewVBoxLayout(),
		widget.NewLabel("Chat Log:"),
		scrollableChatLog,
		container.NewHBox(userInput, sendButton),
	)
	go func() {
		for {
			<-models.GotMessage
			chatLog.SetText(chatLog.Text + models.MessageReceived + "\n")
		}

	}()

	return content
}
