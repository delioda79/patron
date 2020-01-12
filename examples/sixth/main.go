package main

import (
	"context"
	"fmt"
	"os"

	"github.com/beatlabs/patron"
	"github.com/beatlabs/patron/grpc"
	"github.com/beatlabs/patron/log"
)

const (
	awsRegion      = "eu-west-1"
	awsID          = "test"
	awsSecret      = "test"
	awsToken       = "token"
	awsSQSEndpoint = "http://localhost:4576"
	awsSQSQueue    = "patron"
)

func init() {
	err := os.Setenv("PATRON_LOG_LEVEL", "debug")
	if err != nil {
		fmt.Printf("failed to set log level env var: %v", err)
		os.Exit(1)
	}
	err = os.Setenv("PATRON_JAEGER_SAMPLER_PARAM", "1.0")
	if err != nil {
		fmt.Printf("failed to set sampler env vars: %v", err)
		os.Exit(1)
	}
	err = os.Setenv("PATRON_HTTP_DEFAULT_PORT", "50005")
	if err != nil {
		fmt.Printf("failed to set default patron port env vars: %v", err)
		os.Exit(1)
	}
}

type greeterServer struct {
	UnimplementedGreeterServer
}

func (gs *greeterServer) SayHello(ctx context.Context, req *HelloRequest) (*HelloReply, error) {
	return &HelloReply{Message: fmt.Sprintf("Hello, %s %s!", req.GetFirstname(), req.GetLastname())}, nil
}

func main() {
	name := "sixth"
	version := "1.0.0"

	err := patron.Setup(name, version)
	if err != nil {
		fmt.Printf("failed to set up logging: %v", err)
		os.Exit(1)
	}

	cmp, err := grpc.New(":50006").Create()
	if err != nil {
		log.Fatalf("failed to create gRPC component: %v", err)
	}

	RegisterGreeterServer(cmp.Server(), &greeterServer{})

	srv, err := patron.New(name, version, patron.Components(cmp))
	if err != nil {
		log.Fatalf("failed to create service: %v", err)
	}

	ctx := context.Background()
	err = srv.Run(ctx)
	if err != nil {
		log.Fatalf("failed to run service: %v", err)
	}
}
