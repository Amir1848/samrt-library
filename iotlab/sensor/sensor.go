package sensor

import (
	"bufio"
	"fmt"
	"net"
	"strings"
	"test/actuator"
)

func ServeSensorServer() {
	fmt.Println("listening to port 86 for sensor server")
	ln, err := net.Listen("tcp", ":86")
	if err != nil {
		fmt.Println(err)
		return
	}

	for {
		conn, err := ln.Accept()
		if err != nil {
			fmt.Println(err)
			continue
		}

		handleConnection(conn)
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

		actuator.SendCommandToActurator(a)
		if strings.EqualFold(a, "end\r\n") {
			break
		}
	}
}
