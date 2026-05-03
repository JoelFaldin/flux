package parser

import (
	"fmt"
	"net"
	"strings"

	"flux/internal/store"
)

type Command int

const (
	Set Command = iota
	Get
	Del
	PING
)

func Parser(conn net.Conn, cm []string, s *store.Store) {
	res := handleCommand(cm[0])

	// Handle health check:
	if res == 3 {
		response := "PONG"

		_, err := conn.Write([]byte(response))
		if err != nil {
			fmt.Printf("Server write error: %v", err)
		}
	}

	// Handle "Get":
	if res == 1 {
		key := cm[1]
		before, _, _ := strings.Cut(key, "\r\n")
		format := before

		val, ok := s.GetValue(format)
		if !ok {
			conn.Write([]byte("(nil)"))
		}

		conn.Write([]byte(val))

		return
	}

	// Handle "Set":
	if res == 0 {
		key := cm[1]
		val := cm[2]
		s.SetValue(key, val)

		conn.Write([]byte("OK"))
		return
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
