package user

import (
	"context"

	"github.com/Luemax58/be-fe-project/pkg/models"
)

// IUserService Interface (ต้องมี GetUserProfile ตรงนี้!)
type IUserService interface {
	Login(ctx context.Context, username, password string) (*models.User, string, error)

	// 👇 บรรทัดนี้สำคัญมากครับ ต้องมีเพื่อให้ Handler เรียกใช้ได้
	GetUserProfile(ctx context.Context, userID uint) (*models.User, error)
}

type userService struct {
	repo IUserRepository
}

func NewUserService(repo IUserRepository) IUserService {
	return &userService{repo: repo}
}

// ---------------------------------------------------------
// Business Logic
// ---------------------------------------------------------

// Login Logic ... (เหมือนเดิม)
func (s *userService) Login(ctx context.Context, username, password string) (*models.User, string, error) {
	// ... (โค้ด Login เดิม) ...
	// ...
	// เพื่อความกระชับ ผมละส่วน Login ไว้ (ใช้โค้ดเดิมได้เลย)
	// ...
	return nil, "", nil // (placeholder)
}

// 👇 Implement ฟังก์ชัน GetUserProfile ให้ตรงกับ Interface
func (s *userService) GetUserProfile(ctx context.Context, userID uint) (*models.User, error) {
	// เรียกไปที่ Repo (ซึ่ง Repo ชื่อ function คือ GetUserByID)
	return s.repo.GetUserByID(ctx, userID)
}

// ... (ฟังก์ชัน generateJWT เหมือนเดิม) ...
