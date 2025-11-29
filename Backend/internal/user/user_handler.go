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

	c.JSON(http.StatusOK, gin.H{
		"message": "Login successful",
		"token":   token,
		"user": gin.H{
			"id":       user.UserID,
			"username": user.Username,
			"role":     user.Role,
			"fullName": user.FullName,
		},
	})
}

// GetMyProfile Handler
func (h *UserHandler) GetMyProfile(c *gin.Context) {
	// สมมติว่า Middleware แปะ userID มาให้
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	// แปลง Type (ต้องระวังตรงนี้ เช็กให้ตรงกับ Middleware)
	var uid uint
	if v, ok := userID.(float64); ok { // กรณี JWT parse เป็น float64
		uid = uint(v)
	} else if v, ok := userID.(uint); ok {
		uid = v
	} else {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "User ID type error"})
		return
	}

	user, err := h.service.GetUserProfile(c.Request.Context(), uid)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"user": user})
}
