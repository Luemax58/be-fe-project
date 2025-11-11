package user

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// --- DTOs (Data Transfer Objects) ---

// 1. RegisterRequest
type RegisterRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
	FullName string `json:"full_name" binding:"required"`
	Phone    string `json:"phone"`
	Role     string `json:"role" binding:"required"`
}

// 2. RegisterResponse
type RegisterResponse struct {
	UserID   uint   `json:"user_id"`
	Username string `json:"username"`
	FullName string `json:"full_name"`
	Role     string `json:"role"`
}

// 3. LoginRequest (สำหรับ Login)
type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// 4. LoginResponse (สำหรับ Login)
type LoginResponse struct {
	Token string `json:"token"`
}


// --- Handler ---

// IUserHandler
type IUserHandler interface {
	Register(c *gin.Context)
	Login(c *gin.Context)
	GetMyProfile(c *gin.Context)
}

// userHandler
type userHandler struct {
	userService IUserService // Handler จะสั่ง Service
}

// NewUserHandler
func NewUserHandler(service IUserService) IUserHandler {
	return &userHandler{userService: service}
}

// --- Logic (Register) ---

func (h *userHandler) Register(c *gin.Context) {
	// 1. Bind JSON
	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 2. สั่ง Service
	// (ตอนนี้ `newUser` จะเป็น `*models.User` ที่ Go รู้จักแล้ว)
	newUser, err := h.userService.Register(
		req.Username,
		req.Password,
		req.FullName,
		req.Phone,
		req.Role,
	)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// 3. สร้าง Response (ตอนนี้ Go จะรู้จัก newUser.UserID แล้ว)
	response := RegisterResponse{
		UserID:   newUser.UserID,
		Username: newUser.Username,
		FullName: newUser.FullName,
		Role:     newUser.Role,
	}

	// 4. ส่งกลับ
	c.JSON(http.StatusCreated, response)
}

// --- Logic (Login) ---

func (h *userHandler) Login(c *gin.Context) {
	// 1. Bind JSON
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 2. สั่ง Service
	tokenString, err := h.userService.Login(req.Username, req.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	// 3. ส่ง Token กลับ
	c.JSON(http.StatusOK, LoginResponse{Token: tokenString})
}

// GetMyProfile คือ Handler ที่ต้อง "ผ่านด่าน" มาก่อน
func (h *userHandler) GetMyProfile(c *gin.Context) {

    // 1. "ดึง" user_id ที่ "ด่านตรวจ" (Middleware) แปะมาให้
    // (เราไม่ต้องเช็ก Token เองแล้ว เพราะ Middleware ทำให้แล้ว!)
    userID, exists := c.Get("user_id")
    if !exists {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "User ID not found in context"})
        return
    }

    // 2. สั่ง Service ให้ไปหา User
    user, err := h.userService.GetUserProfile(userID.(uint))
    if err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
        return
    }

    // 3. ส่งข้อมูลกลับ (ใช้ RegisterResponse ซ้ำได้เลย เพราะมันไม่มีรหัสผ่าน)
    response := RegisterResponse{
        UserID:   user.UserID,
        Username: user.Username,
        FullName: user.FullName,
        Role:     user.Role,
    }

    c.JSON(http.StatusOK, response)
}