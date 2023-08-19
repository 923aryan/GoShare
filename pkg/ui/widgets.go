package ui

import (
	// "TerminalChat/pkg/ui"

	"fmt"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	// "fyne.io/fyne/v2/container"
)

func MainWindow() (*widget.Button, *fyne.Container) {

	startButton := widget.NewButton("Start Discovery", nil)
	MainContainer := container.NewVBox(startButton)
	return startButton, MainContainer

}

func UiLayout(startButton *widget.Button, MainContainer *fyne.Container) *fyne.Container {

	DeviceContainer := container.NewVBox()
	MainContainer.Add(DeviceContainer)
	return DeviceContainer
}

func CreateProgressBar(MainContainer *fyne.Container, mainApp fyne.App, onTap func()) (*fyne.Container, *widget.Button, *widget.ProgressBar) {
	progressBar := widget.NewProgressBar()
	sendButton := widget.NewButton("Send File", onTap)
	sffButton := widget.NewButton("Select File", func() {
		sendButton.Disable()
		SelectFileWindow := mainApp.NewWindow("Select File")
		ShowFilePickerWindow(SelectFileWindow, func() {
			fmt.Println("reaching in callbakc")
			SelectFileWindow.Close()
			sendButton.Enable()
		})
		// SelectFileWindow.Close()

	})

	buttons := container.NewHBox(sendButton, sffButton)

	content := container.NewVBox(buttons, progressBar)
	return content, sendButton, progressBar
}

func CreateTabs(MainContainer *fyne.Container, SendUi *fyne.Container) *container.AppTabs {

	tabs := container.NewAppTabs(
		container.NewTabItem("Discover PC", MainContainer),
		container.NewTabItem("Connected", ChatScreen(SendUi)),
	)
	return tabs
}
