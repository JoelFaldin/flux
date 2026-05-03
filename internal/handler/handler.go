package handler

import "net"

// Each new connection spawns a new gorouting that executes this function.
func Handler(con net.Conn) {
	defer con.Close()

}
