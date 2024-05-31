package main

import (
	"fmt"
	"net"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/dziedzicgrzegorz/Tic-Tac-Toe/constants"
)

func setupConnection(wait bool, ip string, port string) (net.Conn, int, error) {
	if wait {
		ln, err := net.Listen("tcp", ":"+port)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to listen on port %v: %w", port, err)
		}
		conn, err := ln.Accept()
		if err != nil {
			return nil, 0, fmt.Errorf("failed to accept a connection: %w", err)
		}
		return conn, constants.PlayerX, nil
	} else {
		conn, err := net.Dial("tcp", ip+":"+port)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to connect to %v:%v: %w", ip, port, err)
		}
		return conn, constants.PlayerO, nil
	}
}
func createReceiveMove(conn net.Conn) func() tea.Msg {
	return func() tea.Msg {
		buffer := make([]byte, 1024)
		len, err := conn.Read(buffer)
		if err != nil {
			return errMsg{err: err}
		}
		return moveMessage{command: string(buffer[:len])}
	}
}
