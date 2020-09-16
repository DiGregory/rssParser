package observer

import (
	"net"

	"github.com/DiGregory/rssParser/proto"
	"google.golang.org/grpc"
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
