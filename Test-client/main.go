package main

import (
	"context"
	"fmt"
	"log"
	"time"

	pb "github.com/Farabi99/Simple-gRPC-Exploration/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
)

const (
	ServerAddress = "localhost:50051"
)

var client pb.UserServiceClient

func main() {
	conn, err := grpc.NewClient(ServerAddress, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()

	client = pb.NewUserServiceClient(conn)

	fmt.Println("======================================")
	fmt.Println("    gRPC User Service Test Client     ")
	fmt.Println("======================================\n")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	testCreateUser(ctx)
	testGetUser(ctx)
	testListUsers(ctx)
	testUpdateUser(ctx)
	testGetUserAfterUpdate(ctx)
	testDeleteUser(ctx)
	testGetUserNotFound(ctx)
	testListUsersWithPagination(ctx)

	fmt.Println("\n======================================")
	fmt.Println("       All Tests Completed!          ")
	fmt.Println("======================================")
}

func testCreateUser(ctx context.Context) {
	fmt.Println("\n--- Test: Create User ---")
	res, err := client.CreateUser(ctx, &pb.CreateUserRequest{
		Name:  "John Doe",
		Email: "john.doe@example.com",
	})
	if err != nil {
		log.Printf("CreateUser failed: %v", err)
		return
	}
	fmt.Printf("Created user:\n  ID:    %s\n  Name:  %s\n  Email: %s\n", res.User.Id, res.User.Name, res.User.Email)
}

func testGetUser(ctx context.Context) {
	fmt.Println("\n--- Test: Get User ---")
	// First create a user to get
	createRes, err := client.CreateUser(ctx, &pb.CreateUserRequest{
		Name:  "Jane Smith",
		Email: "jane.smith@example.com",
	})
	if err != nil {
		log.Printf("CreateUser failed: %v", err)
		return
	}

	res, err := client.GetUser(ctx, &pb.GetUserRequest{
		Id: createRes.User.Id,
	})
	if err != nil {
		log.Printf("GetUser failed: %v", err)
		return
	}
	fmt.Printf("Retrieved user:\n  ID:    %s\n  Name:  %s\n  Email: %s\n", res.User.Id, res.User.Name, res.User.Email)
}

func testListUsers(ctx context.Context) {
	fmt.Println("\n--- Test: List Users ---")
	res, err := client.ListUsers(ctx, &pb.ListUsersRequest{
		Limit: 10,
	})
	if err != nil {
		log.Printf("ListUsers failed: %v", err)
		return
	}
	fmt.Printf("Retrieved %d users:\n", len(res.Users))
	for i, u := range res.Users {
		fmt.Printf("  %d. [%s] %s (%s)\n", i+1, u.Id[:8], u.Name, u.Email)
	}
}

func testUpdateUser(ctx context.Context) {
	fmt.Println("\n--- Test: Update User ---")
	// First create a user to update
	createRes, err := client.CreateUser(ctx, &pb.CreateUserRequest{
		Name:  "Old Name",
		Email: "old.email@example.com",
	})
	if err != nil {
		log.Printf("CreateUser failed: %v", err)
		return
	}

	res, err := client.UpdateUser(ctx, &pb.UpdateUserRequest{
		Id:    createRes.User.Id,
		Name:  "New Name",
		Email: "new.email@example.com",
	})
	if err != nil {
		log.Printf("UpdateUser failed: %v", err)
		return
	}
	fmt.Printf("Updated user:\n  ID:    %s\n  Name:  %s\n  Email: %s\n", res.User.Id, res.User.Name, res.User.Email)
}

func testGetUserAfterUpdate(ctx context.Context) {
	fmt.Println("\n--- Test: Get User After Update ---")
	// Create and update a user
	createRes, _ := client.CreateUser(ctx, &pb.CreateUserRequest{
		Name:  "To Update",
		Email: "to.update@example.com",
	})

	client.UpdateUser(ctx, &pb.UpdateUserRequest{
		Id:    createRes.User.Id,
		Name:  "Updated Successfully",
		Email: "updated@example.com",
	})

	res, err := client.GetUser(ctx, &pb.GetUserRequest{
		Id: createRes.User.Id,
	})
	if err != nil {
		log.Printf("GetUser failed: %v", err)
		return
	}
	fmt.Printf("Retrieved updated user:\n  ID:    %s\n  Name:  %s\n  Email: %s\n", res.User.Id, res.User.Name, res.User.Email)
}

func testDeleteUser(ctx context.Context) {
	fmt.Println("\n--- Test: Delete User ---")
	// First create a user to delete
	createRes, err := client.CreateUser(ctx, &pb.CreateUserRequest{
		Name:  "To Delete",
		Email: "to.delete@example.com",
	})
	if err != nil {
		log.Printf("CreateUser failed: %v", err)
		return
	}

	res, err := client.DeleteUser(ctx, &pb.DeleteUserRequest{
		Id: createRes.User.Id,
	})
	if err != nil {
		log.Printf("DeleteUser failed: %v", err)
		return
	}
	fmt.Printf("DeleteUser response: success=%v\n", res.Success)

	// Verify user is deleted
	_, err = client.GetUser(ctx, &pb.GetUserRequest{Id: createRes.User.Id})
	if err != nil {
		st, _ := status.FromError(err)
		fmt.Printf("User successfully deleted (GetUser returned: %s)\n", st.Message())
	}
}

func testGetUserNotFound(ctx context.Context) {
	fmt.Println("\n--- Test: Get User (Not Found) ---")
	_, err := client.GetUser(ctx, &pb.GetUserRequest{
		Id: "non-existent-id-12345",
	})
	if err != nil {
		st, _ := status.FromError(err)
		fmt.Printf("Expected error received: %s\n", st.Message())
	} else {
		fmt.Println("ERROR: Expected not found error, but got success")
	}
}

func testListUsersWithPagination(ctx context.Context) {
	fmt.Println("\n--- Test: Cursor-Based Pagination ---")
	fmt.Println("Fetching users in pages of 5...\n")

	var currentCursor string
	pageNumber := 1
	totalFetched := 0

	for {
		res, err := client.ListUsers(ctx, &pb.ListUsersRequest{
			Limit:  5,
			Cursor: currentCursor,
		})
		if err != nil {
			log.Printf("ListUsers failed: %v", err)
			return
		}

		fmt.Printf("Page %d (%d users):\n", pageNumber, len(res.Users))
		for _, u := range res.Users {
			fmt.Printf("  - [%s] %s (%s)\n", u.Id[:8], u.Name, u.Email)
		}
		fmt.Printf("  Next cursor: '%s'\n\n", res.NextCursor)

		totalFetched += len(res.Users)

		if res.NextCursor == "" {
			break
		}
		currentCursor = res.NextCursor
		pageNumber++
	}

	fmt.Printf("Pagination test complete. Total users fetched: %d\n", totalFetched)
}
