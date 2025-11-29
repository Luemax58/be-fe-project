package maintenance

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

type MaintenanceHandler struct {
	service IMaintenanceService
}

func NewMaintenanceHandler(service IMaintenanceService) *MaintenanceHandler {
	return &MaintenanceHandler{service: service}
}

func (h *MaintenanceHandler) CreateMaintenanceRequest(c *gin.Context) {
	var req struct {
		RoomID           uint   `json:"room_id"`
		TenantID         uint   `json:"tenant_id"`
		IssueDescription string `json:"issue_description"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if req.RoomID == 0 || req.TenantID == 0 || strings.TrimSpace(req.IssueDescription) == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ข้อมูลไม่ครบถ้วน"})
		return
	}

	mr, err := h.service.CreateRequest(req.RoomID, req.TenantID, req.IssueDescription)
	if err != nil {
		// แยก Status Code ตาม Error จริงๆ ควรละเอียดกว่านี้แต่เอาคร่าวๆ
		if err.Error() == "มี pending request สำหรับ room นี้แล้ว" {
			c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
		} else {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		}
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "สร้างคำร้องซ่อมบำรุงแล้ว",
		"request": mr,
	})
}

func (h *MaintenanceHandler) ListMaintenanceRequests(c *gin.Context) {
	filters := map[string]interface{}{
		"room_id":   c.Query("room_id"),
		"tenant_id": c.Query("tenant_id"),
		"status":    c.Query("status"),
	}

	requests, err := h.service.ListRequests(filters)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"count": len(requests),
		"data":  requests,
	})
}
