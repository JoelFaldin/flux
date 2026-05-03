package handler

import (
	"bufio"
	"flux/internal/parser"
	"log"
	"net"
	"strings"
)

// Each new connection spawns a new gorouting that executes this function.
func Handler(conn net.Conn) {
	defer conn.Close()

	reader := bufio.NewReader(conn)
	message, err := reader.ReadString('\n')
	if err != nil {
		log.Printf("Read error: %v", err)
		return
	}

	res := strings.Split(message, " ")
	parser.Parser(res[0], conn)

	// response := fmt.Sprintf("ACK: %s\n", ackMsg)
	// if err != nil {
	// 	log.Printf("Server write error: %v", err)
	// }
}
