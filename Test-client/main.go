package main

import (
	"context"
	"log"
	"time"

	pb "github.com/Farabi99/Simple-gRPC-Exploration/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	conn, err := grpc.Dial("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()

	client := pb.NewUserServiceClient(conn)

	// We'll use a longer timeout since we are looping
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	log.Println("--- Testing Cursor-Based Pagination ---")

	var currentCursor string
	pageNumber := 1
	totalFetched := 0

	for {
		log.Printf("\nFetching Page %d (Cursor: '%s')", pageNumber, currentCursor)

		res, err := client.ListUsers(ctx, &pb.ListUsersRequest{
			Limit:  15, // Get 15 users at a time
			Cursor: currentCursor,
		})
		if err != nil {
			log.Fatalf("Error fetching list: %v", err)
		}

		// Print the users in this page
		for _, u := range res.Users {
			log.Printf(" - [%s] %s (%s)", u.Id[:8], u.Name, u.Email) // Print just first 8 chars of ID for readability
		}

		totalFetched += len(res.Users)
		log.Printf("-> Received %d users. Next Cursor: '%s'", len(res.Users), res.NextCursor)

		// If the server returns an empty NextCursor, we have reached the end of the database
		if res.NextCursor == "" {
			break
		}

		// Update the cursor for the next iteration
		currentCursor = res.NextCursor
		pageNumber++
	}

	log.Printf("\n--- Finished! Total Users Fetched: %d ---", totalFetched)
}
