package repository

import (
	"github.com/aungmyozaw92/go-grpc-starter/internal/entity"
	"gorm.io/gorm"
)

type UserRepository interface {
	Create(user *entity.User) error
	FindByUsername(username string) (*entity.User, error)
	FindByID(id int) (*entity.User, error)
	Update(user *entity.User) error
	Delete(id int) error
}

type userRepository struct {
	db *gorm.DB
}