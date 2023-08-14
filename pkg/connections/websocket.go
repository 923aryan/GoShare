package connections

import (
	"TerminalChat/pkg/filetransfer"
	"TerminalChat/pkg/models"
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
	//TODO

	wg := &sync.WaitGroup{}
	wg.Add(2)
	go filetransfer.SendText(conn, wg)
	go filetransfer.ReceiveText(conn, wg)
	wg.Wait()

}

func InitiateConnection() {
	http.HandleFunc("/socket", Handler)
	log.Fatal(http.ListenAndServe(":8000", nil))
}

func ClientConnection(entry *zeroconf.ServiceEntry, send chan int) {

	interrupt := make(chan os.Signal, 1)

	signal.Notify(interrupt, os.Interrupt)
	u := url.URL{Scheme: "ws",
		Host: fmt.Sprintf("%s:%d", entry.AddrIPv4[0], entry.Port), Path: "/socket"}
	conn, _, err := websocket.DefaultDialer.Dial(u.String(), nil)

	if err != nil {
		fmt.Println("Websocket connection failed", err)
		return
	}
	models.ConnectionFormed <- true
	models.ConnectionEstablished = true
	defer conn.Close()

	done := make(chan struct{})

	// Function to read message
	wg := &sync.WaitGroup{}
	wg.Add(2)
	go filetransfer.ReceiveText(conn, wg)
	go filetransfer.SendText(conn, wg)
	wg.Wait()
	// go func() {
	// 	defer close(done)
	// 	for {
	// 		_, message, err := conn.ReadMessage()
	// 		if err != nil {
	// 			fmt.Println("read:", err)
	// 			return
	// 		}
	// 		fmt.Printf("recv: %s", string(message))
	// 	}
	// }()

	// go func() {
	// 	reader := bufio.NewReader(os.Stdin)
	// 	for {
	// 		fmt.Println("Enter Message: ")
	// 		msg, _ := reader.ReadString('\n')
	// 		msg = strings.TrimSpace(msg)
	// 		err := conn.WriteMessage(websocket.TextMessage, []byte(msg))
	// 		if err != nil {
	// 			fmt.Println("write:", err)
	// 			return
	// 		}
	// 	}
	// }()

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
