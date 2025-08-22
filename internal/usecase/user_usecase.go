package usecase

import (
	"errors"
	"math"

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

func (u *UserUseCase) Register(user *entity.User) (string, error) {
	// Check if username already exists
	exists, err := u.userRepo.ExistsByUsername(user.Username)
	if err != nil {
		return "", err
	}
	if exists {
		return "", errors.New("username already exists")
	}

	// Check if email already exists (if email is provided)
	if user.Email != nil && *user.Email != "" {
		emailExists, err := u.userRepo.ExistsByEmail(*user.Email)
		if err != nil {
			return "", err
		}
		if emailExists {
			return "", errors.New("email already exists")
		}
	}

	hashedPassword, err := infrastructure.HashPassword(user.Password)
	if err != nil {
		return "", err
	}
	user.Password = string(hashedPassword)

	if err := u.userRepo.Create(user); err != nil {
		return "", err
	}

	return infrastructure.GenerateJWT(int(user.ID))
}

func (u *UserUseCase) Login(username, password string) (string, error) {
	user, err := u.userRepo.FindByUsername(username)
	if err != nil {
		return "", err
	}
	if !infrastructure.CheckPasswordHash(user.Password, password) {
		return "", errors.New("invalid credentials")
	}
	return infrastructure.GenerateJWT(int(user.ID))
}

func (u *UserUseCase) GetProfile(token string) (*entity.User, error) {
	userID, err := infrastructure.ValidateToken(token)
	if err != nil {
		return nil, err
	}
	return u.userRepo.FindByID(userID)
}

type UserListResult struct {
	Users      []*entity.User
	Pagination PaginationInfo
}

type PaginationInfo struct {
	CurrentPage int
	PerPage     int
	TotalPages  int
	TotalCount  int64
	HasNext     bool
	HasPrev     bool
}

func (u *UserUseCase) GetUserList(token string, page, limit int, search string) (*UserListResult, error) {
	// Validate token
	_, err := infrastructure.ValidateToken(token)
	if err != nil {
		return nil, err
	}

	// Set default values
	if page <= 0 {
		page = 1
	}
	if limit <= 0 || limit > 100 {
		limit = 10 // Default limit
	}

	// Get users and total count
	users, total, err := u.userRepo.GetUserList(page, limit, search)
	if err != nil {
		return nil, err
	}

	// Calculate pagination info
	totalPages := int(math.Ceil(float64(total) / float64(limit)))
	hasNext := page < totalPages
	hasPrev := page > 1

	pagination := PaginationInfo{
		CurrentPage: page,
		PerPage:     limit,
		TotalPages:  totalPages,
		TotalCount:  total,
		HasNext:     hasNext,
		HasPrev:     hasPrev,
	}

	return &UserListResult{
		Users:      users,
		Pagination: pagination,
	}, nil
}

func (u *UserUseCase) GetUser(token string, userID int) (*entity.User, error) {
	// Validate token
	_, err := infrastructure.ValidateToken(token)
	if err != nil {
		return nil, err
	}

	// Get user by ID
	user, err := u.userRepo.FindByID(userID)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (u *UserUseCase) CreateUser(token string, user *entity.User) (*entity.User, error) {
	// Validate token (only authenticated users can create users)
	_, err := infrastructure.ValidateToken(token)
	if err != nil {
		return nil, err
	}

	// Check if username already exists
	exists, err := u.userRepo.ExistsByUsername(user.Username)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, errors.New("username already exists")
	}

	// Check if email already exists (if email is provided)
	if user.Email != nil && *user.Email != "" {
		emailExists, err := u.userRepo.ExistsByEmail(*user.Email)
		if err != nil {
			return nil, err
		}
		if emailExists {
			return nil, errors.New("email already exists")
		}
	}

	// Hash password
	hashedPassword, err := infrastructure.HashPassword(user.Password)
	if err != nil {
		return nil, err
	}
	user.Password = string(hashedPassword)

	// Create user
	if err := u.userRepo.Create(user); err != nil {
		return nil, err
	}

	return user, nil
}

func (u *UserUseCase) UpdateUser(token string, userID int, updateData *entity.User) (*entity.User, error) {
	// Validate token
	_, err := infrastructure.ValidateToken(token)
	if err != nil {
		return nil, err
	}

	// Get existing user
	existingUser, err := u.userRepo.FindByID(userID)
	if err != nil {
		return nil, err
	}

	// Check if username already exists (excluding current user)
	if updateData.Username != existingUser.Username {
		exists, err := u.userRepo.ExistsByUsernameExcludeID(updateData.Username, existingUser.ID)
		if err != nil {
			return nil, err
		}
		if exists {
			return nil, errors.New("username already exists")
		}
	}

	// Check if email already exists (excluding current user)
	if updateData.Email != nil && *updateData.Email != "" {
		// Only check if email is different from current
		currentEmail := ""
		if existingUser.Email != nil {
			currentEmail = *existingUser.Email
		}

		if *updateData.Email != currentEmail {
			emailExists, err := u.userRepo.ExistsByEmailExcludeID(*updateData.Email, existingUser.ID)
			if err != nil {
				return nil, err
			}
			if emailExists {
				return nil, errors.New("email already exists")
			}
		}
	}

	// Update fields
	existingUser.Username = updateData.Username
	existingUser.Name = updateData.Name
	existingUser.Email = updateData.Email
	existingUser.Phone = updateData.Phone
	existingUser.Mobile = updateData.Mobile
	existingUser.ImageURL = updateData.ImageURL
	existingUser.IsActive = updateData.IsActive
	existingUser.RoleID = updateData.RoleID

	// Update user
	if err := u.userRepo.Update(existingUser); err != nil {
		return nil, err
	}

	return existingUser, nil
}

func (u *UserUseCase) DeleteUser(token string, userID int) error {
	// Validate token
	_, err := infrastructure.ValidateToken(token)
	if err != nil {
		return err
	}

	// Check if user exists
	_, err = u.userRepo.FindByID(userID)
	if err != nil {
		return err
	}

	// Delete user
	return u.userRepo.Delete(userID)
}
