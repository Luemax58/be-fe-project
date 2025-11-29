package user

import (
	"context"
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/Luemax58/be-fe-project/pkg/models"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt" // ⚠️ อย่าลืม go get library นี้นะครับ
)

type IUserService interface {
	Login(ctx context.Context, username, password string) (*models.User, string, error)
	GetUserProfile(ctx context.Context, userID uint) (*models.User, error)
}

type userService struct {
	repo IUserRepository
}

func NewUserService(repo IUserRepository) IUserService {
	return &userService{repo: repo}
}

// Login
func (s *userService) Login(ctx context.Context, username, password string) (*models.User, string, error) {
	// 1. ค้นหา User
	user, err := s.repo.GetUserByUsername(ctx, username)
	if err != nil {
		return nil, "", errors.New("invalid username or password")
	}

	// 2. เช็ก Password
	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password))
	if err != nil {
		return nil, "", errors.New("invalid username or password")
	}

	// 3. ✅ สร้าง JWT Token ของจริง! (ตรงนี้แหละที่แก้)
	tokenString, err := generateJWT(user)
	if err != nil {
		return nil, "", fmt.Errorf("failed to generate token: %v", err)
	}

	return user, tokenString, nil
}

// ฟังก์ชันช่วยสร้าง JWT (Private Function)
func generateJWT(user *models.User) (string, error) {
	// อ่าน Secret Key จาก Env (เหมือนที่ Middleware อ่าน)
	secretKey := os.Getenv("JWT_SECRET_KEY")
	if secretKey == "" {
		// Fallback: ถ้าลืมตั้งใน .env ให้ใช้ค่า Default นี้ไปก่อน (จะได้ไม่พังตอน dev)
		secretKey = "my-secret-key-1234"
		fmt.Println("⚠️ Warning: JWT_SECRET_KEY not set, using default key.")
	}

	// สร้าง Claims (ข้อมูลในบัตร)
	claims := jwt.MapClaims{
		"user_id":  user.UserID,                           // ตรงกับที่ Middleware รอรับ
		"username": user.Username,                         // แถมให้
		"role":     user.Role,                             // แถมให้ (เอาไว้เช็กสิทธิ์)
		"exp":      time.Now().Add(time.Hour * 24).Unix(), // หมดอายุใน 24 ชม.
	}

	// สร้าง Token พร้อมเซ็นลายเซ็น
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(secretKey))

	return tokenString, err
}

// GetUserProfile
func (s *userService) GetUserProfile(ctx context.Context, userID uint) (*models.User, error) {
	return s.repo.GetUserByID(ctx, userID)
}
