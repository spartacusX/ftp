package main

import (
	"bufio"
	"errors"
	"fmt"
	"github.com/spartacusX/ftp/util"
	"net"
	"os"
	"path/filepath"
	"strings"
)

const (
	CMD_UNKNOWN = iota
	CMD_LOCAL
	CMD_REMOTE
)

var commands []string // Command list
var strCurrentDir string

func main() {
	conn, err := net.Dial("tcp", "localhost:21")
	if err != nil {
		fmt.Println("Failed to connect to server!")
		return
	}

	defer conn.Close()

	strCurrentDir, err = filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		fmt.Println("Get current directory failed, error: ", err.Error())
		return
	}

	fmt.Println("Connected to server: ", conn.RemoteAddr().String())

	cmdList, err := util.ReceiveData(conn)
	if err != nil {
		return
	}

	CacheCmdList(cmdList)
	ShowCmdList()

	for {
		// Get user's input
		fmt.Printf("> ")
		scanner := bufio.NewScanner(os.Stdin)
		if scanner.Scan() {
			cmd := scanner.Text()

			if ProcessCmd(conn, cmd) == true { // exit
				return
			}
		}
	}
}

func ProcessCmd(conn net.Conn, cmd string) (bExit bool) {
	bExit = false
	if cmdType, cmdName, cmdArgs := ParseCmd(cmd); cmdType != CMD_UNKNOWN {
		if cmdType == CMD_LOCAL {
			ProcessLocalCmd(cmdName, cmdArgs)
		} else if cmdType == CMD_REMOTE {
			if err := ForwordRemoteCmd(conn, cmd); err != nil {
				bExit = true
				fmt.Println(err.Error())
			}
		} else {
			fmt.Println("Invalid command!")
		}

	}

	return bExit
}

func CacheCmdList(cmdList string) {
	cmds := strings.Split(cmdList, " ")
	for _, cmd := range cmds {
		commands = append(commands, strings.TrimSpace(cmd))
	}
}

func ShowCmdList() {
	for i, cmd := range commands {
		fmt.Printf("%-15s", cmd)
		if (i+1)%9 == 0 {
			fmt.Printf("\n")
		}
	}
	fmt.Printf("\n")
}

func ParseCmd(strCmd string) (cmdType int, cmdName string, cmdArgs []string) {
	strCmdName := strings.Split(strCmd, " ")[0]
	slArgs := strings.Split(strCmd, " ")[1:]
	cmdType = CMD_UNKNOWN

	for _, cmd := range commands {
		if strCmdName == cmd {
			cmdName = strCmdName
			cmdArgs = slArgs
			if strCmdName == "help" || strCmdName == "lcd" {
				cmdType = CMD_LOCAL
			} else {
				cmdType = CMD_REMOTE
			}
			break
		}
	}
	return
}

func ProcessLocalCmd(cmdName string, cmdArgs []string) (err error) {
	switch cmdName {
	case "help":
		ShowCmdList()
	case "lcd":
		err = os.Chdir(cmdArgs[0])
		if err == nil {
			strCurrentDir = cmdArgs[0]
		}
	case "ldir":
		fmt.Println(strCurrentDir)
	default:
		fmt.Println("Unknown local command")
	}
	return
}

func ForwordRemoteCmd(conn net.Conn, cmd string) (err error) {
	// Send request
	if err = util.SendData(conn, cmd); err != nil {
		return
	}

	// Get response
	res, err := util.ReceiveData(conn)
	if err != nil {
		return
	}

	fmt.Println("Server Response:", res)
	// Todo		Process server response

	// Close connection since server said "Bye"
	if res == "bye" {
		err = errors.New("Disconnected by server.")
	}

	return
}
