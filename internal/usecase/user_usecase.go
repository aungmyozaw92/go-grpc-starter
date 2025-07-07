package usecase

import (
	"errors"

	"github.com/aungmyozaw92/go-grpc-starter/internal/entity"
	"github.com/aungmyozaw92/go-grpc-starter/internal/infrastructure"
	"github.com/aungmyozaw92/go-grpc-starter/internal/repository"
)

type UserUseCase struct {
	userRepo repository.UserRepository
}

func NewUserUseCase(userRepo repository.UserRepository) *UserUseCase {
	return &UserUseCase{userRepo: userRepo}
}

func (u *UserUseCase) Register(user *entity.User) (string, error)  {
	hashedPassword, err := infrastructure.HashPassword(user.Password)
	if err != nil {
		return "", err
	}
	user.Password = string(hashedPassword)
	 if err := u.userRepo.Create(user); err != nil {
		return "", err
	}

	return infrastructure.GenerateJWT(user.ID)
}


func (u *UserUseCase) Login(username, password string) (string, error) {
  user, err := u.userRepo.FindByUsername(username)
  if err != nil {
    return "", err
  }
  if !infrastructure.CheckPasswordHash(user.Password, password) {
    return "", errors.New("invalid credentials")
  }
  return infrastructure.GenerateJWT(user.ID)
}

func (u *UserUseCase) GetProfile(token string) (*entity.User, error) {
  userID, err := infrastructure.ValidateToken(token)
  if err != nil {
    return nil, err
  }
  return u.userRepo.FindByID(userID)
}