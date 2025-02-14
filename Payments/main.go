package main

import (
	"context"
	"log"
	"net"

	common "github.com/mrd1920/oms-common"
	"github.com/mrd1920/oms-common/broker"
	"github.com/mrd1920/oms-common/discovery"
	"github.com/mrd1920/oms-common/discovery/consul"
	stripeProcessor "github.com/mrd1920/oms-payments/processor/stripe"

	"github.com/stripe/stripe-go/v78"
	"google.golang.org/grpc"
)

var (
	serviceName = "payments"
	amqpUser    = common.EnvString("AMQP_USER", "user")
	amqpPass    = common.EnvString("AMQP_PASS", "password")
	amqpPort    = common.EnvString("AMQP_PORT", "5672")
	amqpHost    = common.EnvString("AMQP_HOST", "localhost")
	grpcAddr    = common.EnvString("GRPC_ADDR", "localhost:2001")
	consulAddr  = common.EnvString("CONSUL_ADDR", "localhost:8500")
	stripeKey   = common.EnvString("STRIPE_KEY", "")
	httpAddr    = common.EnvString("HTTP_ADDR", "localhost:8081")
)

func main() {

	//Register with consul
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
				log.Fatalf("Health check failed: %v", err.Error())
			}
		}
	}()

	defer registry.Deregister(ctx, instanceID, serviceName)

	//Stripe setup
	stripe.Key = stripeKey

	//Message broker connection
	channel, close := broker.Connect(amqpUser, amqpPass, amqpHost, amqpPort)
	defer func() {
		close()
		channel.Close()
	}()

	stripeprocessor := stripeProcessor.NewProcessor()
	svc := NewService(stripeprocessor)
	amqpConsumer := NewConsumer(svc)
	amqpConsumer.Listen(channel)

	//gRPC server
	grpcServer := grpc.NewServer()

	l, err := net.Listen("tcp", grpcAddr)
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}
	defer l.Close()

	log.Printf("Starting gRPC server on %s", grpcAddr)
	if err := grpcServer.Serve(l); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}

}
