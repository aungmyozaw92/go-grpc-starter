package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/aungmyozaw92/go-grpc-starter/proto/userpb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	// Connect to the gRPC server
	conn, err := grpc.Dial("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Failed to connect: %v", err)
	}
	defer conn.Close()

	client := userpb.NewUserServiceClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// First, register multiple users for testing
	fmt.Println("=== Creating Test Users ===")

	users := []struct {
		username string
		name     string
		email    string
	}{
		{"john_doe", "John Doe", "john@example.com"},
		{"jane_smith", "Jane Smith", "jane@example.com"},
		{"bob_wilson", "Bob Wilson", "bob@example.com"},
		{"alice_brown", "Alice Brown", "alice@example.com"},
		{"charlie_davis", "Charlie Davis", "charlie@example.com"},
	}

	var token string

	// Register users
	for i, user := range users {
		registerReq := &userpb.RegisterRequest{
			Username: user.username,
			Name:     user.name,
			Email:    user.email,
			Phone:    fmt.Sprintf("555-000%d", i+1),
			Mobile:   fmt.Sprintf("555-100%d", i+1),
			ImageUrl: "https://example.com/avatar.jpg",
			Password: "password123",
			IsActive: true,
			RoleId:   1,
		}

		registerResp, err := client.Register(ctx, registerReq)
		if err != nil {
			fmt.Printf("Failed to register %s: %v\n", user.username, err)
			continue
		}

		fmt.Printf("✅ Registered: %s\n", user.username)

		// Save token from first user for testing
		if i == 0 {
			token = registerResp.Token
		}
	}

	if token == "" {
		log.Fatal("Failed to get a valid token for testing")
	}

	fmt.Printf("\nUsing token: %s...\n\n", token[:20])

	// Test GetUserList with different scenarios
	testScenarios := []struct {
		name   string
		page   int32
		limit  int32
		search string
	}{
		{"Default pagination (page 1, limit 10)", 1, 10, ""},
		{"Page 1 with limit 3", 1, 3, ""},
		{"Page 2 with limit 3", 2, 3, ""},
		{"Search for 'john'", 1, 10, "john"},
		{"Search for 'example.com'", 1, 10, "example.com"},
		{"Search for 'Smith'", 1, 10, "Smith"},
		{"No results search", 1, 10, "nonexistent"},
		{"Large limit", 1, 100, ""},
	}

	for _, scenario := range testScenarios {
		fmt.Printf("=== %s ===\n", scenario.name)

		userListReq := &userpb.UserListRequest{
			Token:  token,
			Page:   scenario.page,
			Limit:  scenario.limit,
			Search: scenario.search,
		}

		userListResp, err := client.GetUserList(ctx, userListReq)
		if err != nil {
			fmt.Printf("❌ Error: %v\n\n", err)
			continue
		}

		if !userListResp.Success {
			fmt.Printf("❌ Failed: %s\n\n", userListResp.Message)
			continue
		}

		data := userListResp.Data
		pagination := data.Pagination

		fmt.Printf("✅ Success: %s\n", userListResp.Message)
		fmt.Printf("📊 Pagination Info:\n")
		fmt.Printf("   Current Page: %d\n", pagination.CurrentPage)
		fmt.Printf("   Per Page: %d\n", pagination.PerPage)
		fmt.Printf("   Total Pages: %d\n", pagination.TotalPages)
		fmt.Printf("   Total Count: %d\n", pagination.TotalCount)
		fmt.Printf("   Has Next: %v\n", pagination.HasNext)
		fmt.Printf("   Has Previous: %v\n", pagination.HasPrev)

		fmt.Printf("👥 Users (%d results):\n", len(data.Users))
		for i, user := range data.Users {
			fmt.Printf("   %d. ID:%d | %s (%s) | %s | Active:%v\n",
				i+1, user.Id, user.Username, user.Name, user.Email, user.IsActive)
		}
		fmt.Println()
	}

	// Test with invalid token
	fmt.Println("=== Testing Invalid Token ===")
	invalidReq := &userpb.UserListRequest{
		Token: "invalid_token",
		Page:  1,
		Limit: 10,
	}

	_, err = client.GetUserList(ctx, invalidReq)
	if err != nil {
		fmt.Printf("✅ Expected error with invalid token: %v\n", err)
	} else {
		fmt.Printf("❌ Should have failed with invalid token\n")
	}
}
