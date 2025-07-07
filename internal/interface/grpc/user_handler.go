package grpc

import (
	"context"

	"github.com/aungmyozaw92/go-grpc-starter/internal/entity"
	"github.com/aungmyozaw92/go-grpc-starter/internal/usecase"
	"github.com/aungmyozaw92/go-grpc-starter/proto/userpb"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type UserHandler struct {
	userpb.UnimplementedUserServiceServer
	UserUseCase *usecase.UserUseCase
}

func NewUserHandler(userUseCase *usecase.UserUseCase) *UserHandler {
	return &UserHandler{UserUseCase: userUseCase}
}

func (h *UserHandler) Register(ctx context.Context, req *userpb.RegisterRequest) (*userpb.AuthResponse, error) {
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
		// return nil, err
		return nil, status.Errorf(codes.Internal, "Failed to register user: %v", err)
	}

	return &userpb.AuthResponse{
		Token: token,
		Message: "User registered successfully",
	}, nil
}	

func (h *UserHandler) Login(ctx context.Context, req *userpb.LoginRequest) (*userpb.AuthResponse, error) {
  token, err := h.UserUseCase.Login(req.Username, req.Password)
  if err != nil {
    return nil, err
  }
  return &userpb.AuthResponse{Token: token, Message: "Login successful"}, nil
}

func (h *UserHandler) GetProfile(ctx context.Context, req *userpb.ProfileRequest) (*userpb.ProfileResponse, error) {
  user, err := h.UserUseCase.GetProfile(req.Token)
  if err != nil {
    return nil, err
  }
  return &userpb.ProfileResponse{
    Id:         int32(user.ID),
    Username:   user.Username,
    Name:       user.Name,
    Email:      deref(user.Email),
    Phone:      user.Phone,
    Mobile:     user.Mobile,
    IsActive:   derefBool(user.IsActive),
    RoleId:     int32(user.RoleID),
    ImageUrl:   user.ImageURL,
    CreatedAt:  user.CreatedAt,
    UpdatedAt:  user.UpdatedAt,
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