package server

import (
	jobsServerProto "github.com/Dpbm/jobsServer/proto"
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

	jobsServerProto.RegisterJobsServer(grpcServer, jobServerDefinition)

	reflection.Register(grpcServer)

	server.TCPServer = tcpServer
	server.GRPCServer = grpcServer
}

func (server *GRPC) Listen() {
	server.GRPCServer.Serve(server.TCPServer.Listener)
}

func (server *GRPC) Close() {
	server.TCPServer.Close()
}
