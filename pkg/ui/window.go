package ui

import (
	"TerminalChat/pkg/models"
	"fmt"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
)

func GetWidget() fyne.Widget {
	return widget.NewLabel("This is screen 1")
}

func ChatScreen(SendUi *fyne.Container) *fyne.Container {
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
		SendUi,
	)
	go func() {
		for {
			<-models.GotMessage
			chatLog.SetText(chatLog.Text + models.MessageReceived + "\n")
		}

	}()

	return content
}

func ShowFilePickerWindow(selctFileWindow fyne.Window, CallBack func()) {

	var selectedPath string
	btnFunc := func() {
		dialog.ShowFileOpen(func(file fyne.URIReadCloser, err error) {
			if err == nil && file != nil {
				selectedPath = file.URI().Path()
				fmt.Println("Selected file/folder path:", selectedPath)
				models.SelectedPath = selectedPath
				CallBack()
			} else if err != nil {
				dialog.ShowError(err, selctFileWindow)
			}
		}, selctFileWindow)
	}
	btn := widget.NewButton("Open File/Folder", btnFunc)
	selctFileWindow.SetContent(container.NewCenter(btn))
	selctFileWindow.Show()
}
