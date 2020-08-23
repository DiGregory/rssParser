package observer

import (
	"google.golang.org/grpc"
	"google.golang.org/grpc/grpclog"
	"github.com/DiGregory/s7testTask/proto"
	"context"
	"fmt"
)

func start(host string)(){
	opts := []grpc.DialOption{
		grpc.WithInsecure(),
	}
	conn, err := grpc.Dial(host, opts...)
	if err != nil {
		grpclog.Fatalf("fail to dial: %v", err)
	}

	defer conn.Close()

	client:= proto.NewObserverClient(conn)
	request:=&proto.Request{Message:"lel"}
	response,err:=client.Do(context.Background(),request)
	if err != nil {
		grpclog.Fatalf("fail to dial: %v", err)
	}

	fmt.Println(response.Message)

}