package main

import (
	"log"
	"net"

	"github.com/Farabi99/Simple-gRPC-Exploration/controllers"
	"github.com/Farabi99/Simple-gRPC-Exploration/models"
	pb "github.com/Farabi99/Simple-gRPC-Exploration/proto"
	"google.golang.org/grpc"
)

func main() {
	// 1. Initialize Model
	repo := models.NewInMemoryUserRepo()

	log.Println("Seeding database with 100 users...")
	repo.Seed(100)

	// 2. Initialize Controller
	userController := controllers.NewUserController(repo)

	// 3. Setup gRPC Server
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()

	// Register the controller with the gRPC server
	pb.RegisterUserServiceServer(grpcServer, userController)

	log.Printf("gRPC server listening at %v", lis.Addr())
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
