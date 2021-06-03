package main

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"log"

	"go-micro/shippy/shippy-service-consignment/proto/consignment"

	"os"

	"google.golang.org/grpc"
)

const (
	address         = "localhost:50051"
	defaultFilename = "consignment.json"
)

func parseFile(file string) (*consignment.Consignment, error) {
	consignmentObj := consignment.Consignment{}
	data, err := ioutil.ReadFile(file)
	if err != nil {
		log.Printf("Failed to Read: %v", err)
		return &consignmentObj, err
	}
	err = json.Unmarshal(data, &consignmentObj)
	if err != nil {
		log.Printf("Failed to Unmarshal: %v", err)
		return &consignmentObj, err
	}

	return &consignmentObj, nil
}
func main() {

	conn, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Failed to Connect: %v", err)
	}
	defer conn.Close()
	client := consignment.NewShippingServiceClient(conn)

	file := defaultFilename

	if len(os.Args) > 2 {
		file = os.Args[1]
	}
	consignment, err := parseFile(file)
	if err != nil {
		log.Fatalf("Failed to Parse: %v", err)
	}

	res, err := client.CreateConsignemnt(context.Background(), consignment)
	if err != nil {
		log.Fatalf("Failed to Create: %v", err)
	}
	log.Printf("Created: %v", res.Created)
}
