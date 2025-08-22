package grpc

import (
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// Response messages for consistency
const (
	// Success messages
	MsgUserRegistered    = "User registered successfully"
	MsgUserLoggedIn      = "User logged in successfully"
	MsgProfileRetrieved  = "User profile retrieved successfully"
	MsgUserListRetrieved = "User list retrieved successfully"
	MsgUserRetrieved     = "User retrieved successfully"
	MsgUserCreated       = "User created successfully"
	MsgUserUpdated       = "User updated successfully"
	MsgUserDeleted       = "User deleted successfully"

	// Error messages - Validation
	MsgUsernameRequired = "Username is required"
	MsgNameRequired     = "Name is required"
	MsgEmailRequired    = "Email is required"
	MsgPasswordRequired = "Password is required"
	MsgTokenRequired    = "Authentication token is required"
	MsgInvalidEmail     = "Invalid email format"
	MsgInvalidUsername  = "Username must be 3-30 characters and contain only letters, numbers, and underscores"
	MsgPasswordTooShort = "Password must be at least 6 characters long"
	MsgInvalidRoleID    = "Role ID must be a positive integer"
	MsgUsernameExists   = "Username already exists"
	MsgEmailExists      = "Email address already exists"

	// Error messages - Authentication/Authorization
	MsgInvalidCredentials = "Invalid username or password"
	MsgInvalidToken       = "Invalid or expired token"
	MsgUnauthorized       = "Unauthorized access"

	// Error messages - Internal/System
	MsgUserRegistrationFailed = "Failed to register user"
	MsgUserLoginFailed        = "Failed to authenticate user"
	MsgProfileRetrievalFailed = "Failed to retrieve user profile"
	MsgUserNotFound           = "User not found"
	MsgUserCreationFailed     = "Failed to create user"
	MsgUserUpdateFailed       = "Failed to update user"
	MsgUserDeletionFailed     = "Failed to delete user"
	MsgInternalError          = "Internal server error"
	MsgDatabaseError          = "Database operation failed"
)

// Error helper functions for consistent error responses
func NewValidationError(message string) error {
	return status.Errorf(codes.InvalidArgument, message)
}

func NewAuthenticationError(message string) error {
	return status.Errorf(codes.Unauthenticated, message)
}

func NewAuthorizationError(message string) error {
	return status.Errorf(codes.PermissionDenied, message)
}

func NewInternalError(message string) error {
	return status.Errorf(codes.Internal, message)
}

func NewNotFoundError(message string) error {
	return status.Errorf(codes.NotFound, message)
}

func NewAlreadyExistsError(message string) error {
	return status.Errorf(codes.AlreadyExists, message)
}

// Response code mapping for different scenarios
type ResponseCode string

const (
	CodeSuccess             ResponseCode = "SUCCESS"
	CodeValidationError     ResponseCode = "VALIDATION_ERROR"
	CodeAuthenticationError ResponseCode = "AUTHENTICATION_ERROR"
	CodeAuthorizationError  ResponseCode = "AUTHORIZATION_ERROR"
	CodeNotFound            ResponseCode = "NOT_FOUND"
	CodeAlreadyExists       ResponseCode = "ALREADY_EXISTS"
	CodeInternalError       ResponseCode = "INTERNAL_ERROR"
)

// Helper function to get response code from gRPC error
func GetResponseCode(err error) ResponseCode {
	if err == nil {
		return CodeSuccess
	}

	st, ok := status.FromError(err)
	if !ok {
		return CodeInternalError
	}

	switch st.Code() {
	case codes.InvalidArgument:
		return CodeValidationError
	case codes.Unauthenticated:
		return CodeAuthenticationError
	case codes.PermissionDenied:
		return CodeAuthorizationError
	case codes.NotFound:
		return CodeNotFound
	case codes.AlreadyExists:
		return CodeAlreadyExists
	case codes.Internal:
		return CodeInternalError
	default:
		return CodeInternalError
	}
}
