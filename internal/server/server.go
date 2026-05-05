package server

import (
	"flux/internal/handler"
	"flux/internal/store"
	"fmt"
	"net"
	"os"
)

func Server(flg string) {
	var port string
	if flg != "" {
		port = flg
	} else {
		port = "8080"
	}

	p := fmt.Sprintf("127.0.0.1:%s", port)
	ln, err := net.Listen("tcp", p)
	if err != nil {
		fmt.Println("Error starting server", err)
		os.Exit(1)
	}
	defer ln.Close()

	fmt.Printf("Server started on port :%s!\n", port)

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
