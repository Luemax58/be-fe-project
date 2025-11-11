package user

import (
	// 1. Import library bcrypt
	"fmt"

	"golang.org/x/crypto/bcrypt"

	// 2. Import repo และ models ของเรา
	// (อย่าลืมเปลี่ยน [username]/[repo-name] เป็นของคุณ)
	"github.com/Luemax58/be-fe-project/pkg/models"
	"gorm.io/gorm"
)

// IUserService คือ "เมนู" ที่บอกว่า Service ทำอะไรได้บ้าง
type IUserService interface {
	// รับข้อมูลสำหรับสมัคร (เราจะเพิ่ม Login ที่นี่ทีหลัง)
	Register(username, password, fullName, phone, role string) (*models.User, error)
	// TODO: Login(username, password) (*string, error) // (คืน token)
}

// ----------------------------------------------------

// userService คือ struct "พ่อครัว" ที่ทำงานจริง
type userService struct {
	// มันต้องคุยกับ DB ได้, เลยต้องมี repo
	userRepo IUserRepository
}

// NewUserService คือ "โรงงาน" สร้าง Service
func NewUserService(repo IUserRepository) IUserService {
	return &userService{userRepo: repo}
}

// ----------------------------------------------------
// VVVV นี่คือ "สมอง" ของการ Register VVVV
// ----------------------------------------------------

func (s *userService) Register(username, password, fullName, phone, role string) (*models.User, error) {

	// 1. (Logic) ตรวจสอบว่ามี username นี้ในระบบหรือยัง
	_, err := s.userRepo.GetUserByUsername(username)

	if err == nil {
		// ถ้า err == nil (ไม่มี Error) -> แปลว่า "เจอ"
		return nil, fmt.Errorf("username '%s' already exists", username)
	}

	// ถ้า err != nil (มี Error) ... ให้เช็กก่อนว่าเป็น Error "ที่ไม่เจอ" รึเปล่า
	if err != gorm.ErrRecordNotFound {
		// ถ้าเป็น Error อื่น (เช่น DB ดับ) -> ให้พัง
		return nil, fmt.Errorf("database error: %w", err)
	}

	// ถ้ามาถึงตรงนี้ได้... แปลว่า err คือ gorm.ErrRecordNotFound
	// (ซึ่งแปลว่า Username ว่าง -> ดี!) ... ให้ทำงานต่อได้เลย

	// 2. (Logic) HASH + SALT รหัสผ่าน
	// นี่คือขั้นตอนที่ "ยุ่งยาก" ที่คุณถามถึง... ซึ่งมีแค่บรรทัดเดียว!
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	// 3. (Logic) เตรียมข้อมูล User object เพื่อลง DB
	newUser := &models.User{
		Username:     username,
		PasswordHash: string(hashedPassword), // <--- นี่คือ Hash+Salt
		FullName:     fullName,
		Phone:        &phone, // (ใช้ & เพราะ Phone เป็น *string (nullable))
		Role:         role,
	}

	// (ถ้าคุณเลือก "ลบ" column password_salt, GORM จะไม่ยุ่งกับมันเลย)
	// (ถ้าคุณเลือก "ทำให้ NULL", GORM ก็จะใส่ NULL ให้)

	// 4. (Logic) สั่ง Repo ให้บันทึก
	if err := s.userRepo.CreateUser(newUser); err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	// 5. คืนค่า User ที่เพิ่งสร้างเสร็จ (เผื่อ Handler อยากเอาไปใช้)
	return newUser, nil
}
