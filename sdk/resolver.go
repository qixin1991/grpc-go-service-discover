package sdk

import (
	"errors"
	"fmt"
	"strings"

	etcd3 "github.com/coreos/etcd/clientv3"
	"google.golang.org/grpc/naming"
)

// resolver is the implementaion of grpc.naming.Resolver
type resolver struct {
	prefix      string
	serviceName string // service name to resolve
}

// NewResolver return resolver with prefix&service name
func NewResolver(prefix string, serviceName string) *resolver {
	return &resolver{prefix: prefix, serviceName: serviceName}
}

var client etcd3.Client

// var serviceKey string

var stopSignal = make(chan bool, 1)

// Resolve to resolve the service from etcd, target is the dial address of etcd
// target example: "http://127.0.0.1:2379,http://127.0.0.1:12379,http://127.0.0.1:22379"
func (re *resolver) Resolve(target string) (naming.Watcher, error) {
	if re.serviceName == "" {
		return nil, errors.New("grpclb: no service name provided")
	}

	// generate etcd client
	client, err := etcd3.New(etcd3.Config{
		Endpoints: strings.Split(target, ","),
	})
	if err != nil {
		return nil, fmt.Errorf("grpclb: creat etcd3 client failed: %s", err.Error())
	}

	// Return watcher
	return &watcher{re: re, client: *client}, nil
}
