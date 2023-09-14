package utils

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"image"
	"image/png"
	"net"
)

const (
	maxIPv4Packet = 1400 // Recommended size for IPv4
	maxIPv6Packet = 1230 // Recommended size for IPv6
)

func SendImageToServer(conn net.Conn, img *image.RGBA) {
	// Encode the image to PNG format
	var buf bytes.Buffer
	if err := png.Encode(&buf, img); err != nil {
		fmt.Println("Error encoding image:", err)
		return
	}

	// Determine the maximum packet size based on IPv4 or IPv6
	maxPacketSize := maxIPv4Packet
	serverIP := conn.RemoteAddr().(*net.TCPAddr).IP
	if serverIP.To4() == nil {
		maxPacketSize = maxIPv6Packet
	}

	// Split the image data into smaller packets and send them
	imageBytes := buf.Bytes()
	totalSize := len(imageBytes)
	offset := 0

	// Send the total size of the image to the server using binary.BigEndian.PutUint32
	{
		packetSize := make([]byte, 4)
		binary.BigEndian.PutUint32(packetSize, uint32(totalSize))
		_, err := conn.Write(packetSize)
		if err != nil {
			fmt.Println("Error sending image packet to server:", err)
			return
		}
	}

	for offset < totalSize {
		end := offset + maxPacketSize
		if end > totalSize {
			end = totalSize
		}

		packet := imageBytes[offset:end]

		_, err := conn.Write(packet)
		if err != nil {
			fmt.Println("Error sending image packet to server:", err)
			return
		}

		offset = end
	}

	// Send a zero-length packet to indicate the end of the image
	{
		_, err := conn.Write([]byte{})
		if err != nil {
			fmt.Println("Error sending image packet to server:", err)
			return
		}
	}

	fmt.Println("Sent image to server")
}
