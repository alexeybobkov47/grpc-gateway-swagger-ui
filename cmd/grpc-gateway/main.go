package main

import (
	"context"
	"flag"
	"log"
	"net/http"
	"os"

	gw "github.com/alexeybobkov47/grpc-gateway-swagger-ui/api/proto"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"google.golang.org/grpc"
)

var (
	// command-line options:
	// gRPC server endpoint
	grpcServer         = os.Getenv("GRPC_HOST") + ":" + os.Getenv("GRPC_PORT")
	grpcServerEndpoint = flag.String("grpc-server-endpoint", grpcServer, "gRPC server endpoint")
)

func run() error {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()
	// Register gRPC server endpoint
	// Note: Make sure the gRPC server is running properly and accessible
	mux := runtime.NewServeMux()
	opts := []grpc.DialOption{grpc.WithInsecure()}
	err := gw.RegisterGetInfoHandlerFromEndpoint(ctx, mux, *grpcServerEndpoint, opts)
	if err != nil {
		return err
	}

	// Start HTTP server (and proxy calls to gRPC server endpoint)
	return http.ListenAndServe(":"+os.Getenv("GRPC_GATEWAY_PORT"), mux)
}

func main() {
	flag.Parse()
	log.Printf("Starting grpc-gateway on localhost:%v", os.Getenv("GRPC_GATEWAY_PORT"))
	if err := run(); err != nil {
		log.Fatalln(err)
	}
}
