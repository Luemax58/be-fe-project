package database

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// ConnectDB เป็นฟังก์ชันที่จะอ่านค่าจาก .env และเชื่อมต่อ Database
// มันจะคืน *gorm.DB object ที่เชื่อมต่อสำเร็จ
func ConnectDB() (*gorm.DB, error) {
	// 1. โหลด .env file (ถ้ามี)
	// เมื่อเรารัน 'go run cmd/api/main.go' จาก root, มันจะหา .env ที่ root
	if err := godotenv.Load(); err != nil {
		log.Println("Warning: .env file not found, using environment variables")
	}

	// 2. ดึงค่าจาก Environment Variables
	dbHost := os.Getenv("DB_HOST")
	dbUser := os.Getenv("DB_USER")
	dbPass := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")
	dbPort := os.Getenv("DB_PORT")

	// 3. (Optional) ถ้าค่าไหนว่าง ให้ตั้ง Default (สำหรับ MySQL)
	if dbHost == "" {
		dbHost = "127.0.0.1"
	}
	if dbPort == "" {
		dbPort = "3306"
	}
	if dbUser == "" {
		dbUser = "root"
	}

	// 4. สร้าง DSN (Data Source Name)
	// format: [user]:[pass]@tcp([host]:[port])/[dbname]?charset=utf8mb4&parseTime=True&loc=Local
	// parseTime=True สำคัญมาก! เพื่อให้ GORM แปลงเวลาใน DB เป็น time.Time ของ Go ได้
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		dbUser,
		dbPass,
		dbHost,
		dbPort,
		dbName,
	)

	// 5. เชื่อมต่อ GORM
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		// ถ้าเชื่อมต่อไม่สำเร็จ ให้ return error
		return nil, fmt.Errorf("failed to connect database: %w", err)
	}

	log.Println("Database connection successful!")
	// คืนค่า DB object ที่เชื่อมต่อสำเร็จ
	return db, nil
}