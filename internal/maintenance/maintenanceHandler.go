package maintenance

import (
	"net/http"
	"strings"
	"time"

	"github.com/Luemax58/be-fe-project/pkg/models"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// Handler struct
type maintHandler struct {
	DB *gorm.DB
}

// Constructor สำหรับสร้าง Handler
func NewMaintenanceHandler(db *gorm.DB) *maintHandler {
	return &maintHandler{DB: db}
}

// ---------------- Maintenance APIs ------------------

// CreateMaintenanceRequest:ฟังก์ชันนี้ใช้สำหรับสร้างคำร้องซ่อมบำรุง
func (m *maintHandler) CreateMaintenanceRequest(c *gin.Context) {
	var req struct {
		RoomID           uint   `json:"room_id"`
		TenantID         uint   `json:"tenant_id"`
		IssueDescription string `json:"issue_description"`
	}

	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// --- Validate input ---
	if req.RoomID == 0 || req.TenantID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "room_id และ tenant_id จำเป็นต้องกรอก"})
		return
	}

	if strings.TrimSpace(req.IssueDescription) == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "issue_description ไม่สามารถเป็นค่าว่าง"})
		return
	}

	// --- ตรวจสอบว่าห้องมีสัญญาเช่าที่ active ไหม ---
	var lease models.Lease
	if err := m.DB.Where("room_id = ? AND status = 'active'", req.RoomID).
		First(&lease).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ไม่มี active lease เจอใน room นี้"})
		return
	}

	// --- ตรวจสอบว่าผู้เช่าตรงกับห้องนั้นจริงไหม ---
	if lease.TenantID != req.TenantID {
		c.JSON(http.StatusBadRequest, gin.H{"error": "tenant ไม่ได้เป็นของ room นี้"})
		return
	}

	// --- เช็คว่ามี request pending อยู่แล้วไหม ---
	var exists int64
	m.DB.Model(&models.MaintenanceRequest{}).
		Where("room_id = ? AND status = 'pending'", req.RoomID).
		Count(&exists)

	if exists > 0 {
		c.JSON(http.StatusConflict, gin.H{"error": "มี pending request สำหรับ room นี้แล้ว"})
		return
	}

	// --- Create new maintenance request ---
	mr := models.MaintenanceRequest{
		RoomID:           req.RoomID,
		TenantID:         req.TenantID,
		IssueDescription: req.IssueDescription,
		RequestDate:      time.Now(),
		Status:           "pending",
	}

	if err := m.DB.Create(&mr).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	m.DB.Preload("Room").Preload("Tenant").First(&mr, mr.RequestID)

	c.JSON(http.StatusCreated, gin.H{
		"message": "สร้างคำร้องซ่อมบำรุงแล้ว",
		"request": mr,
	})
}

// ListMaintenanceRequests: ลิสการซ่อมบำรุงของทุกห้อง
func (m *maintHandler) ListMaintenanceRequests(c *gin.Context) {
	var requests []models.MaintenanceRequest

	// --- Query Parameters ---
	roomID := c.Query("room_id")
	tenantID := c.Query("tenant_id")
	status := c.Query("status")
	fromDate := c.Query("from")
	toDate := c.Query("to")

	// --- Base Query + Preload ---
	db := m.DB.
		Model(&models.MaintenanceRequest{}).
		Preload("Room").
		Preload("Tenant")

	// --- Filters ---
	if roomID != "" {
		db = db.Where("room_id = ?", roomID)
	}

	if tenantID != "" {
		db = db.Where("tenant_id = ?", tenantID)
	}

	if status != "" {
		db = db.Where("status = ?", status)
	}

	if fromDate != "" {
		db = db.Where("request_date >= ?", fromDate)
	}

	if toDate != "" {
		db = db.Where("request_date <= ?", toDate)
	}

	// --- Sorting (latest first) ---
	db = db.Order("request_date DESC")

	// --- Execute query ---
	if err := db.Find(&requests).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"count": len(requests),
		"data":  requests,
	})
}
