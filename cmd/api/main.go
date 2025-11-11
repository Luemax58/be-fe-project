package main

import (
	"log"

	"github.com/gin-gonic/gin"

	// (อย่าลืมเปลี่ยน [username]/[repo-name] เป็นของคุณ)
	"github.com/Luemax58/be-fe-project/pkg/database"

	// Import "user" (ที่จะดึง Repo, Service, Handler)
	"github.com/Luemax58/be-fe-project/internal/user"
)

func main() {
	// 1. เชื่อมต่อ Database (จาก pkg/database/database.go)
	db, err := database.ConnectDB()
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// (ไม่จำเป็นต้อง AutoMigrate แล้ว เพราะเรา Import .sql ที่สมบูรณ์แบบไปแล้ว)
	// db.AutoMigrate(&models.User{}, &models.Room{}, ...)

	// 2. สร้าง Server (Gin)
	r := gin.Default()

	// 3. VVVV "เดินสายไฟ" (Dependency Injection) VVVV
	// นี่คือหัวใจที่เชื่อมทุก Layer เข้าด้วยกัน

	// --- (ส่วนของคุณ A: User) ---
	userRepo := user.NewUserRepository(db)
	userService := user.NewUserService(userRepo)
	userHandler := user.NewUserHandler(userService)

	// --- (ส่วนของคุณ A: Room) ---
	// TODO: สร้าง roomRepo, roomService, roomHandler ที่นี่

	// --- (ส่วนของคน B: Booking) ---
	// TODO: เพื่อนคุณจะมาสร้าง bookingRepo, bookingService, bookingHandler ที่นี่

	// 4. ตั้งค่า API Routes
	api := r.Group("/api/v1")
	{
		// API ของคุณ A (User)
		api.POST("/register", userHandler.Register)
		// TODO: api.POST("/login", userHandler.Login)
		// TODO: api.GET("/users", ...)

		// API ของคุณ A (Room)
		// TODO: api.GET("/rooms", ...)

		// API ของคน B (Booking)
		// TODO: api.POST("/bookings", ...)
	}

	// 5. รัน Server
	log.Println("Starting server on :8080...")
	if err := r.Run(":8080"); err != nil {
		log.Fatal("Failed to start server: ", err)
	}
}
