package discovery

import (
	"TerminalChat/pkg/connections"
	"TerminalChat/pkg/models"
	"context"
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/widget"
	"github.com/grandcat/zeroconf"
)

func DiscoverServices(wg *sync.WaitGroup, content *fyne.Container, ctx context.Context, cancel context.CancelFunc) {
	defer wg.Done()
	var status *widget.Label = widget.NewLabel("")
	resolver, err := zeroconf.NewResolver(nil)
	if err != nil {
		fmt.Println("Failed to initialize resolver", err.Error())
		status.SetText(err.Error())
		status.Refresh()
		return
	}
	notify := make(chan int)
	entries := make(chan *zeroconf.ServiceEntry)
	go func(results <-chan *zeroconf.ServiceEntry) {
		for entry := range results {
			if entry.TTL == 0 {
				fmt.Println("Service gone:", entry.ServiceInstanceName())

			} else {
				fmt.Println("Found service:", entry.ServiceInstanceName())
			}
			//fmt.Printf("Found services: %v\n", entry)

			NewLabel := widget.NewLabel(entry.HostName)
			clickable := widget.NewButton(NewLabel.Text, nil)
			clickable.OnTapped = func() {
				fmt.Println("clicked", entry.AddrIPv4)
				go connections.ClientConnection(entry, notify)

			}

			models.Mu.Lock()
			models.ServiceEntries = append(models.ServiceEntries, entry)
			clickable.SetIcon(clickable.Icon)
			content.Add(clickable)
			models.Mu.Unlock()

			content.Refresh()
		}
		if !models.ConnectionEstablished {
			fmt.Println("No more entries")
		}

	}(entries)

	status.SetText("Finding")
	content.Refresh()
	err = resolver.Browse(ctx, "Goshare._tcp", "local.", entries)
	if err != nil {
		fmt.Println("Failed to browse:", err.Error())
		status.SetText(err.Error())
		status.Refresh()
		return
	}
	select {
	case <-ctx.Done():
		if !models.ConnectionEstablished {
			fmt.Println("DiscoverService Cancelled")
			return
		}
		return
	case <-time.After(time.Second * 20):
		if !models.ConnectionEstablished {
			fmt.Println("DisoverService Timeout")
			return
		}
		return

	}
}

func RegisterServices(wg *sync.WaitGroup, content *fyne.Container, ctx context.Context, cancel context.CancelFunc) {
	defer wg.Done()
	defer cancel()
	var status *widget.Label = widget.NewLabel("")
	hostname, _ := os.Hostname()
	service, err := zeroconf.Register(
		hostname,
		"Goshare._tcp",
		"local.",
		8000,
		nil,
		nil)

	if err != nil {
		fmt.Println("Failed to register service", err.Error())
		status.SetText(err.Error())
		status.Refresh()
		return
	}

	defer service.Shutdown()

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt, syscall.SIGTERM)

	select {
	case <-sig:
		status.SetText("Process Cancelled")
		status.Refresh()
	case <-ctx.Done():
		if !models.ConnectionEstablished {
			fmt.Println("RegisterService Cancelled")
			return
		}
		return
	case <-time.After(time.Second * 20):
		if !models.ConnectionEstablished {
			fmt.Println("RegisterService Timeout")
			return
		}
		return
	}

}
