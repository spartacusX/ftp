package util

import (
	"bytes"
	"errors"
	"fmt"
	"net"
)

func SendData(conn net.Conn, strData string) (err error) {
	dataLen := len(strData)
	if dataLen <= 0 {
		fmt.Println("Nothing to send, length = 0")
		return errors.New("Invalid data to send.")
	}

	// Send data length
	buf := make([]byte, 1)
	buf[0] = byte(dataLen)
	_, err = conn.Write(buf)
	if err != nil {
		fmt.Println("Failed to send data length, error: ", err.Error())
		return err
	}

	// Send data
	buf = bytes.NewBufferString(strData).Bytes()
	_, err = conn.Write(buf)
	if err != nil {
		fmt.Println("Failed to send data, error: ", err.Error())
		return err
	}

	// Done
	return nil
}

func ReceiveData(conn net.Conn) (strData string, err error) {
	// Receive data length
	buf := make([]byte, 1)
	_, err = conn.Read(buf)
	if err != nil {
		fmt.Println("Failed to receive data length, error: ", err.Error())
		return strData, err
	}

	dataLen := int(buf[0])
	if dataLen <= 0 {
		fmt.Println("Nothing to receive, length = 0")
		return strData, errors.New("Invalid data length.")
	}

	// Receive data
	buf = make([]byte, dataLen)
	_, err = conn.Read(buf)
	if err != nil {
		fmt.Println("Failed to receive data, error: ", err.Error())
		return strData, err
	}

	// Done
	return bytes.NewBuffer(buf).String(), nil
}
