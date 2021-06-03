package main

import (
	"context"
	"go-micro/shippy/shippy-service-consignment/proto/consignment"
	"log"
	"net"
	"sync"

	"google.golang.org/grpc"
)

type repository interface {
	Create(*consignment.Consignment) (*consignment.Consignment, error)
}

type Repository struct {
	mu           sync.RWMutex
	consignments []*consignment.Consignment
}

func (repo *Repository) Create(consignmentObj *consignment.Consignment) (*consignment.Consignment, error) {
	repo.mu.Lock()
	defer repo.mu.Unlock()
	updated := append(repo.consignments, consignmentObj)
	repo.consignments = updated
	return consignmentObj, nil
}

type service struct {
	repo repository
	consignment.UnimplementedShippingServiceServer
}

func (s *service) CreateConsignemnt(ctx context.Context, req *consignment.Consignment) (*consignment.Response, error) {
	consignmentObj, err := s.repo.Create(req)
	if err != nil {
		return nil, err
	}

	return &consignment.Response{
		Created:     true,
		Consignment: consignmentObj,
	}, nil
}

func main() {

	repo := &Repository{}
	unimplementedShippingServiceServer := consignment.UnimplementedShippingServiceServer{}

	lis, err := net.Listen("tcp", ":8080")
	if err != nil {
		log.Fatalf("Failed to Listen: %v", err)
	}
	s := grpc.NewServer()

	consignment.RegisterShippingServiceServer(s, &service{repo, unimplementedShippingServiceServer})

	log.Println("Running on port: 8080")
	if err := s.Serve(lis); err != nil {
		log.Fatalf("Failed to Serve: %v", err)
	}
}
