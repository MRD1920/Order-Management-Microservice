package main

import (
	"context"
	"log"
	"net"
	"time"

	common "github.com/mrd1920/oms-common"
	"github.com/mrd1920/oms-common/discovery"
	"github.com/mrd1920/oms-common/discovery/consul"
	"google.golang.org/grpc"
)

var (
	grpcAddr    = common.EnvString("GRPC_ADDR", "localhost:2000")
	serviceName = "orders"
	consulAddr  = common.EnvString("CONSUL_ADDR", "localhost:8500")
)

func main() {
	registry, err := consul.NewRegistry(consulAddr, serviceName)
	if err != nil {
		panic(err)
	}

	ctx := context.Background()
	instanceID := discovery.GenerateInstanceID(serviceName)
	if err := registry.Register(ctx, instanceID, serviceName, grpcAddr); err != nil {
		panic(err)
	}

	go func() {
		for {
			if err := registry.HealthCheck(instanceID, serviceName); err != nil {
				log.Printf("Health check failed: %v", err)
			}
			time.Sleep(time.Second * 1)
		}
	}()

	defer registry.Deregister(ctx, instanceID, serviceName)

	grpcServer := grpc.NewServer()
	l, err := net.Listen("tcp", grpcAddr)
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}
	defer l.Close()

	store := NewStore()
	svc := NewService(store)
	NewGrpcHandler(grpcServer, svc)

	svc.CreateOrder(context.Background())

	log.Println("Starting gRPC server on", grpcAddr)

	if err := grpcServer.Serve(l); err != nil {
		log.Fatal(err.Error())
	}
}
