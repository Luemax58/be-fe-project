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
        IssueDescription string `json:"issue_description"`
    }

    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    // ดึง tenant จาก JWT
    tenantID := c.GetUint("user_id")

    // Validate input (ไม่มี TenantID ใน body แล้ว)
    if req.RoomID == 0 || strings.TrimSpace(req.IssueDescription) == "" {
        c.JSON(http.StatusBadRequest, gin.H{"error": "ข้อมูลไม่ครบถ้วน"})
        return
    }

    // ส่งข้อมูลไป service
    mr, err := h.service.CreateRequest(req.RoomID, tenantID, req.IssueDescription)
    if err != nil {
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
    filters := make(map[string]interface{})

    if status := c.Query("status"); status != "" {
        filters["status"] = status
    }

    reqs, err := h.service.ListRequests(filters)
    if err != nil {
        c.JSON(500, gin.H{"error": err.Error()})
        return
    }

    c.JSON(200, gin.H{
        "count": len(reqs),
        "data": reqs,
    })
}


func (h *MaintenanceHandler) UpdateStatus(c *gin.Context) {
	id := c.Param("id")

	var req struct {
		Status string `json:"status"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	if err := h.service.UpdateStatus(id, req.Status); err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{"message": "updated"})
}


