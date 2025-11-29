package main

import (
	"fmt"
	"log"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"

	// Import แพ็กเกจภายในให้ครบ
	"github.com/Luemax58/be-fe-project/internal/billing"
	"github.com/Luemax58/be-fe-project/internal/health"
	"github.com/Luemax58/be-fe-project/internal/maintenance"
	"github.com/Luemax58/be-fe-project/internal/middleware"
	"github.com/Luemax58/be-fe-project/internal/room"
	"github.com/Luemax58/be-fe-project/internal/user"
	"github.com/Luemax58/be-fe-project/pkg/database"
	"golang.org/x/crypto/bcrypt"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("⚠️  Warning: .env file not found (using system env or defaults)")
	} else {
		log.Println("✅ .env loaded successfully")
	}
	// 1. เชื่อมต่อ Database
	db, err := database.ConnectDB()
	password := "123456" // รหัสที่เราต้องการ

	// Gen Hash
	bytes, _ := bcrypt.GenerateFromPassword([]byte(password), 14)

	fmt.Println("--- Copy Hash ด้านล่างนี้ไปใส่ใน Database ---")
	fmt.Println(string(bytes))
	fmt.Println("-------------------------------------------")
	fmt.Println(string(bytes))
	if err != nil {
		log.Fatal(err)
	}

	sqlDB, _ := db.DB()
	defer sqlDB.Close()

	// 2. Setup Server
	r := gin.Default()
	r.Use(cors.Default())

	// ---------------------------------------------------------
	// 3. Wiring (ต่อสายไฟ Dependency Injection)
	// ---------------------------------------------------------

	// --- Health (แบบเดิม) ---
	healthHandler := health.NewHealthHandler(db)

	// --- User (แบบใหม่: Repo -> Service -> Handler) ---
	userRepo := user.NewUserRepository(db)
	userService := user.NewUserService(userRepo)
	userHandler := user.NewUserHandler(userService)

	// --- Room (แบบใหม่: Repo -> Service -> Handler) ---
	roomRepo := room.NewRoomRepository(db)
	roomService := room.NewRoomService(roomRepo)
	roomHandler := room.NewRoomHandler(roomService)

	// --- Maintenance (แบบใหม่: Repo -> Service -> Handler) --- ✅ จุดที่เพิ่งแก้
	maintRepo := maintenance.NewMaintenanceRepository(db)
	maintService := maintenance.NewMaintenanceService(maintRepo)
	maintHandler := maintenance.NewMaintenanceHandler(maintService)

	// --- Billing (แบบเดิม: ต่อ DB ตรง) ---
	billingHandler := billing.NewBillingHandler(db)

	// ---------------------------------------------------------
	// 4. Routes
	// ---------------------------------------------------------
	r.Use(middleware.TimeoutMiddleware(10 * time.Second))

	// Health Check
	r.GET("/health", healthHandler.HealthCheck)

	apiV1 := r.Group("/api/v1")
	{
		// --- Public Routes ---
		// (Register ลบออกแล้วตามแผน)
		apiV1.POST("/login", userHandler.Login)

		// --- Protected Routes ---
		protected := apiV1.Group("")
		protected.Use(middleware.AuthMiddleware())
		{
			// User
			protected.GET("/users/me", userHandler.GetMyProfile)

			// Room
			protected.GET("/rooms", roomHandler.GetAllRooms)
			// protected.POST("/rooms", roomHandler.CreateRoom)

			// Maintenance
			maint := protected.Group("/maintenance")
			{
				maint.POST("/creates", maintHandler.CreateMaintenanceRequest)
				maint.GET("/requests", maintHandler.ListMaintenanceRequests)
			}

			// Billing
			bill := protected.Group("/billing")
			{
				bill.POST("/invoices/generate", billingHandler.GenerateInvoices)
				bill.POST("/utilities/record", billingHandler.RecordUtilityUsage)
				bill.POST("/payments/record", billingHandler.RecordPayment)
			}
		}
	}

	// 5. Run Server
	log.Println("Server running on port :8080")
	r.Run(":8080")
}
