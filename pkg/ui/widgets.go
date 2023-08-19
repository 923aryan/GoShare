package ui

import (
	// "TerminalChat/pkg/ui"

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

func CreateProgressBar(MainContainer *fyne.Container, onTap func()) (*fyne.Container, *widget.Button, *widget.ProgressBar) {
	progressBar := widget.NewProgressBar()
	sendButton := widget.NewButton("Send File", onTap)
	content := container.NewVBox(sendButton, progressBar)
	return content, sendButton, progressBar
}

func CreateTabs(MainContainer *fyne.Container, SendUi *fyne.Container) *container.AppTabs {

	tabs := container.NewAppTabs(
		container.NewTabItem("Discover PC", MainContainer),
		container.NewTabItem("Connected", ChatScreen(SendUi)),
	)
	return tabs
}
