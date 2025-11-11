package database

import (
	"context" // <--- (1) เพิ่ม
	"fmt"
	"log"
	"os"
	"time" // <--- (2) เพิ่ม

	"github.com/joho/godotenv"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// ConnectDB (เวอร์ชันอัปเกรด)
func ConnectDB() (*gorm.DB, error) {
	// 1. โหลด .env
	if err := godotenv.Load(); err != nil {
		log.Println("Warning: .env file not found, using environment variables")
	}

	// 2. ดึงค่าจาก Environment Variables (เหมือนเดิม)
	dbHost := os.Getenv("DB_HOST")
	dbUser := os.Getenv("DB_USER")
	dbPass := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")
	dbPort := os.Getenv("DB_PORT")

	// 3. สร้าง DSN (เหมือนเดิม)
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		dbUser, dbPass, dbHost, dbPort, dbName,
	)

	// 4. เชื่อมต่อ GORM (เหมือนเดิม)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("failed to connect database: %w", err)
	}

	// 5. VVVV ส่วนที่อัปเกรด (ตามสไลด์อาจารย์) VVVV
	// ดึง "sql.DB" (ตัวจริง) ออกมาจาก GORM
	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to get underlying sql.DB: %w", err)
	}

	// (A) ตั้งค่า Connection Pool (เหมือน db.SetMaxOpenConns)
	sqlDB.SetMaxOpenConns(100) // (ตั้งค่าตามความเหมาะสม)
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetConnMaxLifetime(5 * time.Minute)

	// (B) ทดสอบการเชื่อมต่อด้วย Context Timeout (เหมือน db.PingContext)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := sqlDB.PingContext(ctx); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}
	// ^^^^ จบส่วนอัปเกรด ^^^^

	log.Println("Database connection successful with pooling!")
	return db, nil
}