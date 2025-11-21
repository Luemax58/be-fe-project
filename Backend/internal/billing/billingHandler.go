package billing

import (
	"errors"
	"net/http"
	"time"

	"github.com/Luemax58/be-fe-project/pkg/models"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// Handler struct
type billHandler struct {
	DB *gorm.DB
}

// Constructor สำหรับสร้าง Handler
func NewBillingHandler(db *gorm.DB) *billHandler {
	return &billHandler{DB: db}
}

// ---------------- Billing APIs ------------------

// GenerateInvoices: สร้างบิลรายเดือนให้กับทุกห้องที่มีสัญญาเช่า
func (b *billHandler) GenerateInvoices(c *gin.Context) {
	var req struct {
		Month        string `json:"month"`
		DueDaysAfter int    `json:"due_days_after"`
	}
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	//ตรวจสอบเดือน ถ้าไม่กรอกจะใช้เดือนปัจจุบัน
	var billingMonth time.Time
	var err error
	if req.Month == "" {
		now := time.Now()
		billingMonth = time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, time.Local)
	} else {
		billingMonth, err = time.Parse("2006-01", req.Month)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "เดือนต้องเป็น YYYY-MM"})
			return
		}
	}

	if req.DueDaysAfter == 0 {
		req.DueDaysAfter = 7
	}

	//ค้นหา Lease ที่ Active อยู่ในเดือนนั้น
	var leases []models.Lease
	start := billingMonth
	end := billingMonth.AddDate(0, 1, -1)
	if err := b.DB.Where("start_date <= ? AND end_date >= ? AND status = ?", end, start, "active").Find(&leases).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if len(leases) == 0 {
		c.JSON(http.StatusOK, gin.H{"message": "ไม่ active leases ที่เจอ"})
		return
	}

	//วนสร้างใบแจ้งหนี้ให้ทุก lease
	created := 0
	for _, l := range leases {
		due := billingMonth.AddDate(0, 0, req.DueDaysAfter)
		mb := models.MonthlyBilling{
			RoomID:           l.RoomID,
			BillingMonth:     billingMonth,
			DueDate:          &due,
			WaterUnits:       0,
			ElectricityUnits: 0,
			Status:           "unpaid",
		}
		//เช็คก่อนว่ามีใบแจ้งหนี้ของห้องนี้ ประจำเดือนนี้แล้วหรือยัง
		var existing models.MonthlyBilling
		err := b.DB.Where("room_id = ? AND billing_month = ?", mb.RoomID, mb.BillingMonth).First(&existing).Error

		if errors.Is(err, gorm.ErrRecordNotFound) {
			if err := b.DB.Create(&mb).Error; err == nil {
				created++
			}
		} else if err != nil {
			continue // หรือ log error
		}
	}

	c.JSON(http.StatusOK, gin.H{"message": "สร้างบิลรายเดือนสำเร็จแล้ว", "created": created})
}

// RecordUtilityUsage:ฟังก์ชันนี้ใช้สำหรับบันทึกค่าน้ำและค่าไฟของบิลรายเดือน
func (b *billHandler) RecordUtilityUsage(c *gin.Context) {
	var req struct {
		BillingID        uint     `json:"billing_id"`
		WaterUnits       *float64 `json:"water_units"`
		ElectricityUnits *float64 `json:"electricity_units"`
	}

	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// --- ตรวจสอบ billing ---
	var mb models.MonthlyBilling
	if err := b.DB.First(&mb, req.BillingID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "ไม่พบ billing record"})
		return
	}

	// --- ป้องกันการแก้ไขบิลที่จ่ายเงินแล้ว ---
	if mb.Status == "paid" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ไม่สามารถแก้ไขบิลที่จ่ายเงินแล้ว"})
		return
	}

	// --- Validate units ---
	if req.WaterUnits != nil && *req.WaterUnits < 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "water_units ไม่สามารถติดลบได้"})
		return
	}
	if req.ElectricityUnits != nil && *req.ElectricityUnits < 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "electricity_units ไม่สามารถติดลบได้"})
		return
	}

	// --- อัปเดตเฉพาะค่าที่ส่งเข้ามา ---
	if req.WaterUnits != nil {
		mb.WaterUnits = *req.WaterUnits
	}
	if req.ElectricityUnits != nil {
		mb.ElectricityUnits = *req.ElectricityUnits
	}

	// --- กำหนดราคาไฟ/น้ำที่นี่ ---
	const waterRate = 18.0
	const electricityRate = 7.0

	// --- คำนวณยอดเงิน ---
	mb.WaterBill = mb.WaterUnits * waterRate
	mb.ElectricityBill = mb.ElectricityUnits * electricityRate
	mb.TotalUtilityBill = mb.WaterBill + mb.ElectricityBill

	// --- บันทึกลง DB ---
	if err := b.DB.Save(&mb).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	b.DB.Preload("Room").Preload("Tenant").First(&mb, mb.BillingID)

	// --- ส่งกลับ ---
	c.JSON(http.StatusOK, gin.H{"message": "บันทึกค่าน้ำและค่าไฟสำเร็จแล้ว", "billing": mb})
}

