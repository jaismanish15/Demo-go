// server/main.go
package main

import (
	"Demo-go/config"
	"Demo-go/proto"
	"context"
	"database/sql"
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var (
	port = flag.Int("port", 50051, "The server port")
	db   *sql.DB
)

// UserService struct
type UserService struct {
	proto.UnimplementedUserServiceServer
	mu    sync.Mutex
	users map[string]*proto.UserDetailsResponse
}

// initializeDatabase function
func initializeDatabase() error {
	_, err := db.Exec(`
	CREATE TABLE IF NOT EXISTS user_details (
		name VARCHAR(255) PRIMARY KEY,
		age INT
	)
	`)
	if err != nil {
		return fmt.Errorf("failed to create user_details table: %v", err)
	}
	return nil
}

// PostDetails function
func (s *UserService) PostDetails(ctx context.Context, req *proto.UserDetailsRequest) (*proto.UserDetailsResponse, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.users == nil {
		s.users = make(map[string]*proto.UserDetailsResponse)
	}
	s.users[req.GetName()] = &proto.UserDetailsResponse{Name: req.GetName(), Age: req.GetAge()}
	log.Printf("Received user details for %s with age %d", req.GetName(), req.GetAge())

	return &proto.UserDetailsResponse{Name: req.GetName(), Age: req.GetAge()}, nil
}

// GetMyDetails function
func (s *UserService) GetMyDetails(ctx context.Context, req *proto.GetMyDetailsRequest) (*proto.UserDetailsResponse, error) {
	userName := "mj"

	s.mu.Lock()
	defer s.mu.Unlock()

	if userDetails, ok := s.users[userName]; ok {
		log.Printf("Retrieved details for user: %s, Age: %d", userName, userDetails.GetAge())
		return userDetails, nil
	}

	return nil, status.Errorf(codes.NotFound, "User not found")
}

// StartGRPCServer function
func StartGRPCServer(port int) {
	flag.Parse()
	userService := &UserService{users: make(map[string]*proto.UserDetailsResponse)}
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	proto.RegisterUserServiceServer(s, userService)
	log.Printf("server listening at %v", lis.Addr())

	// Use a WaitGroup to wait for the server to finish serving
	var wg sync.WaitGroup
	wg.Add(1)

	// Run gRPC server in a goroutine
	go func() {
		defer wg.Done()
		if err := s.Serve(lis); err != nil {
			log.Fatalf("failed to serve: %v", err)
		}
	}()

	// Set up signal handling for graceful shutdown
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	// Wait for a signal to stop the server
	<-stop
	log.Println("Shutting down server...")
	s.GracefulStop()
	log.Println("Server stopped gracefully")

	// Wait for the server goroutine to finish
	wg.Wait()
}

// startPostgres function
func startPostgres() {
	flag.Parse()

	// Get PostgreSQL connection string from config
	connectionString, err := config.GetConnectionString()
	if err != nil {
		log.Fatalf("failed to get PostgreSQL connection string: %v", err)
	}

	// Initialize PostgreSQL connection
	db, err = sql.Open("postgres", connectionString)
	if err != nil {
		log.Fatalf("failed to connect to PostgreSQL: %v", err)
	}
	defer db.Close()

	// Initialize PostgreSQL database
	if err := initializeDatabase(); err != nil {
		log.Fatalf("failed to initialize database: %v", err)
	}
}

func main() {
	StartGRPCServer(*port)
	go startPostgres()
}
