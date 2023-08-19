package server

import (
	"encoding/binary"
	"io"
	"log"
	"net"
	"os"
	"path/filepath"
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

}

// func handleConnection(conn net.Conn) {
// 	defer conn.Close()

// 	var fileSize int64
// 	err := binary.Read(conn, binary.BigEndian, &fileSize)
// 	if err != nil {
// 		log.Fatalf("Reading size error: %v", err)
// 	}

// 	outputFile, err := os.Create("received_file.iso")
// 	if err != nil {
// 		log.Fatalf("File create error: %v", err)
// 	}
// 	defer outputFile.Close()

//		_, err = io.CopyN(outputFile, conn, fileSize)
//		if err != nil {
//			log.Fatalf("File write error: %v", err)
//		}
//	}
func handleConnection(conn net.Conn) {
	defer conn.Close()

	var filenameLen int32
	if err := binary.Read(conn, binary.BigEndian, &filenameLen); err != nil {
		log.Fatalf("Error reading filename length: %v", err)
	}

	filenameBuffer := make([]byte, filenameLen)
	if _, err := io.ReadFull(conn, filenameBuffer); err != nil {
		log.Fatalf("Error reading filename: %v", err)
	}
	filename := string(filenameBuffer)

	var fileSize int64
	if err := binary.Read(conn, binary.BigEndian, &fileSize); err != nil {
		log.Fatalf("Reading size error: %v", err)
	}

	folderPath := "./" + filename
	if err := os.MkdirAll(folderPath, os.ModePerm); err != nil {
		log.Fatalf("Error creating directory: %v", err)
	}

	outputFilePath := filepath.Join(folderPath, filename)
	outputFile, err := os.Create(outputFilePath)
	if err != nil {
		log.Fatalf("File create error: %v", err)
	}
	defer outputFile.Close()

	_, err = io.CopyN(outputFile, conn, fileSize)
	if err != nil {
		log.Fatalf("File write error: %v", err)
	}
}