// RecordPayment:ฟังก์ชันนี้ใช้สำหรับบันทึกข้อมูลการชำระเงินของผู้เช่า
func (b *billHandler) RecordPayment(c *gin.Context) {
	var req struct {
		BillingID     uint    `json:"billing_id"`
		TenantID      uint    `json:"tenant_id"`
		AmountPaid    float64 `json:"amount_paid"`
		PaymentMethod string  `json:"payment_method"`
		Notes         *string `json:"notes"`
	}

	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// --- Validate input ---
	if req.AmountPaid <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "amount_paid ต้องมากกว่า 0"})
		return
	}

	if req.PaymentMethod != "cash" && req.PaymentMethod != "transfer" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "payment_method ไม่ถูกต้อง (ต้องเป็น 'cash' หรือ 'transfer')"})
		return
	}

	// --- โหลด billing ---
	var mb models.MonthlyBilling
	if err := b.DB.First(&mb, req.BillingID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "billing ไม่พบ"})
		return
	}

	// --- ห้ามจ่ายเพิ่มหลังบิลเป็น paid ---
	if mb.Status == "paid" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "bill นี้ถูกจ่ายไปเรียบร้อยแล้ว"})
		return
	}

	// --- ตรวจสอบ tenant ให้ตรงกับห้องนั้น ---
	var lease models.Lease
	if err := b.DB.Where("room_id = ? AND status = 'active'", mb.RoomID).
		First(&lease).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "active lease ไม่เจอใน room นี้"})
		return
	}

	if lease.TenantID != req.TenantID {
		c.JSON(http.StatusBadRequest, gin.H{"error": "tenant ไม่ได้เป็นเจ้าของ bill นี้"})
		return
	}

	// --- Begin Transaction ---
	tx := b.DB.Begin()
	if tx.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "เริ่ม transaction ล้มเหลว"})
		return
	}

	// --- Create Payment ---
	p := models.Payment{
		BillingID:     req.BillingID,
		TenantID:      req.TenantID,
		AmountPaid:    req.AmountPaid,
		PaymentDate:   time.Now(),
		PaymentMethod: req.PaymentMethod,
		Notes:         req.Notes,
	}

	if err := tx.Create(&p).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// --- คำนวณยอดชำระรวม ---
	var sum float64
	if err := tx.Model(&models.Payment{}).
		Where("billing_id = ?", req.BillingID).
		Select("SUM(amount_paid)").Scan(&sum).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// --- ถ้าจ่ายครบให้เปลี่ยนสถานะเป็น paid ---
	if sum >= mb.TotalUtilityBill && mb.TotalUtilityBill > 0 {
		mb.Status = "paid"
		if err := tx.Save(&mb).Error; err != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	}

	// --- Commit Transaction ---
	if err := tx.Commit().Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// --- Respond ---
	c.JSON(http.StatusOK, gin.H{
		"message":        "payment recorded",
		"billing_status": mb.Status,
		"total_paid":     sum,
		"billing_total":  mb.TotalUtilityBill,
	})
}
