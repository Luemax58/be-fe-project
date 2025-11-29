package database

import (
	"fmt"
	"os"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func ConnectDB() (*gorm.DB, error) {
	// อ่านค่าจาก .env
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbName := os.Getenv("DB_NAME")

	// ถ้าไม่มีใน .env ให้ใช้ค่า Default (กันเหนียว)
	if dbUser == "" {
		dbUser = "root"
	}
	if dbPassword == "" {
		dbPassword = "Lue5548"
	} // รหัสของคุณ
	if dbHost == "" {
		dbHost = "localhost"
	}
	if dbPort == "" {
		dbPort = "3306"
	}
	if dbName == "" {
		dbName = "mydatabase"
	}

	// สร้าง DSN (Data Source Name)
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		dbUser, dbPassword, dbHost, dbPort, dbName,
	)

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	// ❌❌❌ คอมเมนต์ส่วนนี้ทิ้งครับ! ❌❌❌
	// เพราะเราใช้ SQL Script สร้างตารางที่ถูกต้องไปแล้ว ไม่ต้องให้ GORM สร้างซ้ำ
	/*
		err = db.AutoMigrate(
			&models.User{},
			&models.Room{},
			&models.Lease{},
			&models.MonthlyBilling{},
			&models.Payment{},
			&models.MaintenanceRequest{},
			&models.Announcement{},
		)
		if err != nil {
			return nil, err
		}
	*/

	return db, nil
}
