package server

import (
	"fmt"
	"net"

	logger "github.com/Dpbm/shared/log"
)

type TCPServer struct {
	ServerURL string
	Listener  net.Listener
}

func (server *TCPServer) Listen(host string, port uint32) {
	serverURL := fmt.Sprintf("%s:%d", host, port)
	listen, err := net.Listen("tcp", serverURL)
	if err != nil {
		logger.LogFatal(err) // it will exit with status 1
	}

	server.ServerURL = serverURL
	server.Listener = listen
}

func (server *TCPServer) Close() {
	server.Listener.Close()
}
