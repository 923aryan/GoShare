package server

import (
	"encoding/binary"
	"fmt"
	"io"
	"log"
	"net"
	"os"
)

func Start() {

	listener, err := net.Listen("tcp", ":8108")
	if err != nil {
		log.Fatalf("Error: %v", err)
	}

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Fatalf("Connection error: %v", err)
		}
		go handleConnection(conn)
	}
	fmt.Println("reached end")
}

func handleConnection(conn net.Conn) {
	defer conn.Close()

	var fileSize int64
	err := binary.Read(conn, binary.BigEndian, &fileSize)
	if err != nil {
		log.Fatalf("Reading size error: %v", err)
	}

	outputFile, err := os.Create("received_file.iso")
	if err != nil {
		log.Fatalf("File create error: %v", err)
	}
	defer outputFile.Close()

	_, err = io.CopyN(outputFile, conn, fileSize)
	if err != nil {
		log.Fatalf("File write error: %v", err)
	}
}
