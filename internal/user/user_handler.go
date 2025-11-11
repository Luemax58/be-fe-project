package user

import (
	// Import fmt
	"net/http"

	"github.com/gin-gonic/gin"
)

// --- DTOs (Data Transfer Objects) ---
// เราสร้าง Structs แยกสำหรับ "รับ" และ "ส่ง" ข้อมูล

// 1. RegisterRequest คือ JSON ที่เรา "คาดหวัง" จะได้รับตอนสมัคร
// `binding:"required"` คือเวทมนตร์ของ Gin ที่จะเช็กให้ว่า "ห้ามว่าง"
type RegisterRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
	FullName string `json:"full_name" binding:"required"`
	Phone    string `json:"phone"`                   // Phone ไม่บังคับ
	Role     string `json:"role" binding:"required"` // ปกติ Frontend ควรส่ง 'tenant'
}

// 2. RegisterResponse คือ JSON ที่เราจะ "ส่งกลับ"
// (สังเกตว่าเราจะไม่ส่ง PasswordHash กลับไปเด็ดขาด!)
type RegisterResponse struct {
	UserID   uint   `json:"user_id"`
	Username string `json:"username"`
	FullName string `json:"full_name"`
	Role     string `json:"role"`
}

// --- Handler ---

// IUserHandler คือ "เมนู" ของ Handler
type IUserHandler interface {
	Register(c *gin.Context)
	// TODO: Login(c *gin.Context)
}

// userHandler คือ "คนรับออเดอร์" ที่ทำงานจริง
type userHandler struct {
	userService IUserService // Handler จะสั่ง Service
}

// NewUserHandler คือ "โรงงาน" สร้าง Handler
func NewUserHandler(service IUserService) IUserHandler {
	return &userHandler{userService: service}
}

// --- VVVV นี่คือ Logic ของ API VVVV ---

func (h *userHandler) Register(c *gin.Context) {
	// 1. รับ JSON Request และตรวจสอบ (Bind & Validate)
	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		// ถ้า JSON ที่ส่งมาไม่ครบ (เช่น ไม่มี "username")
		// Gin จะโยน Error ให้ตรงนี้เลย
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 2. ถ้า JSON ถูกต้อง, ส่งข้อมูลไปให้ "Service" (สมอง)
	newUser, err := h.userService.Register(
		req.Username,
		req.Password,
		req.FullName,
		req.Phone,
		req.Role,
	)

	if err != nil {
		// ถ้า Service ตอบกลับมาว่ามี Error (เช่น "username already exists")
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// 3. ถ้าสำเร็จ! (Service คืน newUser มาให้)
	// ให้เราแปลง newUser (จาก DB) เป็น Response (ที่จะส่งกลับ)
	response := RegisterResponse{
		UserID:   newUser.UserID,
		Username: newUser.Username,
		FullName: newUser.FullName,
		Role:     newUser.Role,
	}

	// 4. ส่ง JSON กลับไปด้วย Status 201 Created
	c.JSON(http.StatusCreated, response)

} // <--- 🚨 นี่คือ } ที่ผมลืมใส่ให้ครับ!
