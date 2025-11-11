package user

import (
	"context"
	"github.com/Luemax58/be-fe-project/pkg/models"

	"gorm.io/gorm"
)

// IUserRepository คือ Interface ที่บอกว่า Repo นี้ทำอะไรได้บ้าง
// (คุณ A ทำเรื่อง User/Room, คุณ B ทำเรื่อง Billing/Booking)
// นี่คือ "สัญญา" ที่ Service (Logic) จะมาเรียกใช้
type IUserRepository interface {
	GetUserByUsername(ctx context.Context, username string) (*models.User, error)
    CreateUser(ctx context.Context, user *models.User) error
    GetUserByID(ctx context.Context, id uint) (*models.User, error)
	// TODO (ของคุณ A):
	// GetUserByID(id uint) (*models.User, error)
	// UpdateUser(user *models.User) error
	// DeleteUser(id uint) error
	// GetAllUsers() ([]models.User, error)
}

// ----------------------------------------------------

// userRepository คือ struct ที่ "ทำจริง" ตาม Interface
// มันจะถือ Connection (DB) ไว้
type userRepository struct {
	db *gorm.DB
}

// NewUserRepository คือ "โรงงาน" สร้าง Repository
// (เราจะเรียกใช้มันใน main.go)
func NewUserRepository(db *gorm.DB) IUserRepository {
	// คืนค่า struct ที่มี DB พร้อมใช้งาน
	return &userRepository{db: db}
}

// ----------------------------------------------------
// VVVV นี่คือโค้ด Logic การดึงข้อมูล VVVV
// ----------------------------------------------------

// 1. GetUserByUsername
func (r *userRepository) GetUserByUsername(ctx context.Context, username string) (*models.User, error) {
    var user models.User
    // (3) เพิ่ม .WithContext(ctx) ต่อท้าย .db
    // นี่คือการ "ส่งต่อ" ตัวจับเวลาให้ GORM
    if err := r.db.WithContext(ctx).Where("username = ?", username).First(&user).Error; err != nil {
        return nil, err
    }
    return &user, nil
}

// 2. CreateUser
func (r *userRepository) CreateUser(ctx context.Context, user *models.User) error {
    // (4) เพิ่ม .WithContext(ctx)
    return r.db.WithContext(ctx).Create(user).Error
}

func (r *userRepository) GetUserByID(ctx context.Context, id uint) (*models.User, error) {
    var user models.User
    // (5) เพิ่ม .WithContext(ctx)
    if err := r.db.WithContext(ctx).Where("user_id = ?", id).First(&user).Error; err != nil {
        return nil, err
    }
    return &user, nil
}