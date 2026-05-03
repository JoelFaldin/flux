package parser

import (
	"fmt"
	"net"
)

type Command int

const (
	Set Command = iota
	Get
	Del
	PING
)

func Parser(cm string, conn net.Conn) {
	res := handleCommand(cm)

	// Handle health check:
	if res == 3 {
		response := "PONG"

		_, err := conn.Write([]byte(response))
		if err != nil {
			fmt.Printf("Server write error: %v", err)
		}
	}
}

func handleCommand(cmd string) Command {
	switch cmd {
	case "Set":
		return Set
	case "Get":
		return Get
	case "Del":
		return Del
	case "PING":
		return PING
	default:
		return 0
	}
}
