package health

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// (เราไม่จำเป็นต้องมี Service/Repo Layer สำหรับการ Ping ง่ายๆ นี้)

// IHealthHandler คือ "เมนู" ของ Health Handler
type IHealthHandler interface {
	HealthCheck(c *gin.Context)
}

// healthHandler จะ "ถือ" GORM DB object โดยตรง
type healthHandler struct {
	db *gorm.DB
}

// NewHealthHandler คือ "โรงงาน" สร้าง Health Handler
func NewHealthHandler(db *gorm.DB) IHealthHandler {
	return &healthHandler{db: db}
}

// HealthCheck คือ Logic ที่จะถูกเรียกโดย API /health
// (นี่คือฟังก์ชันที่ตรงกับตัวอย่าง "HealthCheck" ของอาจารย์คุณเป๊ะๆ)
func (h *healthHandler) HealthCheck(c *gin.Context) {
	
	// 1. ดึง sql.DB (ตัวเชื่อมต่อดิบ) ออกมาจาก GORM
	sqlDB, err := h.db.DB()
	if err != nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{
			"status": "unhealthy",
			"reason": "Failed to get sql.DB from gorm",
		})
		return
	}

	// 2. สร้าง "ตัวจับเวลา" (Context) 2 วินาที
	// (เราใช้ c.Request.Context() เพื่อให้มัน "ต่อ" มาจาก Middleware หลัก)
	ctx, cancel := context.WithTimeout(c.Request.Context(), 2*time.Second)
	defer cancel()

	// 3. สั่ง Ping (แบบมีตัวจับเวลา)
	err = sqlDB.PingContext(ctx)
	if err != nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{
			"status": "unhealthy",
			"reason": "Database connection failed",
			"error":  err.Error(),
		})
		return
	}

	// 4. ถ้า Ping ผ่าน!
	c.JSON(http.StatusOK, gin.H{"status": "healthy"})
}