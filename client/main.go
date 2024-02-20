package main

import (
	"context"
	"flag"
	"fmt"
	"log"

	pb "Demo-go/proto"
	"google.golang.org/grpc"
)

var (
	addr = flag.String("addr", "localhost:50051", "the address to connect to")
	name = flag.String("name", "mj", "Name to greet") // Updated the default name
)

func main() {
	flag.Parse()

	// Set up a connection to the server
	conn, err := grpc.Dial(*addr, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()

	// Create a client instance
	client := pb.NewUserServiceClient(conn)

	// Call PostDetails
	postDetailsReq := &pb.UserDetailsRequest{Name: *name, Age: 25}
	postDetailsResp, err := client.PostDetails(context.Background(), postDetailsReq)
	if err != nil {
		log.Fatalf("PostDetails failed: %v", err)
	}
	fmt.Printf("PostDetails Response: %+v\n", postDetailsResp)

	// Call GetMyDetails
	getDetailsReq := &pb.GetMyDetailsRequest{}
	getDetailsResp, err := client.GetMyDetails(context.Background(), getDetailsReq)
	if err != nil {
		log.Fatalf("GetMyDetails failed: %v", err)
	}
	fmt.Printf("GetMyDetails Response: %+v\n", getDetailsResp)
}
