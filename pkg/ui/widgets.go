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

func CreateTabs(MainContainer *fyne.Container) *container.AppTabs {

	tabs := container.NewAppTabs(
		container.NewTabItem("Discover PC", MainContainer),
		container.NewTabItem("Connected", ChatScreen()),
	)
	return tabs
}
