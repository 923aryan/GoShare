package filetransfer

import (
	"TerminalChat/pkg/models"
	"encoding/binary"
	"fmt"
	"io"
	"net"
	"os"
	"sync"

	"fyne.io/fyne/v2/widget"
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

func SendFile(filepath string, progressBar *widget.ProgressBar, wg *sync.WaitGroup) {
	defer wg.Done()
	file, err := os.Open(filepath)
	if err != nil {
		fmt.Println("File open error:", err)
		return
	}
	defer file.Close()

	fileInfo, err := file.Stat()
	if err != nil {
		fmt.Println("File stat error:", err)
		return
	}
	fmt.Println((models.Entry.AddrIPv4[0]).String())
	conn, err := net.Dial("tcp", (models.Entry.AddrIPv4[0]).String()+":8108")
	if err != nil {
		fmt.Println("Dial error:", err)
		return
	}

	defer conn.Close()

	err = binary.Write(conn, binary.BigEndian, fileInfo.Size())
	if err != nil {
		fmt.Println("Size send error:", err)
		return
	}

	sentBytes := int64(0)
	buffer := make([]byte, 1024)
	for {
		bytesRead, err := file.Read(buffer)
		if err != nil && err != io.EOF {
			fmt.Println("File read error:", err)
			return
		}
		if bytesRead == 0 {
			break
		}

		_, err = conn.Write(buffer[:bytesRead])
		if err != nil {
			fmt.Println("Send error:", err)
			return
		}

		sentBytes += int64(bytesRead)
		progressBar.SetValue(float64(sentBytes) / float64(fileInfo.Size()))
	}

}
