package librarynetwork

import (
	"bufio"
	"context"
	"fmt"
	"net"
	"os"
	"strconv"
	"strings"

	"github.com/Amir1848/samrt-library/services/library"
	"github.com/Amir1848/samrt-library/utils/dbutil"
)

func ServeLibrarySockets() {
	selectedPort := os.Getenv("LibraryPort")

	ln, err := net.Listen("tcp", ":"+selectedPort)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println("listening on port " + selectedPort + " for tcp socket library:")
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
	ctx := context.Background()
	db, err := dbutil.GetDBConnection(ctx)
	if err != nil {
		fmt.Println(err)
		return
	}
	var libraryId int64 = 0
	defer func() {
		if libraryId > 0 {
			library.SetLibraryStatus(ctx, db, libraryId, false)
			library.SetLibraryItemsAsUnknown(ctx, db, libraryId)
		}
		conn.Close()
	}()

	token, err := readStringFromConnection(conn)
	if err != nil {
		fmt.Println(err)
		return
	}

	lib, found, err := library.GetLibraryWithToken(ctx, db, token)
	if err != nil {
		fmt.Println(err)
		return
	}

	if !found {
		return
	}

	libraryId = lib.Id

	err = library.SetLibraryStatus(ctx, db, lib.Id, true)
	if err != nil {
		fmt.Println(err)
		return
	}

	for {
		command, err := readStringFromConnection(conn)
		if err != nil {
			fmt.Println(err)
			return
		}

		if command == "" {
			continue
		}

		commandParts := strings.Split(command, " ")
		if len(commandParts) == 0 {
			continue
		}

		switch commandParts[0] {
		case "set":
			if len(commandParts) != 3 {
				continue
			}

			itemName := commandParts[1]
			status := commandParts[2]
			statusCode, err := strconv.Atoi(status)
			if err != nil {
				conn.Write([]byte("value " + status + " is not valid status"))
				continue
			}

			err = library.SetLibItemStatus(ctx, db, libraryId, itemName, statusCode)
			if err != nil {
				fmt.Print(err)
				return
			}

		}

		fmt.Println(command)
	}
}

func readStringFromConnection(conn net.Conn) (string, error) {
	token, err := bufio.NewReader(conn).ReadString('\n')
	if err != nil {
		return "", err
	}

	token = strings.Trim(token, "\r\n")
	return token, nil
}
