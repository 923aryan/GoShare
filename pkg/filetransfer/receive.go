package filetransfer

import (
	"TerminalChat/pkg/models"
	"fmt"
	"sync"

	"github.com/gorilla/websocket"
)

func ReceiveText(conn *websocket.Conn, wg *sync.WaitGroup) {
	go func() {
		defer wg.Done()
		for {
			_, message, err := conn.ReadMessage()

			if err != nil {
				fmt.Println("read:", err)
				return
			}
			models.GotMessage <- true
			models.MessageReceived = string(message)
			fmt.Println("recvInServer", string(message))

			if err != nil {
				fmt.Println("write", err)
				break
			}
		}
		fmt.Println("yesssss")
	}()

}
