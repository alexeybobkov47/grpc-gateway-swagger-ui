package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	pb "test-grpc/api/proto"
	p "test-grpc/internal/parser"

	"google.golang.org/grpc"
)

var port = os.Getenv("GRPC_PORT")

type Service struct {
	pb.UnimplementedGetInfoServer
	p p.ParserInterface
}

func main() {
	srv := grpc.NewServer()

	s := &Service{p: &p.ParserImpl{}}
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
	url := fmt.Sprintf("https://www.rusprofile.ru/search?query=%v", req.Inn)
	info, err := s.p.ParsePage(url)
	if err != nil {
		return nil, err
	}
	return &pb.GetInfoResponse{Inn: info.INN, Kpp: info.KPP, CompanyName: info.CompanyName,
		ChiefName: info.ChiefName}, nil
}
