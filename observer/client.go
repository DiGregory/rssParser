package observer

import (
	"google.golang.org/grpc"
	"google.golang.org/grpc/grpclog"
	"github.com/DiGregory/s7testTask/proto"
	"fmt"
	"context"
)

func ClientStart(host string) () {
	opts := []grpc.DialOption{
		grpc.WithInsecure(),
	}
	conn, err := grpc.Dial(host, opts...)
	if err != nil {
		grpclog.Fatalf("fail to dial: %v", err)
	}

	defer conn.Close()

	client := proto.NewNewsServiceClient(conn)
	request := &proto.GetNewsRequest{}
	response, err := client.GetNews(context.Background(), request)
	if err != nil {
		grpclog.Fatalf("fail to dial: %v", err)
	}

	fmt.Println(response.News)
}
