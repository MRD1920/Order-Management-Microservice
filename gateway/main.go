package main

import (
	"log"
	"net/http"

	_ "github.com/joho/godotenv/autoload"
	common "github.com/mrd1920/oms-common"
	pb "github.com/mrd1920/oms-common/api"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var (
	httpAddr          = common.EnvString("HTTP_ADDR", ":3000")
	ordersServiceAddr = "localhost:3000"
)

func main() {
	conn, err := grpc.NewClient(ordersServiceAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Failed to dial server: %v", err)
	}
	defer conn.Close()

	log.Println("Dialing orders service at ", ordersServiceAddr)

	c := pb.NewOrderServiceClient(conn)

	mux := http.NewServeMux()
	handler := NewHandler(c)
	handler.registerRoutes(mux)

	log.Printf("Starting HTTP server at %s", httpAddr)

	if err := http.ListenAndServe(httpAddr, mux); err != nil {
		log.Fatal("Faild to create a http server: ", err)
	}

}
