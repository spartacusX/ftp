package main

import (
	"fmt"
	"github.com/spartacusX/ftp/util"
	"net"
	"time"
)

func main() {
	conn, err := net.Dial("tcp", "localhost:21")
	if err != nil {
		fmt.Println("Failed to connect to BidServer!")
		return
	}

	defer conn.Close()
	fmt.Println("Connected to server: ", conn.RemoteAddr().String())

	cmds, err := util.ReceiveData(conn)
	if err != nil {
		return
	}
	fmt.Println(cmds)

	for {
		// Get user's input
		var cmd string
		_, err = fmt.Scanln(&cmd)
		if err != nil {
			fmt.Println("Scan input failed, error: ", err.Error())
			break
		}

		// Send user's request
		err = util.SendData(conn, cmd)
		if err != nil {
			break
		}

		// Get server's response
		res, err := util.ReceiveData(conn)
		if err != nil {
			return
		}

		fmt.Println(res)
		// Todo		Process server response

		// Close connection since server said "Bye"
		if res == "Bye" {
			return
		}
	}
}
