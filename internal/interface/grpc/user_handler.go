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
		Username: req.Username,
		Name:     req.Name,
		Email:    &req.Email,
		Phone:    req.Phone,
		Mobile:   req.Mobile,
		ImageURL: req.ImageUrl,
		Password: req.Password,
		IsActive: &req.IsActive,
		RoleID:   int(req.RoleId),
	}

	token, err := h.UserUseCase.Register(user)
	if err != nil {
		if err.Error() == "username already exists" {
			return nil, NewAlreadyExistsError(MsgUsernameExists)
		}
		if err.Error() == "email already exists" {
			return nil, NewAlreadyExistsError(MsgEmailExists)
		}
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
			CreatedAt: user.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
			UpdatedAt: user.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
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

func (h *UserHandler) GetUserList(ctx context.Context, req *userpb.UserListRequest) (*userpb.UserListResponse, error) {
	// Validate token
	if strings.TrimSpace(req.Token) == "" {
		return nil, NewValidationError(MsgTokenRequired)
	}

	// Set default pagination values
	page := int(req.Page)
	limit := int(req.Limit)
	search := strings.TrimSpace(req.Search)

	if page <= 0 {
		page = 1
	}
	if limit <= 0 {
		limit = 10
	}

	// Get user list from usecase
	result, err := h.UserUseCase.GetUserList(req.Token, page, limit, search)
	if err != nil {
		return nil, NewAuthenticationError(MsgInvalidToken)
	}

	// Convert users to protobuf format
	var pbUsers []*userpb.UserData
	for _, user := range result.Users {
		pbUser := &userpb.UserData{
			Id:        int32(user.ID),
			Username:  user.Username,
			Name:      user.Name,
			Email:     deref(user.Email),
			Phone:     user.Phone,
			Mobile:    user.Mobile,
			IsActive:  derefBool(user.IsActive),
			RoleId:    int32(user.RoleID),
			ImageUrl:  user.ImageURL,
			CreatedAt: user.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
			UpdatedAt: user.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
		}
		pbUsers = append(pbUsers, pbUser)
	}

	// Create pagination metadata
	pagination := &userpb.PaginationMeta{
		CurrentPage: int32(result.Pagination.CurrentPage),
		PerPage:     int32(result.Pagination.PerPage),
		TotalPages:  int32(result.Pagination.TotalPages),
		TotalCount:  int32(result.Pagination.TotalCount),
		HasNext:     result.Pagination.HasNext,
		HasPrev:     result.Pagination.HasPrev,
	}

	return &userpb.UserListResponse{
		Success: true,
		Code:    string(CodeSuccess),
		Message: MsgUserListRetrieved,
		Data: &userpb.UserListData{
			Users:      pbUsers,
			Pagination: pagination,
		},
	}, nil
}

func (h *UserHandler) GetUser(ctx context.Context, req *userpb.GetUserRequest) (*userpb.GetUserResponse, error) {
	// Validate token
	if strings.TrimSpace(req.Token) == "" {
		return nil, NewValidationError(MsgTokenRequired)
	}

	// Validate user ID
	if req.UserId <= 0 {
		return nil, NewValidationError("User ID must be positive")
	}

	// Get user from usecase
	user, err := h.UserUseCase.GetUser(req.Token, int(req.UserId))
	if err != nil {
		return nil, NewNotFoundError(MsgUserNotFound)
	}

	// Convert to protobuf format
	userData := &userpb.UserData{
		Id:        int32(user.ID),
		Username:  user.Username,
		Name:      user.Name,
		Email:     deref(user.Email),
		Phone:     user.Phone,
		Mobile:    user.Mobile,
		IsActive:  derefBool(user.IsActive),
		RoleId:    int32(user.RoleID),
		ImageUrl:  user.ImageURL,
		CreatedAt: user.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt: user.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}

	return &userpb.GetUserResponse{
		Success: true,
		Code:    string(CodeSuccess),
		Message: MsgUserRetrieved,
		Data:    userData,
	}, nil
}

func (h *UserHandler) CreateUser(ctx context.Context, req *userpb.CreateUserRequest) (*userpb.CreateUserResponse, error) {
	// Validate token
	if strings.TrimSpace(req.Token) == "" {
		return nil, NewValidationError(MsgTokenRequired)
	}

	// Validate required fields
	if strings.TrimSpace(req.Username) == "" {
		return nil, NewValidationError(MsgUsernameRequired)
	}
	if strings.TrimSpace(req.Name) == "" {
		return nil, NewValidationError(MsgNameRequired)
	}
	if strings.TrimSpace(req.Email) == "" {
		return nil, NewValidationError(MsgEmailRequired)
	}
	if strings.TrimSpace(req.Password) == "" {
		return nil, NewValidationError(MsgPasswordRequired)
	}

	// Create user entity
	user := &entity.User{
		Username: req.Username,
		Name:     req.Name,
		Email:    &req.Email,
		Phone:    req.Phone,
		Mobile:   req.Mobile,
		ImageURL: req.ImageUrl,
		Password: req.Password,
		IsActive: &req.IsActive,
		RoleID:   int(req.RoleId),
	}

	// Create user via usecase
	createdUser, err := h.UserUseCase.CreateUser(req.Token, user)
	if err != nil {
		if err.Error() == "username already exists" {
			return nil, NewAlreadyExistsError(MsgUsernameExists)
		}
		if err.Error() == "email already exists" {
			return nil, NewAlreadyExistsError(MsgEmailExists)
		}
		return nil, NewInternalError(MsgUserCreationFailed)
	}

	// Convert to protobuf format
	userData := &userpb.UserData{
		Id:        int32(createdUser.ID),
		Username:  createdUser.Username,
		Name:      createdUser.Name,
		Email:     deref(createdUser.Email),
		Phone:     createdUser.Phone,
		Mobile:    createdUser.Mobile,
		IsActive:  derefBool(createdUser.IsActive),
		RoleId:    int32(createdUser.RoleID),
		ImageUrl:  createdUser.ImageURL,
		CreatedAt: createdUser.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt: createdUser.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}

	return &userpb.CreateUserResponse{
		Success: true,
		Code:    string(CodeSuccess),
		Message: MsgUserCreated,
		Data:    userData,
	}, nil
}

func (h *UserHandler) UpdateUser(ctx context.Context, req *userpb.UpdateUserRequest) (*userpb.UpdateUserResponse, error) {
	// Validate token
	if strings.TrimSpace(req.Token) == "" {
		return nil, NewValidationError(MsgTokenRequired)
	}

	// Validate user ID
	if req.UserId <= 0 {
		return nil, NewValidationError("User ID must be positive")
	}

	// Validate required fields
	if strings.TrimSpace(req.Username) == "" {
		return nil, NewValidationError(MsgUsernameRequired)
	}
	if strings.TrimSpace(req.Name) == "" {
		return nil, NewValidationError(MsgNameRequired)
	}
	if strings.TrimSpace(req.Email) == "" {
		return nil, NewValidationError(MsgEmailRequired)
	}

	// Create update data
	updateData := &entity.User{
		Username: req.Username,
		Name:     req.Name,
		Email:    &req.Email,
		Phone:    req.Phone,
		Mobile:   req.Mobile,
		ImageURL: req.ImageUrl,
		IsActive: &req.IsActive,
		RoleID:   int(req.RoleId),
	}

	// Update user via usecase
	updatedUser, err := h.UserUseCase.UpdateUser(req.Token, int(req.UserId), updateData)
	if err != nil {
		if err.Error() == "username already exists" {
			return nil, NewAlreadyExistsError(MsgUsernameExists)
		}
		if err.Error() == "email already exists" {
			return nil, NewAlreadyExistsError(MsgEmailExists)
		}
		return nil, NewInternalError(MsgUserUpdateFailed)
	}

	// Convert to protobuf format
	userData := &userpb.UserData{
		Id:        int32(updatedUser.ID),
		Username:  updatedUser.Username,
		Name:      updatedUser.Name,
		Email:     deref(updatedUser.Email),
		Phone:     updatedUser.Phone,
		Mobile:    updatedUser.Mobile,
		IsActive:  derefBool(updatedUser.IsActive),
		RoleId:    int32(updatedUser.RoleID),
		ImageUrl:  updatedUser.ImageURL,
		CreatedAt: updatedUser.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt: updatedUser.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}

	return &userpb.UpdateUserResponse{
		Success: true,
		Code:    string(CodeSuccess),
		Message: MsgUserUpdated,
		Data:    userData,
	}, nil
}

func (h *UserHandler) DeleteUser(ctx context.Context, req *userpb.DeleteUserRequest) (*userpb.DeleteUserResponse, error) {
	// Validate token
	if strings.TrimSpace(req.Token) == "" {
		return nil, NewValidationError(MsgTokenRequired)
	}

	// Validate user ID
	if req.UserId <= 0 {
		return nil, NewValidationError("User ID must be positive")
	}

	// Delete user via usecase
	err := h.UserUseCase.DeleteUser(req.Token, int(req.UserId))
	if err != nil {
		return nil, NewInternalError(MsgUserDeletionFailed)
	}

	return &userpb.DeleteUserResponse{
		Success: true,
		Code:    string(CodeSuccess),
		Message: MsgUserDeleted,
	}, nil
}
