package user

import (
	"context"

	"github.com/Luemax58/be-fe-project/pkg/models"
	"gorm.io/gorm"
)

// IUserRepository interface
type IUserRepository interface {
	GetUserByUsername(ctx context.Context, username string) (*models.User, error)
	GetUserByID(ctx context.Context, id uint) (*models.User, error)
	// CreateUser เอาออกก็ได้ถ้าไม่ใช้ หรือเก็บไว้เป็น Admin Tool
}

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) IUserRepository {
	return &userRepository{db: db}
}

// 1. GetUserByUsername
func (r *userRepository) GetUserByUsername(ctx context.Context, username string) (*models.User, error) {
	var user models.User
	if err := r.db.WithContext(ctx).Where("username = ?", username).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

// 2. GetUserByID
func (r *userRepository) GetUserByID(ctx context.Context, id uint) (*models.User, error) {
	var user models.User
	if err := r.db.WithContext(ctx).Where("user_id = ?", id).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}
