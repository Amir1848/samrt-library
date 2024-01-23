package main

import (
	"bufio"
	"fmt"
	"net"
	"strings"
)

func main() {
	// Listen for incoming connections on port 8080
	ln, err := net.Listen("tcp", ":8081")
	if err != nil {
		fmt.Println(err)
		return
	}

	// Accept incoming connections and handle them
	for {
		conn, err := ln.Accept()
		if err != nil {
			fmt.Println(err)
			continue
		}

		// Handle the connection in a new goroutine
		go handleConnection(conn)
	}
}

func handleConnection(conn net.Conn) {
	defer func() {
		fmt.Print("closing connection")
		conn.Close()
	}()

	for {
		a, err := bufio.NewReader(conn).ReadString('\n')
		if err != nil {
			fmt.Println(err)
			return
		}

		fmt.Printf("Received: %s", a)
		if strings.EqualFold(a, "end\r\n") {
			break
		}
	}
}
