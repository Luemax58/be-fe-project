package main

import (
	"log"
	"time" // (1) เพิ่ม

	"github.com/gin-gonic/gin"

	// (2) VVVV Import "บ้าน" ใหม่ของเราทั้งหมด VVVV
	"github.com/Luemax58/be-fe-project/internal/billing"
	"github.com/Luemax58/be-fe-project/internal/health"
	"github.com/Luemax58/be-fe-project/internal/maintenance"
	"github.com/Luemax58/be-fe-project/internal/middleware"
	"github.com/Luemax58/be-fe-project/internal/room"
	"github.com/Luemax58/be-fe-project/internal/user"
	"github.com/Luemax58/be-fe-project/pkg/database"
)

func main() {
	// 1. เชื่อมต่อ Database (เวอร์ชันอัปเกรด: พร้อม Pooling)
	db, err := database.ConnectDB()
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// (ตาม Best Practice ของอาจารย์: เราต้อง Close() เมื่อแอปปิด)
	sqlDB, err := db.DB()
	if err != nil {
		log.Fatalf("Failed to get underlying sql.DB: %v", err)
	}
	defer sqlDB.Close()

	// 2. สร้าง Server (Gin)
	// (ตาม Best Practice ของอาจารย์: ใช้ ReleaseMode ใน Production)
	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()

	// 3. VVVV "เดินสายไฟ" (Dependency Injection) VVVV
	// นี่คือ "หัวใจ" ที่ตรงกับหลักการ DIP (BookStore) ของอาจารย์เป๊ะๆ

	// --- (ส่วนของ Health Check) ---
	// (Health Handler จะ "ถือ" DB ตรงๆ เพื่อ Ping)
	healthHandler := health.NewHealthHandler(db)

	// --- (ส่วนของ User) ---
	userRepo := user.NewUserRepository(db)
	userService := user.NewUserService(userRepo)
	userHandler := user.NewUserHandler(userService)

	// --- (ส่วนของ Room) ---
	roomRepo := room.NewRoomRepository(db)
	roomService := room.NewRoomService(roomRepo)
	roomHandler := room.NewRoomHandler(roomService)

	// --- (ส่วนของ billing) ---
	billingHandler := billing.NewBillingHandler(db)

	// --- (ส่วนของ maintenance) ---
	maintenanceHandler := maintenance.NewMaintenanceHandler(db)

	// 4. ตั้งค่า Middleware (Global)
	// (ตาม Best Practice ของอาจารย์: ใช้ Timeout Middleware)
	r.Use(middleware.TimeoutMiddleware(10 * time.Second))

	// 5. ตั้งค่า API Routes (Public)
	// (ตาม Best Practice ของอาจารย์: /health อยู่นอกสุด)
	r.GET("/health", healthHandler.HealthCheck)

	apiV1 := r.Group("/api/v1")
	{
		// --- Group 1: Public Routes (ไม่ต้องใช้ Token) ---
		public := apiV1.Group("")
		{
			public.POST("/register", userHandler.Register)
			public.POST("/login", userHandler.Login)
		}

		// --- Group 2: Protected Routes (ต้องใช้ Token) ---
		protected := apiV1.Group("")
		protected.Use(middleware.AuthMiddleware()) // <--- ใช้ "ด่านตรวจ"
		{
			// User Routes
			protected.GET("/users/me", userHandler.GetMyProfile)

			// Room Routes (ที่เราเพิ่งทำเสร็จ!)
			protected.GET("/rooms", roomHandler.GetAllRooms)
			// TODO: protected.POST("/rooms", roomHandler.CreateRoom)

			billing := protected.Group("/billing")
			{
				billing.POST("/invoices/generate", billingHandler.GenerateInvoices)
				billing.POST("/utilities/record", billingHandler.RecordUtilityUsage)
				billing.POST("/payments/record", billingHandler.RecordPayment)
				// billing.GET("/history", handler.GetBillingHistory)
			}

			maint := protected.Group("/maintenance")
			{
				maint.POST("/creates", maintenanceHandler.CreateMaintenanceRequest)
				// maint.PUT("/update/:id", handler.UpdateMaintenanceStatus)
				maint.GET("/requests", maintenanceHandler.ListMaintenanceRequests)
			}

		}
	}

	// 6. รัน Server
	log.Println("Starting server on :8080...")
	if err := r.Run(":8080"); err != nil {
		log.Fatal("Failed to start server: ", err)
	}
}
