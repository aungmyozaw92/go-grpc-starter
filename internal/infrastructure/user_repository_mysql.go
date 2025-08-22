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

func (r *UserRepository) FindByEmail(email string) (*entity.User, error) {
	var user entity.User
	err := r.DB.Where("email = ?", email).First(&user).Error
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

func (r *UserRepository) GetUserList(page, limit int, search string) ([]*entity.User, int64, error) {
	var users []*entity.User
	var total int64

	query := r.DB.Model(&entity.User{})

	// Apply search filter if provided
	if search != "" {
		searchPattern := "%" + search + "%"
		query = query.Where("username LIKE ? OR name LIKE ? OR email LIKE ?",
			searchPattern, searchPattern, searchPattern)
	}

	// Get total count
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Apply pagination
	offset := (page - 1) * limit
	if err := query.Offset(offset).Limit(limit).Order("created_at DESC").Find(&users).Error; err != nil {
		return nil, 0, err
	}

	return users, total, nil
}

func (r *UserRepository) ExistsByUsername(username string) (bool, error) {
	var count int64
	err := r.DB.Model(&entity.User{}).Where("username = ?", username).Count(&count).Error
	return count > 0, err
}

func (r *UserRepository) ExistsByEmail(email string) (bool, error) {
	var count int64
	err := r.DB.Model(&entity.User{}).Where("email = ?", email).Count(&count).Error
	return count > 0, err
}

func (r *UserRepository) ExistsByUsernameExcludeID(username string, excludeID uint) (bool, error) {
	var count int64
	err := r.DB.Model(&entity.User{}).Where("username = ? AND id != ?", username, excludeID).Count(&count).Error
	return count > 0, err
}

func (r *UserRepository) ExistsByEmailExcludeID(email string, excludeID uint) (bool, error) {
	var count int64
	err := r.DB.Model(&entity.User{}).Where("email = ? AND id != ?", email, excludeID).Count(&count).Error
	return count > 0, err
}
