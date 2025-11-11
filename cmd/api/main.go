package main

import (
    "log"
    "github.com/gin-gonic/gin"
    "github.com/Luemax58/be-fe-project/pkg/database"
    "github.com/Luemax58/be-fe-project/internal/user"
    "github.com/Luemax58/be-fe-project/internal/middleware"
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

	// -- Group 1: Public Routes (ไม่ต้องใช้ Token) --
	public := r.Group("/api/v1")
	{
		public.POST("/register", userHandler.Register)
		public.POST("/login", userHandler.Login)
	}

	// -- Group 2: Protected Routes (ต้องใช้ Token) --
	protected := r.Group("/api/v1")

	// VVVV สั่งให้ Group นี้ "ทั้งหมด" ต้องผ่าน "ด่านตรวจ" ก่อน VVVV
	protected.Use(middleware.AuthMiddleware()) 
	{
		// /api/v1/users/me
		protected.GET("/users/me", userHandler.GetMyProfile)

		// TODO: (ของคุณ A)
		// protected.GET("/rooms", roomHandler.GetAllRooms)
		// protected.POST("/rooms", roomHandler.CreateRoom)

		// TODO: (ของคน B)
		// protected.POST("/bookings", bookingHandler.CreateBooking)
	}

	// 5. รัน Server
	log.Println("Starting server on :8080...")
	if err := r.Run(":8080"); err != nil {
		log.Fatal("Failed to start server: ", err)
	}
}
