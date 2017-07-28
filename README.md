# grpc-go-service-discover

gRPC golang service discover sdk.

# Usage Example

_client.go_

```golang
package main

import (
	"flag"
	"fmt"
	"time"

	"strconv"

	"github.com/qixin1991/grpc-go-service-discover/rpc"

	"github.com/qixin1991/grpc-go-service-discover/sdk"

	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

var (
	serv = flag.String("service", "hello_service", "service name")
	reg  = flag.String("reg", "http://172.20.9.101:2379,http://172.20.9.103:2379,http://172.20.9.105:2379", "register etcd address")
)

func main() {
	flag.Parse()
	r := sdk.NewResolver(*serv)
	b := grpc.RoundRobin(r)

	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	conn, err := grpc.DialContext(ctx, *reg, grpc.WithInsecure(), grpc.WithBalancer(b))
	if err != nil {
		panic(err)
	}

	ticker := time.NewTicker(1 * time.Second)
	for t := range ticker.C {
		client := rpc.NewGreeterClient(conn)
		resp, err := client.SayHello(context.Background(), &rpc.HelloRequest{Name: "world " + strconv.Itoa(t.Second())})
		if err == nil {
			fmt.Printf("%v: Reply is %s\n", t, resp.Message)
		}
		resp2, err2 := client.SayHelloAgain(context.Background(), &rpc.HelloRequest{Name: "world " + strconv.Itoa(t.Second())})
		if err2 == nil {
			fmt.Printf("%v: Reply Again is %s\n", t, resp2.Message)
		} else {
			fmt.Println(err2)
		}
	}
}

```

Issue the following command

```shell
go run client.go --reg [etcd3_cluser_string]
```

> NOTE: Start your server instance before running the client.
>
> For grpc service registry, see https://github.com/qixin1991/grpc-go-register.git