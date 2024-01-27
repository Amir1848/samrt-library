package actuator

import (
	"fmt"
	"net"
)

var commandSenderChannel chan string

func SendCommandToActurator(command string) {
	commandSenderChannel <- command
}

func ServeActuatorServer() {
	commandSenderChannel = make(chan string)

	fmt.Println("listening to port 85 for actuator server")
	ln, err := net.Listen("tcp", ":85")
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
		command := <-commandSenderChannel
		fmt.Print(command)
		conn.Write([]byte(command))
	}
}
