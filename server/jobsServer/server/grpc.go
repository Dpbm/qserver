package server

import (
	"errors"
	"os"

	jobsServerProto "github.com/Dpbm/jobsServer/proto"
	logger "github.com/Dpbm/shared/log"
	"google.golang.org/grpc"
	_ "google.golang.org/grpc/encoding/gzip"
	"google.golang.org/grpc/reflection"
)

type GRPC struct {
	TCPServer  *TCPServer
	GRPCServer *grpc.Server
}

func (server *GRPC) Create(host string, port uint16, jobServerDefinition *JobsServer) {
	tcpServer := &TCPServer{}
	tcpServer.Listen(host, port)

	grpcServer := grpc.NewServer()

	if grpcServer == nil {
		logger.LogFatal(errors.New("failed on create new grpc server"))
		os.Exit(1) // just to ensure it will exit
	}

	jobsServerProto.RegisterJobsServer(grpcServer, jobServerDefinition)

	reflection.Register(grpcServer)

	server.TCPServer = tcpServer
	server.GRPCServer = grpcServer
}

func (server *GRPC) Listen() {
	if server.GRPCServer == nil {
		logger.LogFatal(errors.New("no grpc server instance was added"))
		os.Exit(1) // just to ensure it will exit
	}

	err := server.GRPCServer.Serve(server.TCPServer.Listener)

	if err != nil {
		logger.LogFatal(err)
		os.Exit(1) // just to ensure it will exit
	}
}

func (server *GRPC) Close() {
	server.TCPServer.Close()
}
