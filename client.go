package main

import (
	"flag"
	"fmt"
	"time"

	"strconv"

	pb "grpc-client/rpc"

	grpclb "grpc-client/etcdv3"

	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

var (
	serv = flag.String("service", "hello_service", "service name")
	reg  = flag.String("reg", "http://172.20.9.101:2379,http://172.20.9.103:2379,http://172.20.9.105:2379", "register etcd address")
)

func main() {
	flag.Parse()
	r := grpclb.NewResolver(*serv)
	b := grpc.RoundRobin(r)

	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	conn, err := grpc.DialContext(ctx, *reg, grpc.WithInsecure(), grpc.WithBalancer(b))
	if err != nil {
		panic(err)
	}

	ticker := time.NewTicker(1 * time.Second)
	for t := range ticker.C {
		client := pb.NewGreeterClient(conn)
		resp, err := client.SayHello(context.Background(), &pb.HelloRequest{Name: "world " + strconv.Itoa(t.Second())})
		if err == nil {
			fmt.Printf("%v: Reply is %s\n", t, resp.Message)
		}
		resp2, err2 := client.SayHelloAgain(context.Background(), &pb.HelloRequest{Name: "world " + strconv.Itoa(t.Second())})
		if err2 == nil {
			fmt.Printf("%v: Reply Again is %s\n", t, resp2.Message)
		} else {
			fmt.Println(err2)
		}
	}
}
