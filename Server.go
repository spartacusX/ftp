package main

import (
	"fmt"
	"github.com/spartacusX/ftp/util"
	"log"
	"net"
	"os"
	"strings"
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
			return
		}

		cmdName, cmdArgs := ParseCmd(cmd)
		switch cmdName {
		case "append":
			util.SendData(conn, "append mode")
		case "ascii":
			util.SendData(conn, "ascii mode")
		case "bell":
			util.SendData(conn, "bell mode")
		case "binary":
			util.SendData(conn, "binary mode")
		case "bye":
			util.SendData(conn, "bye")
			return
		case "cd":
			err = os.Chdir(cmdArgs[0])
			if err != nil {
				log.Println(err.Error())
			}
		case "close":
			util.SendData(conn, "bye")
			return
		case "delete":
			// Delete file
		case "help":
			DisplayCmdList(conn)
		default:
			util.SendData(conn, "Invalid command!")
		}
	}
}

func DisplayCmdList(conn net.Conn) (err error) {
	strCmd := "append ascii bell binary bye close delete help put get mls status lcd rcd lls rls rrmdir lrmdir rrename lrrename rmkdir lmkdir"
	return util.SendData(conn, strCmd)
}

func ParseCmd(cmd string) (Name string, args []string) {
	command := strings.Split(cmd, " ")
	return command[0], command[1:]
}
