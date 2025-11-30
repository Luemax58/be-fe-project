package billing

import (
	"net/http"

	"github.com/Luemax58/be-fe-project/pkg/models"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type TenantBillingHandler struct {
	DB *gorm.DB
}

func NewTenantBillingHandler(db *gorm.DB) *TenantBillingHandler {
	return &TenantBillingHandler{DB: db}
}

// ผู้เช่าดูบิลของตนเอง
func (h *TenantBillingHandler) GetMyInvoices(c *gin.Context) {
    tenantID := c.GetUint("user_id")

    var invoices []models.MonthlyBilling
    err := h.DB.
        Table("monthly_billing").
        Select(`
            monthly_billing.billing_id,
            monthly_billing.billing_month,
            monthly_billing.total_utility_bill,
            monthly_billing.status,
            monthly_billing.room_id
        `).
        Joins("JOIN leases ON leases.room_id = monthly_billing.room_id").
        Where("leases.tenant_id = ?", tenantID).
        Order("monthly_billing.billing_month DESC").
        Find(&invoices).Error

    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusOK, invoices)
}


// ผู้เช่าดูประวัติการชำระเงินของตนเอง
func (h *TenantBillingHandler) GetMyPayments(c *gin.Context) {
	tenantID := c.GetUint("user_id")

	var payments []models.Payment
	err := h.DB.Where("tenant_id = ?", tenantID).Find(&payments).Error

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, payments)
}

// ผู้เช่าจ่ายบิลของตนเอง
func (h *TenantBillingHandler) PayInvoice(c *gin.Context) {
	tenantID := c.GetUint("user_id")

	var req struct {
		BillingID  uint    `json:"billing_id"`
		AmountPaid float64 `json:"amount_paid"`
		Method     string  `json:"method"`
		Notes      string  `json:"notes"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	p := models.Payment{
		BillingID:     req.BillingID,
		TenantID:      tenantID,
		AmountPaid:    req.AmountPaid,
		PaymentMethod: req.Method,
		Notes:         &req.Notes,
	}

	if err := h.DB.Create(&p).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// update bill status
	h.DB.Model(&models.MonthlyBilling{}).
		Where("billing_id = ?", req.BillingID).
		Update("status", "paid")

	c.JSON(http.StatusOK, gin.H{
		"message": "payment recorded",
		"payment": p,
	})
}
