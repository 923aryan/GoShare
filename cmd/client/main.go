package main

import (
	"TerminalChat/pkg/connections"
	"TerminalChat/pkg/discovery"
	"TerminalChat/pkg/models"
	"TerminalChat/pkg/ui"
	"context"
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"fyne.io/fyne/v2/app"
)

func main() {

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	mainApp := app.New()
	myWindow := mainApp.NewWindow("Service Discovery")
	startButton, MainContainer := ui.MainWindow()
	DeviceContainer := ui.UiLayout(startButton, MainContainer)
	Tabs := ui.CreateTabs(MainContainer)
	var cancel context.CancelFunc = nil
	var ctx context.Context
	go connections.InitiateConnection()
	startButton.OnTapped = func() {
		if cancel != nil {
			cancel()
			cancel = nil
			DeviceContainer.RemoveAll()
			time.Sleep(time.Millisecond * 100)
		}

		ctx, cancel = context.WithCancel(context.Background())

		wg := &sync.WaitGroup{}
		wg.Add(2)
		go discovery.DiscoverServices(wg, DeviceContainer, ctx, cancel)
		go discovery.RegisterServices(wg, DeviceContainer, ctx, cancel)

	}
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
			return
		}
	}()
	myWindow.ShowAndRun()
}
