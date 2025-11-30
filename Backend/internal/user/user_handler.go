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
			"user_id":   user.UserID,
			"username":  user.Username,
			"full_name": user.FullName,
			"phone":     user.Phone,
			"role":      user.Role,
			"room_id":   user.RoomID,
		},
	})
}

// GetMyProfile Handler
func (h *UserHandler) GetMyProfile(c *gin.Context) {
    userID := c.GetUint("user_id")

	var room models.Room
	h.repo.DB.Where("tenant_id = ?", user.UserID).First(&room)

    user, err := h.service.GetUserByID(c.Request.Context(), userID)
    if err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
        return
    }

    c.JSON(http.StatusOK, gin.H{
        "user_id":   user.UserID,
        "username":  user.Username,
        "full_name": user.FullName,
        "phone":     user.Phone,
        "role":      user.Role,
        "room_id":   user.RoomID,
    })
}
