package rpc

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/leexsh/go-todolist/config"
	pb_todolist "github.com/leexsh/go-todolist/idl/pb"
	"github.com/leexsh/go-todolist/pkg/discovery"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/resolver"
)

var (
	UserClient pb_todolist.UserServiceClient
	TaskClient pb_todolist.TaskServiceClient

	CancleFunc context.CancelFunc
	ctx        context.Context
)

func Init() {
	Register := discovery.NewResolver([]string{config.Conf.Etcd.Address})
	resolver.Register(Register)
	ctx, _ = context.WithTimeout(context.Background(), 3*time.Second)

	defer Register.Close()
	initClient(config.Conf.Domains["user"].Name, &UserClient)
	initClient(config.Conf.Domains["task"].Name, &TaskClient)
}
func initClient(serviceName string, client interface{}) {
	conn, err := connectServer(serviceName)

	if err != nil {
		panic(err)
	}

	switch c := client.(type) {
	case *pb_todolist.UserServiceClient:
		*c = pb_todolist.NewUserServiceClient(conn)
	case *pb_todolist.TaskServiceClient:
		*c = pb_todolist.NewTaskServiceClient(conn)
	default:
		panic("unsupported client type")
	}
}

func connectServer(serviceName string) (conn *grpc.ClientConn, err error) {
	opts := []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	}
	addr := fmt.Sprintf("%s:///%s", "etcd", serviceName)

	// Load balance
	if config.Conf.Services[serviceName].LoadBalance {
		log.Printf("load balance enabled for %s\n", serviceName)
		opts = append(opts, grpc.WithDefaultServiceConfig(fmt.Sprintf(`{"LoadBalancingPolicy": "%s"}`, "round_robin")))
	}

	conn, err = grpc.DialContext(ctx, addr, opts...)
	return
}
