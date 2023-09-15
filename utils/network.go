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
)

const (
	maxPacketSize = 2048
)

type ScreemClient struct {
	Conn net.Conn
    Ack chan bool
}

func NewClient(host string, port string) *ScreemClient {
	serverAddr := fmt.Sprintf("%s:%s", host, port)

	conn, err := net.Dial("tcp", serverAddr)
	if err != nil {
		log.Panicln("Error connecting to server:", err)
	}

	fmt.Println("Connected to server at", serverAddr)

	return &ScreemClient{
        Conn: conn,
        Ack: make(chan bool),
    }
}

func (self *ScreemClient) Close() {
	self.Conn.Close()
	self.Conn = nil
}

func (self *ScreemClient) SendImageToServer(img *image.RGBA) {
	var buf bytes.Buffer
	png.Encode(&buf, img)

	bufferLen := int64(buf.Len())

	binary.Write(self.Conn, binary.BigEndian, bufferLen)

	_, err := io.CopyN(self.Conn, &buf, bufferLen)
	if err != nil {
		log.Fatalln(err)
	}

	// await ack from server
    <-self.Ack
}
