package parser

import (
	"fmt"
	"net"
	"strings"
	"time"

	"flux/internal/loader"
	"flux/internal/models"
	"flux/internal/store"
)

type Command int

const (
	Set Command = iota
	Get
	Del
	PING
	LPUSH
)

func Parser(conn net.Conn, cm []string, s *store.Store, globalConfig *models.Data) {
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

			before := strings.TrimSpace(cm[4])
			format := fmt.Sprintf("%ss", before)

			res, err := time.ParseDuration(format)
			if err != nil {
				conn.Write([]byte("Error: invalid duration"))
			}

			s.SetTemporalValue(key, val, res)
		} else {
			key := cm[1]
			val := cm[2]
			cutVal := strings.TrimSpace(val)

			s.SetValue(key, cutVal)
		}

		// Writting to yaml:
		loader.WriteData(globalConfig, s.GetAllValues())

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
		format := strings.TrimSpace(key)

		val, ok := s.GetValue(format)
		if !ok {
			conn.Write([]byte("(nil)"))
		}

		// Type switch:
		switch v := val.(type) {
		case string:
			conn.Write([]byte(v + "\r\n"))
		case []string:
			result := strings.Join(v, " ")
			conn.Write([]byte("[" + result + "]\r\n"))
		case []any:
			var strItems []string
			for _, item := range v {
				strItems = append(strItems, fmt.Sprintf("%v", item))
			}

			conn.Write([]byte("[" + strings.Join(strItems, ", ") + "]\r\n"))
		default:
			fmt.Fprintf(conn, "%v\r\n", v)
		}

		return
	}

	// Handle "Del":
	if res == 2 {
		if len(cm) < 2 {
			conn.Write([]byte("Error: wrong number of arguments for DEL"))
			return
		}

		key := cm[1]
		before := strings.TrimSpace(key)
		format := before

		s.DeleteValue(format)

		// Writing without the deleted value:
		loader.WriteData(globalConfig, s.GetAllValues())

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

	// Handle LPUSH:
	if res == 4 {
		if len(cm) < 3 {
			conn.Write([]byte("Error: wrong number of arguments for LPUSH"))
			return
		}

		key := cm[1]
		val := cm[2]
		cutVal := strings.TrimSpace(val)

		var sl []string
		sl = append(sl, cutVal)

		s.LPush(key, sl)

		d := s.GetAllValues()
		loader.WriteData(globalConfig, d)

		conn.Write([]byte("OK\r\n"))
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
	case "LPUSH":
		return LPUSH
	default:
		return -1
	}
}
