package main

import (
	"fmt"
	"os"
    "net"

	"screem.frankmayer.io/screen"
	"screem.frankmayer.io/ui"
)

func main() {
	argsCount := len(os.Args) - 1
	if argsCount != 3 {
		fmt.Fprintf(
			os.Stderr,
			"Wrong number of arguments. Expected 3, got %d.\nExpected argument format: <host|join> <server address> <server port>\n",
			argsCount,
		)
		os.Exit(1)
		return
	}

	mode := os.Args[1]
    host := os.Args[2]
    port := os.Args[3]

    serverAddr := fmt.Sprintf("%s:%s", host, port)

    conn, err := net.Dial("tcp", serverAddr)
    if err != nil {
        fmt.Println("Error connecting to server:", err)
        os.Exit(1)
    }
    defer conn.Close()

    fmt.Println("Connected to server at", serverAddr)

	switch mode {
	case "host":
		screen.InitHosting(&conn)
		ui.Host()
	case "join":
		ui.Guest()
	default:
		fmt.Fprintf(
			os.Stderr,
			"Unknown mode \"%s\". Available modes are 'host' and 'join'.\n",
			mode,
		)
		os.Exit(1)
	}
}
