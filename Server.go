package main

import (
	"fmt"
	"github.com/spartacusX/ftp/util"
	"net"
	"strconv"
)

var clients map[net.Conn]int

func main() {
	ln, err := net.Listen("tcp", ":21")
	if err != nil {
		fmt.Println("Listen on port 21 failed, error: ", err.Error())
		return
	}

	clients = make(map[net.Conn]int)

	for {
		fmt.Println("Waiting for new request...")
		conn, err := ln.Accept()
		if err != nil {
			fmt.Println("Accept connection reqest failed!")
			continue
		}

		go ConnHandler(conn)
	}
}

func RecordClient(conn net.Conn) {
	clients[conn] = len(clients)
}

func DeleteClient(conn net.Conn) {
	delete(clients, conn)
}

func ConnHandler(conn net.Conn) {
	RecordClient(conn)
	fmt.Println("Welcome ", conn.RemoteAddr().String())
	defer conn.Close()
	defer DeleteClient(conn)
	defer fmt.Println("Connection closed from: ", conn.RemoteAddr().String())

	err := DisplayCmdList(conn)
	if err != nil {
		return
	}

	for {
		cmd, err := util.ReceiveData(conn)
		if err != nil {
			continue
		}

		switch cmd {
		case "1":
			util.SendData(conn, "Sign in successfully!")
		case "2":
			util.SendData(conn, "Sign up successfully!")
		case "3":
			util.SendData(conn, strconv.Itoa(len(clients)))
		case "4":
			util.SendData(conn, "Sign out successfully!")
		case "5":
			util.SendData(conn, "Bye")
			return
		default:
			util.SendData(conn, "Invalid command!")
		}
	}
}

func DisplayCmdList(conn net.Conn) (err error) {
	strCmd := "1.Sign In\n2.Sign Up\n3.ClientCount\n4.Sign Out\n5.Bye"
	return util.SendData(conn, strCmd)
}
