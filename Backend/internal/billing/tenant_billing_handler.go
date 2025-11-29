package billing

import (
    "net/http"
    "github.com/gin-gonic/gin"
    "gorm.io/gorm"
    "github.com/Luemax58/be-fe-project/pkg/models"
)

type TenantBillingHandler struct {
    db *gorm.DB
}

func NewTenantBillingHandler(db *gorm.DB) *TenantBillingHandler {
    return &TenantBillingHandler{db}
}

// ✔ ผู้เช่าดูใบแจ้งหนี้ของตนเอง
func (h *TenantBillingHandler) GetMyInvoices(c *gin.Context) {
    tenantID := c.GetUint("user_id")

    var invoices []models.MonthlyBilling
    err := h.db.
        Joins("JOIN leases ON leases.room_id = monthly_billing.room_id").
        Where("leases.tenant_id = ?", tenantID).
        Preload("Room").
        Preload("Payments").
        Find(&invoices).Error

    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusOK, invoices)
}

// ✔ ผู้เช่าดูประวัติการชำระเงิน
func (h *TenantBillingHandler) GetMyPayments(c *gin.Context) {
    tenantID := c.GetUint("user_id")

    var payments []models.Payment
    err := h.db.Where("tenant_id = ?", tenantID).Find(&payments).Error

    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusOK, payments)
}

// ✔ ผู้เช่าบันทึกการชำระเงิน
func (h *TenantBillingHandler) PayInvoice(c *gin.Context) {
    tenantID := c.GetUint("user_id")

    var input struct {
        BillingID  uint    `json:"billing_id"`
        AmountPaid float64 `json:"amount_paid"`
        Method     string  `json:"method"` // cash / transfer
        Notes      string  `json:"notes"`
    }
    if err := c.ShouldBindJSON(&input); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    payment := models.Payment{
        BillingID:     input.BillingID,
        TenantID:      tenantID,
        AmountPaid:    input.AmountPaid,
        PaymentMethod: input.Method,
        Notes:         &input.Notes,
    }

    if err := h.db.Create(&payment).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    // อัปเดตสถานะ billing
    h.db.Model(&models.MonthlyBilling{}).
        Where("billing_id = ?", input.BillingID).
        Update("status", "paid")

    c.JSON(http.StatusOK, gin.H{
        "message": "payment recorded",
        "payment": payment,
    })
}
