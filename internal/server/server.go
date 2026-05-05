package server

import (
	"flux/internal/handler"
	"flux/internal/store"
	"fmt"
	"net"
	"os"
)

func Server() {
	ln, err := net.Listen("tcp", "127.0.0.1:8080")
	if err != nil {
		fmt.Println("Error starting server", err)
		os.Exit(1)
	}
	defer ln.Close()

	fmt.Println("Server started on port :8080!")

	s := store.NewStore()
	s.StartCleaner(10)

	for {
		con, err := ln.Accept()
		if err != nil {
			fmt.Println("Error handling connection", err)
			os.Exit(1)
		}

		// Goroutine :D
		go handler.Handler(con, s)
	}
}
