package connections

import (
	"TerminalChat/pkg/filetransfer"
	"TerminalChat/pkg/models"
	"context"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"sync"
	"time"

	"github.com/gorilla/websocket"
	"github.com/grandcat/zeroconf"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func Handler(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println("WebSocket upgrade failed:", err)
		return
	}

	fmt.Println("Connection established")
	models.ConnectionEstablished = true
	models.ConnectionFormed <- true
	defer conn.Close()
	wg := &sync.WaitGroup{}
	wg.Add(2)
	go filetransfer.SendText(conn, wg)
	go filetransfer.ReceiveText(conn, wg)
	wg.Wait()
	go func() {
		wg.Wait()
		<-models.ConnectionAborted
		conn.Close()
	}()
}

func InitiateConnection() {
	http.HandleFunc("/socket", Handler)
	log.Fatal(http.ListenAndServe(":8000", nil))
}

func ClientConnection(entry *zeroconf.ServiceEntry, send chan int, ctx context.Context) {

	interrupt := make(chan os.Signal, 1)

	signal.Notify(interrupt, os.Interrupt)
	u := url.URL{Scheme: "ws",
		Host: fmt.Sprintf("%s:%d", entry.AddrIPv4[0], entry.Port), Path: "/socket"}
	conn, _, err := websocket.DefaultDialer.Dial(u.String(), nil)

	if err != nil {
		fmt.Println("Websocket connection failed", err)
		return
	}
	models.Entry = entry
	models.ConnectionFormed <- true
	models.ConnectionEstablished = true

	defer conn.Close()

	done := make(chan struct{})

	wg := &sync.WaitGroup{}
	wg.Add(2)
	defer wg.Done()
	go filetransfer.ReceiveText(conn, wg)
	go filetransfer.SendText(conn, wg)
	wg.Wait()

	for {
		select {
		case <-done:
			return
		case <-interrupt:
			fmt.Println("Interrupt")
			go func() {
				time.Sleep(time.Second * 1)
				err := conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
				if err != nil {
					fmt.Println("write clsoe", err)
					return
				}

			}()

			<-done
			return
		}
	}

}
