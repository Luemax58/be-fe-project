package middleware

import (
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

// AuthMiddleware คือ "ด่านตรวจ" ของเรา
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 1. ดึง Header "Authorization"
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header is required"})
			c.Abort() // หยุด! ห้ามไปต่อ
			return
		}

		// 2. เช็กว่ามันขึ้นต้นด้วย "Bearer " หรือไม่
		// (รูปแบบมาตรฐานคือ "Bearer [token...]")
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid authorization header format"})
			c.Abort()
			return
		}

		tokenString := parts[1]

		// 3. ตรวจสอบ Token (ไขกุญแจ)
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			// ดึง Secret Key
			secretKey := os.Getenv("JWT_SECRET_KEY")
			if secretKey == "" {
				return nil, fmt.Errorf("JWT_SECRET_KEY is not set")
			}
			// เช็กว่าใช้วิธีเซ็น (HS256) ตรงกับเราไหม
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return []byte(secretKey), nil
		})

		if err != nil || !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}

		// 4. ถ้า Token ถูกต้อง -> ดึง "Claims" (ข้อมูลในบัตร) ออกมา
		if claims, ok := token.Claims.(jwt.MapClaims); ok {
			// ดึง user_id (ที่เรายัดเข้าไปตอน Login) ออกมา
			// (JWT จะเก็บตัวเลขเป็น float64 เสมอ)
			userIDFloat, ok := claims["user_id"].(float64)
			if !ok {
				c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token claims"})
				c.Abort()
				return
			}

			// "แปะ" user_id ไว้ใน Context ของ Gin
			// เพื่อให้ Handler ที่อยู่ "หลังด่าน" ดึงไปใช้ได้
			c.Set("user_id", uint(userIDFloat))

			// 5. "เชิญครับ!" (ปล่อยให้ไปที่ Handler ตัวจริง)
			c.Next()

		} else {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token claims"})
			c.Abort()
		}
	}
}