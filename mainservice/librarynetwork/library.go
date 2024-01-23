package librarynetwork

import (
	"bufio"
	"context"
	"fmt"
	"net"
	"os"
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
			library.SetLibraryStatus(ctx, db, libraryId, true)
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

		fmt.Println(command)
	}
}

func readStringFromConnection(conn net.Conn) (string, error) {
	token, err := bufio.NewReader(conn).ReadString('\n')
	if err != nil {
		return "", nil
	}

	token = strings.Trim(token, "\r\n")
	return token, nil
}
