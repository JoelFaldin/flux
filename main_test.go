package main

import (
	"bufio"
	"fmt"
	"net"
	"strings"
	"sync"
	"testing"
)

func TestMain(t *testing.T) {
	var wg sync.WaitGroup

	for i := range 10 {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()

			conn, err := net.Dial("tcp", "localhost:8080")
			if err != nil {
				t.Errorf("Connection failed: %v\n", err)
				return
			}
			defer conn.Close()

			// Sending SET command:
			fmt.Fprintf(conn, "SET key%d value%d\r\n", id, id)

			// Reading response:
			response, err := bufio.NewReader(conn).ReadString('\n')
			if err != nil {
				t.Errorf("Read failed: %v", err)
			}

			if strings.TrimSpace(response) != "OK" {
				t.Errorf("Expected OK, got %s", response)
			}
		}(i)
	}

	wg.Wait()
}
