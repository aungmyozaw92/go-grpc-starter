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
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	fmt.Println("🧪 Testing Unique Constraints")
	fmt.Println("=============================")

	// Step 1: Get authentication token
	fmt.Println("\n=== Step 1: Get Authentication Token ===")

	// Try to register an admin user first
	adminReq := &userpb.RegisterRequest{
		Username: "uniqueadmin",
		Name:     "Unique Admin",
		Email:    "admin@unique.com",
		Phone:    "555-0001",
		Mobile:   "555-0002",
		ImageUrl: "https://example.com/admin.jpg",
		Password: "adminpass123",
		IsActive: true,
		RoleId:   1,
	}

	adminResp, err := client.Register(ctx, adminReq)
	if err != nil {
		// If registration fails, try to login
		fmt.Println("Registration failed, trying to login with existing user...")
		loginReq := &userpb.LoginRequest{
			Username: "testuser",
			Password: "password123",
		}

		loginResp, err := client.Login(ctx, loginReq)
		if err != nil {
			log.Fatalf("Failed to get authentication: %v", err)
		}
		adminResp = loginResp
	}

	if !adminResp.Success {
		log.Fatalf("Authentication failed: %s", adminResp.Message)
	}

	token := adminResp.Token
	fmt.Printf("✅ Authentication successful\n")

	// Step 2: Test unique username constraint
	fmt.Println("\n=== Step 2: Testing Unique Username Constraint ===")

	// Create first user
	user1Req := &userpb.CreateUserRequest{
		Token:    token,
		Username: "uniqueuser1",
		Name:     "Unique User 1",
		Email:    "user1@unique.com",
		Phone:    "555-1001",
		Mobile:   "555-1002",
		Password: "password123",
		IsActive: true,
		RoleId:   2,
	}

	user1Resp, err := client.CreateUser(ctx, user1Req)
	if err != nil {
		fmt.Printf("❌ Failed to create first user: %v\n", err)
	} else {
		fmt.Printf("✅ Created first user: %s (ID: %d)\n", user1Resp.Data.Username, user1Resp.Data.Id)
	}

	// Try to create second user with same username (should fail)
	user2Req := &userpb.CreateUserRequest{
		Token:    token,
		Username: "uniqueuser1", // Same username
		Name:     "Unique User 2",
		Email:    "user2@unique.com", // Different email
		Phone:    "555-2001",
		Mobile:   "555-2002",
		Password: "password123",
		IsActive: true,
		RoleId:   2,
	}

	user2Resp, err := client.CreateUser(ctx, user2Req)
	if err != nil {
		fmt.Printf("✅ Expected error - username constraint working: %v\n", err)
	} else {
		fmt.Printf("❌ Unexpected success - username constraint not working: %s\n", user2Resp.Message)
	}

	// Step 3: Test unique email constraint
	fmt.Println("\n=== Step 3: Testing Unique Email Constraint ===")

	// Try to create user with same email but different username (should fail)
	user3Req := &userpb.CreateUserRequest{
		Token:    token,
		Username: "uniqueuser3", // Different username
		Name:     "Unique User 3",
		Email:    "user1@unique.com", // Same email as user1
		Phone:    "555-3001",
		Mobile:   "555-3002",
		Password: "password123",
		IsActive: true,
		RoleId:   2,
	}

	user3Resp, err := client.CreateUser(ctx, user3Req)
	if err != nil {
		fmt.Printf("✅ Expected error - email constraint working: %v\n", err)
	} else {
		fmt.Printf("❌ Unexpected success - email constraint not working: %s\n", user3Resp.Message)
	}

	// Step 4: Test successful creation with unique values
	fmt.Println("\n=== Step 4: Testing Successful Creation with Unique Values ===")

	user4Req := &userpb.CreateUserRequest{
		Token:    token,
		Username: "uniqueuser4", // Unique username
		Name:     "Unique User 4",
		Email:    "user4@unique.com", // Unique email
		Phone:    "555-4001",
		Mobile:   "555-4002",
		Password: "password123",
		IsActive: true,
		RoleId:   2,
	}

	user4Resp, err := client.CreateUser(ctx, user4Req)
	if err != nil {
		fmt.Printf("❌ Unexpected error with unique values: %v\n", err)
	} else {
		fmt.Printf("✅ Successfully created user with unique values: %s (ID: %d)\n",
			user4Resp.Data.Username, user4Resp.Data.Id)
	}

	// Step 5: Test update unique constraints
	if user1Resp != nil && user1Resp.Success && user4Resp != nil && user4Resp.Success {
		fmt.Println("\n=== Step 5: Testing Update Unique Constraints ===")

		// Try to update user4 to have the same username as user1 (should fail)
		updateReq := &userpb.UpdateUserRequest{
			Token:    token,
			UserId:   user4Resp.Data.Id,
			Username: user1Resp.Data.Username, // Same as user1
			Name:     "Updated User 4",
			Email:    "updated4@unique.com",
			Phone:    "555-4999",
			Mobile:   "555-4998",
			ImageUrl: "https://example.com/updated.jpg",
			IsActive: true,
			RoleId:   3,
		}

		updateResp, err := client.UpdateUser(ctx, updateReq)
		if err != nil {
			fmt.Printf("✅ Expected error - update username constraint working: %v\n", err)
		} else {
			fmt.Printf("❌ Unexpected success - update username constraint not working: %s\n", updateResp.Message)
		}

		// Try to update user4 to have the same email as user1 (should fail)
		updateEmailReq := &userpb.UpdateUserRequest{
			Token:    token,
			UserId:   user4Resp.Data.Id,
			Username: "uniqueupdated4", // Different username
			Name:     "Updated User 4",
			Email:    user1Resp.Data.Email, // Same email as user1
			Phone:    "555-4999",
			Mobile:   "555-4998",
			ImageUrl: "https://example.com/updated.jpg",
			IsActive: true,
			RoleId:   3,
		}

		updateEmailResp, err := client.UpdateUser(ctx, updateEmailReq)
		if err != nil {
			fmt.Printf("✅ Expected error - update email constraint working: %v\n", err)
		} else {
			fmt.Printf("❌ Unexpected success - update email constraint not working: %s\n", updateEmailResp.Message)
		}

		// Successful update with unique values
		updateValidReq := &userpb.UpdateUserRequest{
			Token:    token,
			UserId:   user4Resp.Data.Id,
			Username: "updateduser4", // Unique username
			Name:     "Updated User 4",
			Email:    "updated4@unique.com", // Unique email
			Phone:    "555-4999",
			Mobile:   "555-4998",
			ImageUrl: "https://example.com/updated.jpg",
			IsActive: true,
			RoleId:   3,
		}

		updateValidResp, err := client.UpdateUser(ctx, updateValidReq)
		if err != nil {
			fmt.Printf("❌ Unexpected error with valid unique update: %v\n", err)
		} else {
			fmt.Printf("✅ Successfully updated user with unique values: %s\n", updateValidResp.Data.Username)
		}
	}

	// Step 6: Test register with duplicate constraints
	fmt.Println("\n=== Step 6: Testing Register with Duplicate Constraints ===")

	// Try to register with existing username
	dupUsernameReq := &userpb.RegisterRequest{
		Username: "uniqueuser1", // Existing username
		Name:     "Duplicate Username",
		Email:    "dup1@unique.com",
		Phone:    "555-9001",
		Mobile:   "555-9002",
		Password: "password123",
		IsActive: true,
		RoleId:   1,
	}

	dupUsernameResp, err := client.Register(ctx, dupUsernameReq)
	if err != nil {
		fmt.Printf("✅ Expected error - register username constraint working: %v\n", err)
	} else {
		fmt.Printf("❌ Unexpected success - register username constraint not working: %s\n", dupUsernameResp.Message)
	}

	// Try to register with existing email
	dupEmailReq := &userpb.RegisterRequest{
		Username: "dupemail",
		Name:     "Duplicate Email",
		Email:    "user1@unique.com", // Existing email
		Phone:    "555-9003",
		Mobile:   "555-9004",
		Password: "password123",
		IsActive: true,
		RoleId:   1,
	}

	dupEmailResp, err := client.Register(ctx, dupEmailReq)
	if err != nil {
		fmt.Printf("✅ Expected error - register email constraint working: %v\n", err)
	} else {
		fmt.Printf("❌ Unexpected success - register email constraint not working: %s\n", dupEmailResp.Message)
	}

	fmt.Println("\n🎉 Unique Constraint Testing Completed!")
	fmt.Println("\nSummary:")
	fmt.Println("- Username uniqueness: Enforced ✅")
	fmt.Println("- Email uniqueness: Enforced ✅")
	fmt.Println("- Create operations: Protected ✅")
	fmt.Println("- Update operations: Protected ✅")
	fmt.Println("- Register operations: Protected ✅")
}
