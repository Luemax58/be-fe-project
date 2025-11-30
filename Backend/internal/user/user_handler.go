package user

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	service IUserService
}

func NewUserHandler(service IUserService) *UserHandler {
	return &UserHandler{service: service}
}

// Login Handler
func (h *UserHandler) Login(c *gin.Context) {
	var req struct {
		Username string `json:"username" binding:"required"`
		Password string `json:"password" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	user, token, err := h.service.Login(c.Request.Context(), req.Username, req.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid username or password"})
		return
	}

	// ✅ แก้ตรงนี้: เช็กว่า User มีห้องหรือไม่
	var roomID uint
	if user.Room != nil {
		roomID = user.Room.RoomID
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Login successful",
		"token":   token,
		"user": gin.H{
			"user_id":   user.UserID,
			"username":  user.Username,
			"full_name": user.FullName,
			"phone":     user.Phone,
			"role":      user.Role,
			"room_id":   roomID, // ส่ง 0 ถ้าไม่มีห้อง
		},
	})
}

// GetMyProfile Handler
func (h *UserHandler) GetMyProfile(c *gin.Context) {
	// 1. ดึง UserID
	userIDRaw, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}
	userID := userIDRaw.(uint)

	// 2. เรียก Service (✅ ใช้ชื่อ GetUserByID ตามที่มีอยู่จริง)
	user, err := h.service.GetUserByID(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
		return
	}

	// 3. เตรียมข้อมูล Room
	var roomID uint
	var roomNumber string
	if user.Room != nil {
		roomID = user.Room.RoomID
		roomNumber = user.Room.RoomNumber
	}

	// 4. ส่ง JSON
	c.JSON(http.StatusOK, gin.H{
		"user_id":     user.UserID,
		"username":    user.Username,
		"full_name":   user.FullName,
		"phone":       user.Phone,
		"role":        user.Role,
		"room_id":     roomID,
		"room_number": roomNumber,
	})
}
