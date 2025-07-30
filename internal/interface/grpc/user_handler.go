package grpc

import (
	"context"
	"fmt"
	"regexp"
	"strings"

	"github.com/aungmyozaw92/go-grpc-starter/internal/entity"
	"github.com/aungmyozaw92/go-grpc-starter/internal/usecase"
	"github.com/aungmyozaw92/go-grpc-starter/proto/userpb"
)

type UserHandler struct {
	userpb.UnimplementedUserServiceServer
	UserUseCase *usecase.UserUseCase
}

func NewUserHandler(userUseCase *usecase.UserUseCase) *UserHandler {
	return &UserHandler{UserUseCase: userUseCase}
}

// validateRegisterRequest validates the registration request
func (h *UserHandler) validateRegisterRequest(req *userpb.RegisterRequest) error {
	// Check required fields
	if strings.TrimSpace(req.Username) == "" {
		return NewValidationError(MsgUsernameRequired)
	}

	if strings.TrimSpace(req.Name) == "" {
		return NewValidationError(MsgNameRequired)
	}

	if strings.TrimSpace(req.Email) == "" {
		return NewValidationError(MsgEmailRequired)
	}

	if strings.TrimSpace(req.Password) == "" {
		return NewValidationError(MsgPasswordRequired)
	}

	// Validate email format
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	if !emailRegex.MatchString(req.Email) {
		return NewValidationError(MsgInvalidEmail)
	}

	// Validate username (alphanumeric and underscore only, 3-30 characters)
	usernameRegex := regexp.MustCompile(`^[a-zA-Z0-9_]{3,30}$`)
	if !usernameRegex.MatchString(req.Username) {
		return NewValidationError(MsgInvalidUsername)
	}

	// Validate password strength (minimum 6 characters)
	if len(req.Password) < 6 {
		return NewValidationError(MsgPasswordTooShort)
	}

	// Validate role_id (should be positive)
	if req.RoleId <= 0 {
		return NewValidationError(MsgInvalidRoleID)
	}

	return nil
}

func (h *UserHandler) Register(ctx context.Context, req *userpb.RegisterRequest) (*userpb.AuthResponse, error) {
	// Validate the request
	if err := h.validateRegisterRequest(req); err != nil {
		return nil, err
	}

	user := &entity.User{
		Username:   req.Username,
		Name:       req.Name,
		Email:      &req.Email,
		Phone:      req.Phone,
		Mobile:     req.Mobile,
		ImageURL:   req.ImageUrl,
		Password:   req.Password,
		IsActive:   &req.IsActive,
		RoleID:     int(req.RoleId),
	}

	token, err := h.UserUseCase.Register(user)
	if err != nil {
		return nil, NewInternalError(MsgUserRegistrationFailed)
	}

	return &userpb.AuthResponse{
		Success: true,
		Code:    string(CodeSuccess),
		Message: MsgUserRegistered,
		Token:   token,
	}, nil
}

func (h *UserHandler) Login(ctx context.Context, req *userpb.LoginRequest) (*userpb.AuthResponse, error) {
	// Debug logging to see what we receive
	fmt.Printf("DEBUG Login - Username: '%s', Password: '%s'\n", req.Username, req.Password)

	// Validate login request
	if strings.TrimSpace(req.Username) == "" {
		return nil, NewValidationError(MsgUsernameRequired)
	}

	if strings.TrimSpace(req.Password) == "" {
		return nil, NewValidationError(MsgPasswordRequired)
	}

	token, err := h.UserUseCase.Login(req.Username, req.Password)
	if err != nil {
		return nil, NewAuthenticationError(MsgInvalidCredentials)
	}

	return &userpb.AuthResponse{
		Success: true,
		Code:    string(CodeSuccess),
		Message: MsgUserLoggedIn,
		Token:   token,
	}, nil
}

func (h *UserHandler) GetProfile(ctx context.Context, req *userpb.ProfileRequest) (*userpb.ProfileResponse, error) {
	// Validate profile request
	if strings.TrimSpace(req.Token) == "" {
		return nil, NewValidationError(MsgTokenRequired)
	}

	user, err := h.UserUseCase.GetProfile(req.Token)
	if err != nil {
		return nil, NewAuthenticationError(MsgInvalidToken)
	}

	return &userpb.ProfileResponse{
		Success: true,
		Code:    string(CodeSuccess),
		Message: MsgProfileRetrieved,
		Data: &userpb.ProfileData{
			Id:        int32(user.ID),
			Username:  user.Username,
			Name:      user.Name,
			Email:     deref(user.Email),
			Phone:     user.Phone,
			Mobile:    user.Mobile,
			IsActive:  derefBool(user.IsActive),
			RoleId:    int32(user.RoleID),
			ImageUrl:  user.ImageURL,
			CreatedAt: user.CreatedAt,
			UpdatedAt: user.UpdatedAt,
		},
	}, nil
}

func deref(s *string) string {
	if s != nil {
		return *s
	}
	return ""
}

func derefBool(b *bool) bool {
	if b != nil {
		return *b
	}
	return false
}