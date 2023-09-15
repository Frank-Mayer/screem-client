package utils

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"image"
	"image/png"
	"io"
	"log"
	"net"
	"os"
)

const (
	maxPacketSize = 2048
)

var (
	isSending = false
)

func Dial(host string, port string) net.Conn {
	serverAddr := fmt.Sprintf("%s:%s", host, port)

	conn, err := net.Dial("tcp", serverAddr)
	if err != nil {
		fmt.Println("Error connecting to server:", err)
		os.Exit(1)
	}

	fmt.Println("Connected to server at", serverAddr)

	return conn
}

func SendImageToServer(conn net.Conn, img *image.RGBA) {
	if isSending {
		return
	}
	isSending = true
	defer func() {
		isSending = false
	}()

	var buf bytes.Buffer
	png.Encode(&buf, img)

	bufferLen := int64(buf.Len())

	binary.Write(conn, binary.BigEndian, bufferLen)

	// n, err := io.Copy(conn, bytes.NewReader(buf.Bytes()))
	n, err := io.CopyN(conn, bytes.NewReader(buf.Bytes()), bufferLen)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Sent", n, "bytes to server.")
}
