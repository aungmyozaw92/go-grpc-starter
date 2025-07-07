package infrastructure

import (
	"github.com/aungmyozaw92/go-grpc-starter/internal/entity"
	"gorm.io/gorm"
)

type UserRepository struct {
	DB *gorm.DB
}

func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{DB: db}
}

func (r *UserRepository) Create(user *entity.User) error {
	return r.DB.Create(user).Error
}

func (r *UserRepository) FindByUsername(username string) (*entity.User, error) {
	var user entity.User
  	err := r.DB.Where("username = ?", username).First(&user).Error
  	return &user, err
}

func (r *UserRepository) FindByID(id int) (*entity.User, error) {
	var user entity.User
	err := r.DB.First(&user, id).Error
	return &user, err
}

func (r *UserRepository) Update(user *entity.User) error {
	return r.DB.Save(user).Error
}

func (r *UserRepository) Delete(id int) error {
	return r.DB.Delete(&entity.User{}, id).Error
}	