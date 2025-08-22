package repository

import (
	"github.com/aungmyozaw92/go-grpc-starter/internal/entity"
	"gorm.io/gorm"
)

type UserRepository interface {
	Create(user *entity.User) error
	FindByUsername(username string) (*entity.User, error)
	FindByEmail(email string) (*entity.User, error)
	FindByID(id int) (*entity.User, error)
	Update(user *entity.User) error
	Delete(id int) error
	GetUserList(page, limit int, search string) ([]*entity.User, int64, error)
	ExistsByUsername(username string) (bool, error)
	ExistsByEmail(email string) (bool, error)
	ExistsByUsernameExcludeID(username string, excludeID uint) (bool, error)
	ExistsByEmailExcludeID(email string, excludeID uint) (bool, error)
}

type userRepository struct {
	db *gorm.DB
}
