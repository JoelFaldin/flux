package parser

import (
	"fmt"
	"net"
	"strings"
	"time"

	"flux/internal/loader"
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

	// Handle default value (-1):
	if res == -1 {
		r := fmt.Sprintf("Error: unknown command '%s'", cm[0])
		conn.Write([]byte(r))
	}

	// Handle "Set":
	if res == 0 {
		if len(cm) < 3 {
			conn.Write([]byte("Error: wrong number of arguments for SET"))
			return
		}

		// Checking for temporal value:
		if len(cm) == 5 {
			if cm[3] != "EX" {
				conn.Write([]byte("Error: invalid action on SET"))
				return
			}

			key := cm[1]
			val := cm[2]

			before, _, _ := strings.Cut(cm[4], "\r\n")
			format := fmt.Sprintf("%ss", before)

			res, err := time.ParseDuration(format)
			if err != nil {
				conn.Write([]byte("Error: invalid duration"))
			}

			s.SetTemporalValue(key, val, res)
		}

		key := cm[1]
		val := cm[2]
		cutVal, _, _ := strings.Cut(val, "\r\n")
		s.SetValue(key, cutVal)

		// Writting to yaml:
		loader.WriteData(loader.Data{Storage: map[string]string{key: val}}, s.GetAllValues())

		conn.Write([]byte("OK\r\n"))
		return
	}

	// Handle "Get":
	if res == 1 {
		if len(cm) < 2 {
			conn.Write([]byte("Error: wrong number of arguments for GET"))
			return
		}

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

	// Handle "Del":
	if res == 2 {
		if len(cm) < 2 {
			conn.Write([]byte("Error: wrong number of arguments for DEL"))
			return
		}

		key := cm[1]
		before, _, _ := strings.Cut(key, "\r\n")
		format := before

		s.DeleteValue(format)

		conn.Write([]byte("OK\r\n"))
	}

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
	case "PING\r\n":
		return PING
	default:
		return -1
	}
}
