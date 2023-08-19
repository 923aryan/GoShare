// package main

// import (
// 	"fyne.io/fyne/v2"
// 	"fyne.io/fyne/v2/app"
// 	"fyne.io/fyne/v2/container"
// 	"fyne.io/fyne/v2/theme"
// 	"fyne.io/fyne/v2/widget"
// )

// func main() {
// 	myApp := app.New()
// 	myWindow := myApp.NewWindow("Fyne Custom Buttons")

// 	// Real-Time Communication Button
// 	rtCommButton := widget.NewButtonWithIcon("Real-Time Communication", theme.AccountIcon(), func() {
// 		// action when pressed
// 	})
// 	rtCommButton.Importance = widget.HighImportance

// 	// Data Streaming Button
// 	dataStreamingButton := widget.NewButtonWithIcon("Data Streaming", theme.StorageIcon(), func() {
// 		// action when pressed
// 	})
// 	dataStreamingButton.Importance = widget.HighImportance

//		// Layout them vertically for now
//		content := container.NewVBox(rtCommButton, dataStreamingButton)
//		myWindow.SetContent(content)
//		myWindow.Resize(fyne.NewSize(300, 200))
//		myWindow.ShowAndRun()
//	}
package main

import (
	"fmt"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
)

func main() {
	myApp := app.New()
	myWindow := myApp.NewWindow("File/Folder Picker")

	var selectedPath string

	// Button function
	btnFunc := func() {
		dialog.ShowFileOpen(func(file fyne.URIReadCloser, err error) {
			if err == nil && file != nil {
				selectedPath = file.URI().Path()
				fmt.Println("Selected file/folder path:", selectedPath)
			} else if err != nil {
				dialog.ShowError(err, myWindow)
			}
		}, myWindow)
	}

	btn := widget.NewButton("Open File/Folder", btnFunc)
	myWindow.SetContent(container.NewCenter(btn))

	myWindow.ShowAndRun()
}
