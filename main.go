package main

import (
	"flag"
	"flux/internal/server"
)

func main() {
	flg := flag.String("port", "", "port for the server")
	flag.Parse()

	server.Server(*flg)
}
