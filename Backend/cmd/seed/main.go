package main

import (
	"log"

	"github.com/Luemax58/be-fe-project/pkg/database"
	"golang.org/x/crypto/bcrypt"
)

func main() {
	// 1. เชื่อมต่อ Database
	db, err := database.ConnectDB()
	if err != nil {
		log.Fatalf("เชื่อมต่อ Database ไม่สำเร็จ: %v", err)
	}

	// 2. สร้าง Hash ของรหัส "123456" (ของจริง)
	log.Println("⏳ กำลังสร้าง Hash สำหรับรหัสผ่าน 123456 ...")
	hashedBytes, err := bcrypt.GenerateFromPassword([]byte("123456"), 14)
	if err != nil {
		log.Fatal("Hash Password พัง:", err)
	}
	realHash := string(hashedBytes)

	// 3. อัปเดต "ทุกคน" (ทั้ง Owner และ Tenant) ให้ใช้รหัสนี้
	// ใช้ Exec SQL ตรงๆ เพื่อความรวดเร็วและชัวร์ที่สุด
	result := db.Exec("UPDATE users SET password_hash = ?", realHash)

	if result.Error != nil {
		log.Fatal("❌ อัปเดตไม่สำเร็จ:", result.Error)
	}

	log.Printf("✅ สำเร็จ! รีเซ็ตรหัสผ่านของ User ทั้งหมด %v คน เป็น '123456' เรียบร้อย", result.RowsAffected)
}
