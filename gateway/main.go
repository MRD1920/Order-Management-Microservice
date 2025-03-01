package main

import (
	"context"
	"log"
	"net/http"
	"time"

	_ "github.com/joho/godotenv/autoload"
	common "github.com/mrd1920/oms-common"
	"github.com/mrd1920/oms-common/discovery"
	"github.com/mrd1920/oms-common/discovery/consul"
	"github.com/mrd1920/oms-gateway/gateway"
)

// Flow Example
// For a create order request:

// 1. Client → Makes HTTP POST request to gateway
// 2. HTTP Handler → Validates request and converts to internal format
// 3. Gateway →
// 		a. Uses Consul to find available order service instance
// 		b. Establishes gRPC connection
// 		c. Forwards request to service
// 4. Order Service → Processes request and returns response
// 5. Gateway → Converts gRPC response to HTTP
// 6. Client → Receives HTTP response
// ---------------------------------------------
// Why Pass Registry to Gateway?
// The gateway needs the registry to:

// 1. Discover service instances dynamically
// 2. Load balance between multiple instances
// 3. Handle service failover
// 4. Maintain up-to-date service endpoints

var (
	serviceName = "gateway"
	httpAddr    = common.EnvString("HTTP_ADDR", ":3000")
	consulAddr  = common.EnvString("CONSUL_ADDR", "localhost:8500")
)

func main() {
	registry, err := consul.NewRegistry(consulAddr, serviceName)
	if err != nil {
		panic(err)
	}

	ctx := context.Background()
	instanceID := discovery.GenerateInstanceID(serviceName)
	if err := registry.Register(ctx, instanceID, serviceName, httpAddr); err != nil {
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

	mux := http.NewServeMux()

	ordersGateway := gateway.NewGRPCGateway(registry)

	handler := NewHandler(ordersGateway)
	handler.registerRoutes(mux)

	log.Printf("Starting HTTP server at %s", httpAddr)

	if err := http.ListenAndServe(httpAddr, mux); err != nil {
		log.Fatal("Failed to create a http server: ", err)
	}

}
