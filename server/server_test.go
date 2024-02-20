// server/user_service_test.go
package main

import (
	"Demo-go/proto"
	"context"
	"fmt"
	"testing"

	"google.golang.org/grpc"
)

const (
	testPort = 50052
)

func TestUserService(t *testing.T) {
	go func() {
		StartGRPCServer(testPort)
	}()

	conn, client, err := createClient(testPort)
	if err != nil {
		t.Fatalf("Failed to create gRPC client: %v", err)
	}
	defer conn.Close()

	t.Run("PostDetails", func(t *testing.T) {
		postDetailsReq := &proto.UserDetailsRequest{Name: "mj", Age: 25}
		resp, err := client.PostDetails(context.Background(), postDetailsReq)
		if err != nil {
			t.Fatalf("PostDetails failed: %v", err)
		}

		expected := &proto.UserDetailsResponse{Name: "mj", Age: 25}
		if resp.GetName() != expected.GetName() || resp.GetAge() != expected.GetAge() {
			t.Errorf("Unexpected response. Expected: %+v, Got: %+v", expected, resp)
		}
	})

	t.Run("GetMyDetails", func(t *testing.T) {
		getDetailsReq := &proto.GetMyDetailsRequest{}
		resp, err := client.GetMyDetails(context.Background(), getDetailsReq)
		if err != nil {
			t.Fatalf("GetMyDetails failed: %v", err)
		}

		expected := &proto.UserDetailsResponse{Name: "mj", Age: 25}
		if resp.GetName() != expected.GetName() || resp.GetAge() != expected.GetAge() {
			t.Errorf("Unexpected response. Expected: %+v, Got: %+v", expected, resp)
		}
	})
}

func createClient(port int) (*grpc.ClientConn, proto.UserServiceClient, error) {
	conn, err := grpc.Dial(fmt.Sprintf("localhost:%d", port), grpc.WithInsecure())
	if err != nil {
		return nil, nil, err
	}
	return conn, proto.NewUserServiceClient(conn), nil
}
