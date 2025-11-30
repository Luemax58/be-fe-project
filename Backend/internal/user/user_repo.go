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
    if err := r.db.WithContext(ctx).
        Where("username = ?", username).
        First(&user).Error; err != nil {
        return nil, err
    }
    return &user, nil
}

func (r *userRepository) GetUserByID(ctx context.Context, id uint) (*models.User, error) {
    var user models.User
    err := r.db.WithContext(ctx).Where("user_id = ?", id).First(&user).Error
    return &user, err
}

