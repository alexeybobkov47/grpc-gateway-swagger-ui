package main

import (
	"context"
	"log"
	"net"
	"os"

	pb "github.com/alexeybobkov47/grpc-gateway-swagger-ui/api/proto"
	p "github.com/alexeybobkov47/grpc-gateway-swagger-ui/internal/parser"

	"google.golang.org/grpc"
)

var port = os.Getenv("GRPC_PORT")

type Service struct {
	pb.UnimplementedGetInfoServer
	p p.ParseInterface
}

func main() {
	srv := grpc.NewServer()

	s := &Service{p: &p.ParseImpl{}}
	pb.RegisterGetInfoServer(srv, s)

	listener, err := net.Listen("tcp", ":"+port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	log.Printf("Starting server on %v", listener.Addr())
	if err := srv.Serve(listener); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

func (s *Service) GetInfoByINN(ctx context.Context, req *pb.GetInfoRequest) (*pb.GetInfoResponse, error) {
	info, err := s.p.ParsePage(req.Inn)
	if err != nil {
		return nil, err
	}
	return &pb.GetInfoResponse{
		Inn: info.INN, Kpp: info.KPP, CompanyName: info.CompanyName,
		ChiefName: info.ChiefName,
	}, nil
}
