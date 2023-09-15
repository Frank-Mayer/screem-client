package main

import (
	"fmt"
	"log"
	"os"

	"screem.frankmayer.io/screen"
	"screem.frankmayer.io/ui"
	"screem.frankmayer.io/utils"
)

func main() {
	argsCount := len(os.Args) - 1
	if argsCount != 3 {
		log.Fatalf(
			"Wrong number of arguments. Expected 3, got %d.\nExpected argument format: <host|join> <server address> <server port>\n",
			argsCount,
		)
	}

	mode := os.Args[1]
	host := os.Args[2]
	port := os.Args[3]

	switch mode {
	case "host":
		client := utils.NewClient(host, port)
        screen.InitHosting(client)
        ui.Host()
        client.Close()
	case "join":
		conn := utils.NewClient(host, port)
        defer conn.Close()
        screen.InitGuest(conn)
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
