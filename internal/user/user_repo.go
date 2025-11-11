package user

import (
	// ต้อง import "models" ที่เราสร้างไว้
	"github.com/Luemax58/be-fe-project/pkg/models"

	"gorm.io/gorm"
)

// IUserRepository คือ Interface ที่บอกว่า Repo นี้ทำอะไรได้บ้าง
// (คุณ A ทำเรื่อง User/Room, คุณ B ทำเรื่อง Billing/Booking)
// นี่คือ "สัญญา" ที่ Service (Logic) จะมาเรียกใช้
type IUserRepository interface {
	// 1. (สำหรับ Register/Login)
	GetUserByUsername(username string) (*models.User, error)
	// 2. (สำหรับ Register)
	CreateUser(user *models.User) error

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
func (r *userRepository) GetUserByUsername(username string) (*models.User, error) {
	var user models.User

	// GORM: "SELECT * FROM users WHERE username = ? LIMIT 1"
	// First() จะคืน error 'record not found' ถ้าไม่เจอ
	if err := r.db.Where("username = ?", username).First(&user).Error; err != nil {
		return nil, err // คืนค่า user เปล่า และ error
	}

	return &user, nil // คืนค่า user ที่เจอ และ nil error
}

// 2. CreateUser
func (r *userRepository) CreateUser(user *models.User) error {
	// GORM: "INSERT INTO users (...) VALUES (...)"
	// Create() จะจัดการเรื่องนี้ทั้งหมด และจะคืน error ถ้า insert ไม่ได้
	if err := r.db.Create(&user).Error; err != nil {
		return err
	}

	return nil // สำเร็จ
}
