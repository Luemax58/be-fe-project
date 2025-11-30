package billing

import (
	"errors"
	"net/http"
	"time"

	"github.com/Luemax58/be-fe-project/pkg/models"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type AdminBillingHandler struct {
	DB *gorm.DB
}

func NewAdminBillingHandler(db *gorm.DB) *AdminBillingHandler {
	return &AdminBillingHandler{DB: db}
}

// --------------------------------------------------------------------
// 1) สร้างใบแจ้งหนี้รายเดือนสำหรับทุกห้องที่มี active lease
// --------------------------------------------------------------------
func (h *AdminBillingHandler) GenerateInvoices(c *gin.Context) {
	var req struct {
		Month        string `json:"month"`          // YYYY-MM
		DueDaysAfter int    `json:"due_days_after"` // default 7 วัน
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// parse เดือน
	var billingMonth time.Time
	var err error
	if req.Month == "" {
		now := time.Now()
		billingMonth = time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, time.Local)
	} else {
		billingMonth, err = time.Parse("2006-01", req.Month)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "รูปแบบเดือนต้องเป็น YYYY-MM"})
			return
		}
	}

	if req.DueDaysAfter == 0 {
		req.DueDaysAfter = 7
	}

	// ค้นหา leases ที่ active ในเดือนนั้น
	var leases []models.Lease
	start := billingMonth
	end := billingMonth.AddDate(0, 1, -1)

	if err := h.DB.Where("start_date <= ? AND end_date >= ? AND status = 'active'", end, start).Find(&leases).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	created := 0
	for _, l := range leases {
		due := billingMonth.AddDate(0, 0, req.DueDaysAfter)

		// check duplicate
		var existing models.MonthlyBilling
		err := h.DB.Where("room_id = ? AND billing_month = ?", l.RoomID, billingMonth).First(&existing).Error
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			continue
		}

		newBill := models.MonthlyBilling{
			RoomID:           l.RoomID,
			BillingMonth:     billingMonth,
			DueDate:          &due,
			Status:           "unpaid",
		}

		if err := h.DB.Create(&newBill).Error; err == nil {
			created++
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "สร้างใบแจ้งหนี้สำเร็จ",
		"created": created,
	})
}

// --------------------------------------------------------------------
// 2) บันทึกค่าน้ำ-ไฟ (admin เท่านั้น)
// --------------------------------------------------------------------
func (h *AdminBillingHandler) RecordUtilityUsage(c *gin.Context) {
	var req struct {
		BillingID        uint     `json:"billing_id"`
		WaterUnits       *float64 `json:"water_units"`
		ElectricityUnits *float64 `json:"electricity_units"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var mb models.MonthlyBilling
	if err := h.DB.First(&mb, req.BillingID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "ไม่พบ billing record"})
		return
	}

	if mb.Status == "paid" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ไม่สามารถแก้ไขบิลที่จ่ายไปแล้ว"})
		return
	}

	// update units
	if req.WaterUnits != nil {
		mb.WaterUnits = *req.WaterUnits
	}
	if req.ElectricityUnits != nil {
		mb.ElectricityUnits = *req.ElectricityUnits
	}

	const waterRate = 18.0
	const electricityRate = 7.0

	mb.WaterBill = mb.WaterUnits * waterRate
	mb.ElectricityBill = mb.ElectricityUnits * electricityRate
	mb.TotalUtilityBill = mb.WaterBill + mb.ElectricityBill

	if err := h.DB.Save(&mb).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "บันทึกค่าน้ำค่าไฟสำเร็จ",
		"billing": mb,
	})
}

// --------------------------------------------------------------------
// 3) Admin บันทึกการชำระเงินแทนผู้เช่า
// --------------------------------------------------------------------
func (h *AdminBillingHandler) RecordPayment(c *gin.Context) {
	var req struct {
		BillingID     uint    `json:"billing_id"`
		TenantID      uint    `json:"tenant_id"`
		AmountPaid    float64 `json:"amount_paid"`
		PaymentMethod string  `json:"payment_method"`
		Notes         *string `json:"notes"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// validate
	if req.AmountPaid <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ยอดชำระต้องมากกว่า 0"})
		return
	}
	if req.PaymentMethod != "cash" && req.PaymentMethod != "transfer" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "payment_method ต้องเป็น cash หรือ transfer"})
		return
	}

	// load billing
	var mb models.MonthlyBilling
	if err := h.DB.First(&mb, req.BillingID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "ไม่พบบิล"})
		return
	}

	if mb.Status == "paid" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "บิลนี้จ่ายไปแล้ว"})
		return
	}

	// load lease เพื่อยืนยัน tenant
	var lease models.Lease
	if err := h.DB.Where("room_id = ? AND status = 'active'", mb.RoomID).First(&lease).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ไม่พบ active lease"})
		return
	}

	if lease.TenantID != req.TenantID {
		c.JSON(http.StatusBadRequest, gin.H{"error": "tenant ไม่ตรงกับห้องนี้"})
		return
	}

	// transaction
	tx := h.DB.Begin()

	payment := models.Payment{
		BillingID:     req.BillingID,
		TenantID:      req.TenantID,
		AmountPaid:    req.AmountPaid,
		PaymentMethod: req.PaymentMethod,
		Notes:         req.Notes,
		PaymentDate:   time.Now(),
	}

	if err := tx.Create(&payment).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// sum payments
	var totalPaid float64
	if err := tx.Model(&models.Payment{}).Where("billing_id = ?", req.BillingID).Select("SUM(amount_paid)").Scan(&totalPaid).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// update status
	if totalPaid >= mb.TotalUtilityBill && mb.TotalUtilityBill > 0 {
		mb.Status = "paid"
		if err := tx.Save(&mb).Error; err != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	}

	tx.Commit()

	c.JSON(http.StatusOK, gin.H{
		"message":       "บันทึกชำระเงินแล้ว",
		"billing_total": mb.TotalUtilityBill,
		"total_paid":    totalPaid,
		"billing_status": mb.Status,
	})
}
