package filetransfer

import (
	"TerminalChat/pkg/models"
	"fmt"
	"sync"

	"github.com/gorilla/websocket"
)

func SendText(conn *websocket.Conn, wg *sync.WaitGroup) {
	go func() {
		defer wg.Done()

		for {

			select {
			case <-models.SendPressed:
				fmt.Println("haha")
				err := conn.WriteMessage(websocket.TextMessage, []byte(models.Message))

				if err != nil {
					fmt.Println("write:", err)
					return
				}

			}
			fmt.Println("done")
		}

	}()
}
