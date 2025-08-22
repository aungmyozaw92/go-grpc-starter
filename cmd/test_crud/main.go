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

	// Step 1: Login to get a token
	fmt.Println("=== Step 1: Login to get authentication token ===")
	loginReq := &userpb.LoginRequest{
		Username: "testuser",
		Password: "password123",
	}

	loginResp, err := client.Login(ctx, loginReq)
	if err != nil {
		// If login fails, try to register first
		fmt.Println("Login failed, trying to register a user first...")
		registerReq := &userpb.RegisterRequest{
			Username: "testuser",
			Name:     "Test Admin",
			Email:    "admin@example.com",
			Phone:    "555-0001",
			Mobile:   "555-0002",
			ImageUrl: "https://example.com/admin.jpg",
			Password: "password123",
			IsActive: true,
			RoleId:   1,
		}

		registerResp, err := client.Register(ctx, registerReq)
		if err != nil {
			log.Fatalf("Failed to register: %v", err)
		}
		loginResp = registerResp
		fmt.Printf("✅ Registered user successfully\n")
	}

	if !loginResp.Success {
		log.Fatalf("Login failed: %s", loginResp.Message)
	}

	token := loginResp.Token
	fmt.Printf("✅ Authentication successful\n")
	fmt.Printf("Token: %s...\n\n", token[:20])

	// Step 2: Create a new user
	fmt.Println("=== Step 2: Create User ===")
	createReq := &userpb.CreateUserRequest{
		Token:    token,
		Username: "newuser123",
		Name:     "New User",
		Email:    "newuser@example.com",
		Phone:    "555-1001",
		Mobile:   "555-1002",
		ImageUrl: "https://example.com/newuser.jpg",
		Password: "newpassword123",
		IsActive: true,
		RoleId:   2,
	}

	createResp, err := client.CreateUser(ctx, createReq)
	if err != nil {
		fmt.Printf("❌ Create user failed: %v\n\n", err)
	} else {
		fmt.Printf("✅ User created successfully\n")
		fmt.Printf("Created User ID: %d\n", createResp.Data.Id)
		fmt.Printf("Username: %s\n", createResp.Data.Username)
		fmt.Printf("Name: %s\n", createResp.Data.Name)
		fmt.Printf("Email: %s\n\n", createResp.Data.Email)
	}

	// Step 3: Get user by ID
	var userID int32 = 1 // Try to get the first user
	if createResp != nil && createResp.Success {
		userID = createResp.Data.Id
	}

	fmt.Printf("=== Step 3: Get User (ID: %d) ===\n", userID)
	getUserReq := &userpb.GetUserRequest{
		Token:  token,
		UserId: userID,
	}

	getUserResp, err := client.GetUser(ctx, getUserReq)
	if err != nil {
		fmt.Printf("❌ Get user failed: %v\n\n", err)
	} else {
		fmt.Printf("✅ User retrieved successfully\n")
		user := getUserResp.Data
		fmt.Printf("ID: %d\n", user.Id)
		fmt.Printf("Username: %s\n", user.Username)
		fmt.Printf("Name: %s\n", user.Name)
		fmt.Printf("Email: %s\n", user.Email)
		fmt.Printf("Phone: %s\n", user.Phone)
		fmt.Printf("Mobile: %s\n", user.Mobile)
		fmt.Printf("Active: %v\n", user.IsActive)
		fmt.Printf("Role ID: %d\n", user.RoleId)
		fmt.Printf("Created: %s\n", user.CreatedAt)
		fmt.Printf("Updated: %s\n\n", user.UpdatedAt)
	}

	// Step 4: Update user
	fmt.Printf("=== Step 4: Update User (ID: %d) ===\n", userID)
	updateReq := &userpb.UpdateUserRequest{
		Token:    token,
		UserId:   userID,
		Username: "updateduser123",
		Name:     "Updated User Name",
		Email:    "updated@example.com",
		Phone:    "555-2001",
		Mobile:   "555-2002",
		ImageUrl: "https://example.com/updated.jpg",
		IsActive: true,
		RoleId:   3,
	}

	updateResp, err := client.UpdateUser(ctx, updateReq)
	if err != nil {
		fmt.Printf("❌ Update user failed: %v\n\n", err)
	} else {
		fmt.Printf("✅ User updated successfully\n")
		user := updateResp.Data
		fmt.Printf("ID: %d\n", user.Id)
		fmt.Printf("Username: %s (updated)\n", user.Username)
		fmt.Printf("Name: %s (updated)\n", user.Name)
		fmt.Printf("Email: %s (updated)\n", user.Email)
		fmt.Printf("Phone: %s (updated)\n", user.Phone)
		fmt.Printf("Role ID: %d (updated)\n\n", user.RoleId)
	}

	// Step 5: Get user list to see all users
	fmt.Println("=== Step 5: Get User List ===")
	userListReq := &userpb.UserListRequest{
		Token: token,
		Page:  1,
		Limit: 10,
	}

	userListResp, err := client.GetUserList(ctx, userListReq)
	if err != nil {
		fmt.Printf("❌ Get user list failed: %v\n\n", err)
	} else {
		fmt.Printf("✅ User list retrieved successfully\n")
		fmt.Printf("Total users: %d\n", userListResp.Data.Pagination.TotalCount)
		fmt.Printf("Users:\n")
		for i, user := range userListResp.Data.Users {
			fmt.Printf("  %d. %s (%s) - %s\n", i+1, user.Username, user.Name, user.Email)
		}
		fmt.Println()
	}

	// Step 6: Create another user to test deletion
	fmt.Println("=== Step 6: Create User for Deletion Test ===")
	createDeleteReq := &userpb.CreateUserRequest{
		Token:    token,
		Username: "deleteme123",
		Name:     "Delete Me User",
		Email:    "deleteme@example.com",
		Phone:    "555-9001",
		Mobile:   "555-9002",
		Password: "deletepassword",
		IsActive: true,
		RoleId:   1,
	}

	createDeleteResp, err := client.CreateUser(ctx, createDeleteReq)
	var deleteUserID int32
	if err != nil {
		fmt.Printf("❌ Create user for deletion failed: %v\n", err)
		deleteUserID = userID // Use existing user ID
	} else {
		deleteUserID = createDeleteResp.Data.Id
		fmt.Printf("✅ User created for deletion test (ID: %d)\n", deleteUserID)
	}
	fmt.Println()

	// Step 7: Delete user
	fmt.Printf("=== Step 7: Delete User (ID: %d) ===\n", deleteUserID)
	deleteReq := &userpb.DeleteUserRequest{
		Token:  token,
		UserId: deleteUserID,
	}

	deleteResp, err := client.DeleteUser(ctx, deleteReq)
	if err != nil {
		fmt.Printf("❌ Delete user failed: %v\n\n", err)
	} else {
		fmt.Printf("✅ User deleted successfully\n")
		fmt.Printf("Message: %s\n\n", deleteResp.Message)
	}

	// Step 8: Try to get the deleted user (should fail)
	fmt.Printf("=== Step 8: Verify Deletion (Get Deleted User) ===\n")
	getDeletedReq := &userpb.GetUserRequest{
		Token:  token,
		UserId: deleteUserID,
	}

	_, err = client.GetUser(ctx, getDeletedReq)
	if err != nil {
		fmt.Printf("✅ Expected error - user not found: %v\n\n", err)
	} else {
		fmt.Printf("❌ Unexpected - deleted user still exists\n\n")
	}

	// Step 9: Final user list
	fmt.Println("=== Step 9: Final User List ===")
	finalListResp, err := client.GetUserList(ctx, userListReq)
	if err != nil {
		fmt.Printf("❌ Get final user list failed: %v\n", err)
	} else {
		fmt.Printf("✅ Final user list:\n")
		fmt.Printf("Total users: %d\n", finalListResp.Data.Pagination.TotalCount)
		for i, user := range finalListResp.Data.Users {
			fmt.Printf("  %d. ID:%d | %s (%s) - %s\n", i+1, user.Id, user.Username, user.Name, user.Email)
		}
	}

	fmt.Println("\n🎉 CRUD Operations Test Completed!")
}
