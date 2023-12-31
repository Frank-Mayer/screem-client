package screen

import (
	"bytes"
	"compress/zlib"
	"encoding/binary"
	"fmt"
	"image/png"
	"io"
	"log"
	"os"
	"time"

	"github.com/kbinani/screenshot"
	"screem.frankmayer.io/ui"
	"screem.frankmayer.io/utils"
)

var (
	Client *utils.ScreemClient
)

func InitHosting(client *utils.ScreemClient) {
	Client = client
	go backgroundLoopHost()
	go backgroundAckLoop()
}

func InitGuest(client *utils.ScreemClient) {
	Client = client
	go backgroundLoopGuest()
}

func backgroundLoopHost() {
	for {
		captureScreen()
		time.Sleep(16 * time.Millisecond)
	}
}

func backgroundAckLoop() {
	for {
		var ack bool
		binary.Read(Client.Conn, binary.BigEndian, &ack)
		if ack {
			Client.Ack <- true
		}
	}
}

func backgroundLoopGuest() {
	var size int64
	for {
		binary.Read(Client.Conn, binary.BigEndian, &size)
		if size == 0 {
			continue
		}
		buffer := new(bytes.Buffer)
		n, err := io.CopyN(buffer, Client.Conn, size)
		if err != nil {
			fmt.Println(err)
            os.Exit(1)
		}
		if n != size {
			log.Panicf("Expected to read %d bytes, got %d\n", size, n)
		}

		// decompress image
		r, err := zlib.NewReader(buffer)
		defer r.Close()
		if err != nil {
			log.Panicln("Error decompressing image:", err)
		}
		uncompressed := new(bytes.Buffer)
		_, err = io.Copy(uncompressed, r)
		if err != nil {
			log.Panicln("Error decompressing image:", err)
		}

		img, err := png.Decode(uncompressed)
		if err != nil {
			log.Panicln("Error decoding image:", err)
		}
		ui.UpdateScreen(&img)

		// send OK to server to continue
		binary.Write(Client.Conn, binary.BigEndian, true)
	}
}

func captureScreen() {
	if ui.ScreenNum < 0 {
		return
	}

	bounds := screenshot.GetDisplayBounds(ui.ScreenNum)

	img, err := screenshot.CaptureRect(bounds)
	if err != nil {
		return
	}

	ui.UpdateScreenPreview(img)
	Client.SendImageToServer(img)
}
