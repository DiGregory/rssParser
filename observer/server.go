package observer

import (
	"net"
	"google.golang.org/grpc"
	"github.com/DiGregory/s7testTask/proto"
)

func Start(address string, service *NewsService) error {
	listener, err := net.Listen("tcp", address)
	if err != nil {
		return err
	}
	server := grpc.NewServer()
	proto.RegisterNewsServiceServer(server, service)
	return server.Serve(listener)
}
