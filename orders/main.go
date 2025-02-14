package main

import (
	"context"
	"log"
	"net"
	"time"

	common "github.com/mrd1920/oms-common"
	"github.com/mrd1920/oms-common/broker"
	"github.com/mrd1920/oms-common/discovery"
	"github.com/mrd1920/oms-common/discovery/consul"
	"google.golang.org/grpc"
)

var (
	grpcAddr     = common.EnvString("GRPC_ADDR", "localhost:2000")
	serviceName  = "orders"
	consulAddr   = common.EnvString("CONSUL_ADDR", "localhost:8500")
	amqpUser     = common.EnvString("RABBITMQ_USER", "user")
	amqpPassword = common.EnvString("RABBITMQ_PASSWORD", "password")
	amqpHost     = common.EnvString("RABBITMQ_HOST", "localhost")
	amqpPort     = common.EnvString("RABBITMQ_PORT", "5672")
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

	//connecting to rabbitmq
	ch, close := broker.Connect(amqpUser, amqpPassword, amqpHost, amqpPort)
	defer func() {
		close()
		ch.Close()
	}()

	//New gRPC server
	grpcServer := grpc.NewServer()
	l, err := net.Listen("tcp", grpcAddr)
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}
	defer l.Close()

	store := NewStore()
	svc := NewService(store)
	NewGrpcHandler(grpcServer, svc, ch)

	// svc.CreateOrder(context.Background())

	log.Println("Starting gRPC server on", grpcAddr)

	if err := grpcServer.Serve(l); err != nil {
		log.Fatal(err.Error())
	}
}
