package middleware

import (
	"context"
	"time"

	"github.com/gin-gonic/gin"
)

// TimeoutMiddleware คือ "ด่านตรวจ" ที่สร้างตัวจับเวลา
// (นี่คือโค้ดจาก "main.go" ในตัวอย่างของอาจารย์คุณเลยครับ)
func TimeoutMiddleware(timeout time.Duration) gin.HandlerFunc {
	return func(c *gin.Context) {
		
		// 1. สร้าง Context (ตัวจับเวลา) ใหม่
		ctx, cancel := context.WithTimeout(c.Request.Context(), timeout)
		defer cancel()

		// 2. "ยัด" Context ใหม่นี้กลับเข้าไปใน Request
		// (นี่คือที่มาของ c.Request.Context() ที่เราใช้ใน Handler ครับ)
		c.Request = c.Request.WithContext(ctx)

		// 3. ปล่อยให้ Request วิ่งไปที่ Handler
		c.Next()
	}
}