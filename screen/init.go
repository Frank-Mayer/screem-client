package screen

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"image/png"
	"io"
	"log"
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
}

func InitGuest(client *utils.ScreemClient) {
	Client = client
	go backgroundLoopGuest()
}

func backgroundLoopHost() {
	for {
		captureScreen()
		time.Sleep(1000 * time.Millisecond)
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
			log.Fatal(err)
		}
		if n != size {
			log.Fatalf("Expected to read %d bytes, got %d\n", size, n)
		}
		img, err := png.Decode(buffer)
		if err != nil {
			log.Fatal(err)
		}
		ui.UpdateScreen(&img)
	}
}

func captureScreen() {
	if ui.ScreenNum < 0 {
		fmt.Println("No screen selected")
		return
	}

	bounds := screenshot.GetDisplayBounds(ui.ScreenNum)

	img, err := screenshot.CaptureRect(bounds)
	if err != nil {
		log.Println("Failed to capture screen:", err)
		return
	}

	ui.UpdateScreenPreview(img)
	fmt.Println("UI updated")
	Client.SendImageToServer(img)
	fmt.Println("Image sent to server")
}
