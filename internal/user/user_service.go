package user

import (
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
	Register(username, password, fullName, phone, role string) (*models.User, error)
	Login(string, string) (string, error)
	GetUserProfile(id uint) (*models.User, error)
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

func (s *userService) Register(username, password, fullName, phone, role string) (*models.User, error) {
	// 1. ตรวจสอบ Username
	_, err := s.userRepo.GetUserByUsername(username)
	if err == nil {
		return nil, fmt.Errorf("username '%s' already exists", username)
	}
	if err != gorm.ErrRecordNotFound {
		return nil, fmt.Errorf("database error: %w", err)
	}

	// 2. Hashing
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	// 3. เตรียม User
	newUser := &models.User{
		Username:     username,
		PasswordHash: string(hashedPassword),
		FullName:     fullName,
		Phone:        &phone,
		Role:         role,
	}

	// 4. บันทึก
	if err := s.userRepo.CreateUser(newUser); err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}
	return newUser, nil
}


// ----------------------------------------------------
// VVVV นี่คือฟังก์ชัน "Login" ที่ขาดหายไป (Error ที่ 3) VVVV
// ----------------------------------------------------

func (s *userService) Login(username string, password string) (string, error) {
	
	// 1. หา User ด้วย Username
	user, err := s.userRepo.GetUserByUsername(username)
	if err != nil {
		return "", fmt.Errorf("invalid username or password")
	}

	// 2. เปรียบเทียบรหัสผ่าน
	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password))
	if err != nil {
		return "", fmt.Errorf("invalid username or password")
	}

	// 3. ถ้าทุกอย่างถูกต้อง -> สร้าง Token
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

func (s *userService) GetUserProfile(id uint) (*models.User, error) {
    user, err := s.userRepo.GetUserByID(id)
    if err != nil {
        // ถ้า gorm.ErrRecordNotFound
        return nil, fmt.Errorf("user not found")
    }
    return user, nil
}