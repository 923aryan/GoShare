package main

import (
	"TerminalChat/cmd/server"
	"TerminalChat/pkg/connections"
	"TerminalChat/pkg/discovery"
	"TerminalChat/pkg/filetransfer"
	"TerminalChat/pkg/models"
	"TerminalChat/pkg/ui"
	"context"
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/widget"
)

func main() {

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	mainApp := app.New()
	myWindow := mainApp.NewWindow("Service Discovery")
	startButton, MainContainer := ui.MainWindow()

	DeviceContainer := ui.UiLayout(startButton, MainContainer)
	var sendButton *widget.Button
	var SendUi *fyne.Container
	var progressBar *widget.ProgressBar
	SendUi, sendButton, progressBar = ui.CreateProgressBar(MainContainer, mainApp, func() {
		// go filetransfer.SendFile("/home/aryan/Downloads/Win11_22H2_English_x64v2.iso", progressBar, wg)
	})

	sendButton.OnTapped = func() {
		fmt.Println("Modified behavior!")
	}

	Tabs := ui.CreateTabs(MainContainer, SendUi)

	var cancel context.CancelFunc = nil
	var ctx context.Context
	go connections.InitiateConnection()
	go server.Start()
	go func() {
		startButton.OnTapped = func() {
			if cancel != nil {
				cancel()
				cancel = nil
				DeviceContainer.RemoveAll()
				//models.Reset()
				time.Sleep(time.Millisecond * 100)
			}

			ctx, cancel = context.WithCancel(context.Background())
			wg := &sync.WaitGroup{}
			wg.Add(2)

			go discovery.DiscoverServices(wg, DeviceContainer, ctx, cancel)
			go discovery.RegisterServices(wg, DeviceContainer, ctx, cancel)
			sendButton.OnTapped = func() {
				go filetransfer.SendFile(models.SelectedPath, progressBar)
			}
			fmt.Println("done 1st iteration")
		}
	}()
	myWindow.SetOnClosed(func() {
		cancel()
		os.Exit(0)
	})
	myWindow.SetContent(Tabs)

	go func() {
		select {
		case <-models.ConnectionFormed:
			Tabs.Select(Tabs.Items[1])
			MainContainer.Refresh()
			cancel()
		case <-sigs:
			fmt.Println("Interrupt received, exiting...")
			cancel()
			//
			return
		}
	}()
	myWindow.ShowAndRun()
}
