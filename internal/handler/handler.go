package handler

import (
	"bufio"
	"fmt"
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

	ackMsg := strings.ToUpper(strings.TrimSpace(message))
	response := fmt.Sprintf("ACK: %s\n", ackMsg)
	r, err := conn.Write([]byte(response))
	if err != nil {
		log.Printf("Server write error: %v", err)
	}

	fmt.Println(r)
}
