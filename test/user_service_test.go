package test

import (
	"context"
	"fmt"
	"log"
	"net"
	"testing"
	"time"

	pb "Demo-go/proto"

	"google.golang.org/grpc"
)

func TestUserService(t *testing.T) {
	port := 50052
	serverAddr := fmt.Sprintf("localhost:%d", port)

	go func() {
		mainWithPort(port)
	}()

	time.Sleep(100 * time.Millisecond)

	conn, err := grpc.Dial(serverAddr, grpc.WithInsecure())
	if err != nil {
		t.Fatalf("Failed to connect to server: %v", err)
	}
	defer conn.Close()

	client := pb.NewUserServiceClient(conn)

	postDetailsReq := &pb.UserDetailsRequest{Name: "John", Age: 30}
	_, err = client.PostDetails(context.Background(), postDetailsReq)
	if err != nil {
		t.Fatalf("PostDetails failed: %v", err)
	}

	getDetailsReq := &pb.GetMyDetailsRequest{}
	details, err := client.GetMyDetails(context.Background(), getDetailsReq)
	if err != nil {
		t.Fatalf("GetMyDetails failed: %v", err)
	}

	expectedDetails := &pb.UserDetailsResponse{Name: "AuthenticatedUser", Age: 25}
	if details.GetName() != expectedDetails.GetName() || details.GetAge() != expectedDetails.GetAge() {
		t.Errorf("Unexpected response. Expected: %+v, Got: %+v", expectedDetails, details)
	}
}

// for test only
func mainWithPort(port int) {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()

	pb.RegisterGreeterServer(s, &server{})

	pb.RegisterUserServiceServer(s, &userService{})

	log.Printf("server listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
