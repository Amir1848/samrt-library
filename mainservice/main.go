package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"

	"github.com/Amir1848/samrt-library/routes"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func init() {
	var envFileName string
	env := os.Getenv("ENV")
	if env == "development" {
		envFileName = "development.env"
	} else if env == "production" {
		envFileName = "production.env"
	} else {
		//todo: enable panic
		// panic("enviroment variable has not been set")
		envFileName = "development.env"
	}

	fmt.Println(envFileName)
	err := godotenv.Load(envFileName)
	if err != nil {
		panic("cannot load env file")
	}
}

func main() {
	r := gin.Default()

	routes.AddMainRoutes(r)

	go serveLibrarySockets()
	r.Run()
}

func serveLibrarySockets() {
	ln, err := net.Listen("tcp", ":8081")
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
