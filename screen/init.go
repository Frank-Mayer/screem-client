package screen

import (
	"log"
	"net"
	"time"

	"github.com/kbinani/screenshot"
	"screem.frankmayer.io/ui"
	"screem.frankmayer.io/utils"
)

var (
	StopHosting chan bool
    Conn *net.Conn = nil
)

func InitHosting(conn *net.Conn) {
    Conn = conn
	StopHosting = make(chan bool, 1)
	go backgroundLoop()
}

func backgroundLoop() {
	for {
		select {
		case <-StopHosting:
			log.Println("Stopping hosting")
			return
		default:
			captureScreen()
		}
		time.Sleep(1000 * time.Millisecond)
	}
}

func captureScreen() {
	if ui.ScreenNum < 0 {
		return
	}

	bounds := screenshot.GetDisplayBounds(ui.ScreenNum)

	img, err := screenshot.CaptureRect(bounds)
	if err != nil {
		log.Println("Failed to capture screen:", err)
		return
	}

	ui.UpdateScreenPreview(img)

    utils.SendImageToServer(*Conn, img)
}
