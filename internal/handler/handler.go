package handler

import (
	"bufio"
	"flux/internal/models"
	"flux/internal/parser"
	"flux/internal/store"
	"log"
	"net"
	"strings"
)

// Each new connection spawns a new gorouting that executes this function.
func Handler(conn net.Conn, s *store.Store, conf *models.Data) {
	defer conn.Close()

	reader := bufio.NewReader(conn)
	message, err := reader.ReadString('\n')
	if err != nil {
		log.Printf("Read error: %v", err)
		return
	}

	res := strings.Split(message, " ")

	parser.Parser(conn, res, s, conf)
}
