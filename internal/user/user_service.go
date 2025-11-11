package user

import (
	"context"
	"fmt"
	"os" // <--- เพิ่ม
	"time" // <--- เพิ่ม

	"golang.org/x/crypto/bcrypt"
	"github.com/golang-jwt/jwt/v5" // <--- เพิ่ม
	"github.com/Luemax58/be-fe-project/pkg/models"
	"gorm.io/gorm"
)

// IUserService คือ "เมนู" ที่บอกว่า Service ทำอะไรได้บ้าง
type IUserService interface {
	Register(ctx context.Context, username, password, fullName, phone, role string) (*models.User, error)
    Login(ctx context.Context, username, password string) (string, error)
    GetUserProfile(ctx context.Context, id uint) (*models.User, error)
}

// ----------------------------------------------------

// userService คือ struct "พ่อครัว" ที่ทำงานจริง
type userService struct {
	userRepo IUserRepository
}

// NewUserService คือ "โรงงาน" สร้าง Service
func NewUserService(repo IUserRepository) IUserService {
	return &userService{userRepo: repo}
}

// ----------------------------------------------------
// VVVV "สมอง" ของการ Register VVVV (อันนี้ของคุณถูกต้องอยู่แล้ว)
// ----------------------------------------------------

func (s *userService) Register(ctx context.Context, username, password, fullName, phone, role string) (*models.User, error) {

    // (3) ส่ง ctx ต่อไปให้ Repo
    _, err := s.userRepo.GetUserByUsername(ctx, username)

    if err == nil {
        return nil, fmt.Errorf("username '%s' already exists", username)
    }
    if err != gorm.ErrRecordNotFound {
        return nil, fmt.Errorf("database error: %w", err)
    }

    // (Hashing ไม่ต้องใช้ ctx เพราะมันทำงานใน CPU, ไม่เกี่ยวกับ DB)
    hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
    if err != nil {
        return nil, fmt.Errorf("failed to hash password: %w", err)
    }

    newUser := &models.User{
        // ... (เหมือนเดิม) ...
        Username:     username,
        PasswordHash: string(hashedPassword),
        FullName:     fullName,
        Phone:        &phone,
        Role:         role,
    }

    // (4) ส่ง ctx ต่อไปให้ Repo
    if err := s.userRepo.CreateUser(ctx, newUser); err != nil {
        return nil, fmt.Errorf("failed to create user: %w", err)
    }

    return newUser, nil
}


// ----------------------------------------------------
// VVVV นี่คือฟังก์ชัน "Login" ที่ขาดหายไป (Error ที่ 3) VVVV
// ----------------------------------------------------

func (s *userService) Login(ctx context.Context, username string, password string) (string, error) {

    // (5) ส่ง ctx ต่อไปให้ Repo
    user, err := s.userRepo.GetUserByUsername(ctx, username)
    if err != nil {
        return "", fmt.Errorf("invalid username or password")
    }

    // (Compare Hashing ไม่ต้องใช้ ctx)
    err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password))
    if err != nil {
        return "", fmt.Errorf("invalid username or password")
    }

    // (Generate JWT ไม่ต้องใช้ ctx)
    tokenString, err := generateJWT(user)
    if err != nil {
        return "", fmt.Errorf("failed to generate token: %w", err)
    }

    return tokenString, nil
}

// ----------------------------------------------------
// VVVV นี่คือฟังก์ชัน "generateJWT" ที่ขาดหายไป VVVV
// ----------------------------------------------------

func generateJWT(user *models.User) (string, error) {
	// 1. ดึง Secret Key จาก .env
	secretKey := os.Getenv("JWT_SECRET_KEY")
	if secretKey == "" {
		return "", fmt.Errorf("JWT_SECRET_KEY is not set")
	}

	// 2. สร้าง "Claims" (ข้อมูลในบัตร)
	claims := jwt.MapClaims{
		"user_id": user.UserID,
		"role":    user.Role,
		"exp":     time.Now().Add(time.Hour * 72).Unix(), // หมดอายุใน 3 วัน
		"iat":     time.Now().Unix(),                      // ออกบัตรเมื่อ
	}

	// 3. สร้าง Token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// 4. "เซ็น" Token ด้วย Secret Key
	tokenString, err := token.SignedString([]byte(secretKey))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func (s *userService) GetUserProfile(ctx context.Context, id uint) (*models.User, error) {
    // (6) ส่ง ctx ต่อไปให้ Repo
    user, err := s.userRepo.GetUserByID(ctx, id)
    if err != nil {
        return nil, fmt.Errorf("user not found")
    }
    return user, nil
}