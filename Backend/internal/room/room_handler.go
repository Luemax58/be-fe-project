package room

import (
	"net/http"
	// (ไม่ต้อง import context ที่นี่)

	_ "github.com/Luemax58/be-fe-project/pkg/models"
	"github.com/gin-gonic/gin"
)

// --- DTOs (เหมือนเดิม) ---
type RoomTenantResponse struct {
	UserID   uint   `json:"user_id"`
	FullName string `json:"full_name"`
}
type RoomResponse struct {
	RoomID       uint                `json:"room_id"`
	RoomNumber   string              `json:"room_number"`
	Floor        *int                `json:"floor"`
	Status       string              `json:"status"`
	BaseRent     float64             `json:"base_rent"`
	FurnitureFee float64             `json:"furniture_fee"`
	Tenant       *RoomTenantResponse `json:"tenant"`
}

// --- Handler (เหมือนเดิม) ---
type IRoomHandler interface {
	GetAllRooms(c *gin.Context)
}
type roomHandler struct {
	roomService IRoomService
}
func NewRoomHandler(service IRoomService) IRoomHandler {
	return &roomHandler{roomService: service}
}

// --- Logic (อัปเกรด) ---

func (h *roomHandler) GetAllRooms(c *gin.Context) {
	
	// (1) ดึง Context (ตัวจับเวลา) จาก Gin...
	// ...แล้ว "ส่ง" ต่อไปให้ Service
	rooms, err := h.roomService.GetAllRooms(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get rooms"})
		return
	}

	// (ส่วนที่เหลือเหมือนเดิมเป๊ะ)
	var roomResponses []RoomResponse
	for _, room := range rooms {
		var tenantResponse *RoomTenantResponse
		if room.Tenant != nil {
			tenantResponse = &RoomTenantResponse{
				UserID:   room.Tenant.UserID,
				FullName: room.Tenant.FullName,
			}
		}
		roomResponses = append(roomResponses, RoomResponse{
			RoomID:       room.RoomID,
			RoomNumber:   room.RoomNumber,
			Floor:        room.Floor,
			Status:       room.Status,
			BaseRent:     room.BaseRent,
			FurnitureFee: room.FurnitureFee,
			Tenant:       tenantResponse,
		})
	}

	c.JSON(http.StatusOK, roomResponses)
}