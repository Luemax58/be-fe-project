package user

import (
	"context"

	"github.com/Luemax58/be-fe-project/pkg/models"
	"gorm.io/gorm"
)

type IUserRepository interface {
	GetUserByUsername(ctx context.Context, username string) (*models.User, error)
	GetUserByID(ctx context.Context, id uint) (*models.User, error)
}

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) IUserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) GetUserByUsername(ctx context.Context, username string) (*models.User, error) {
	var user models.User

	// ✅ เพิ่ม .Preload("Room") ตรงนี้ครับ
	if err := r.db.WithContext(ctx).
		Preload("Room"). // <--- สำคัญมาก! สั่งให้ดึงข้อมูลห้องมาด้วย
		Where("username = ?", username).
		First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) GetUserByID(ctx context.Context, id uint) (*models.User, error) {
	var user models.User
	// เพิ่ม Preload("Room") เพื่อให้มันไปดึงข้อมูลห้องของคนนี้มาด้วย
	if err := r.db.WithContext(ctx).Preload("Room").Where("user_id = ?", id).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}
