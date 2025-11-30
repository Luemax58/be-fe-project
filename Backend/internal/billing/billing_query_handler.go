package billing

import (
    "net/http"
    "github.com/gin-gonic/gin"
    "gorm.io/gorm"
)

type BillingQueryHandler struct {
    DB *gorm.DB
}

func NewBillingQueryHandler(db *gorm.DB) *BillingQueryHandler {
    return &BillingQueryHandler{DB: db}
}

// ดึงบิลทั้งหมด (Admin)
func (h *BillingQueryHandler) GetAllInvoices(c *gin.Context) {
    var invoices []struct {
        BillingID        uint    `json:"billing_id"`
        BillingMonth     string  `json:"billing_month"`
        RoomNumber       string  `json:"room_number"`
        FullName         string  `json:"full_name"`
        TotalUtilityBill float64 `json:"total_utility_bill"`
        Status           string  `json:"status"`
    }

    err := h.DB.Table("monthly_billing mb").
        Select(`
            mb.billing_id,
            DATE_FORMAT(mb.billing_month, '%Y-%m-%d') AS billing_month,
            r.room_number,
            u.full_name,
            mb.total_utility_bill,
            mb.status
        `).
        Joins("LEFT JOIN rooms r ON mb.room_id = r.room_id").
        Joins("LEFT JOIN users u ON r.tenant_id = u.user_id").
        Scan(&invoices).Error

    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusOK, invoices)
}

// ดึงบิลเฉพาะผู้เช่า (Tenant)
func (h *BillingQueryHandler) GetInvoicesByTenant(c *gin.Context) {
    tenantID := c.Param("tenant_id")

    var invoices []struct {
        BillingID        uint      `json:"billing_id"`
        BillingMonth     string    `json:"billing_month"`
        RoomNumber       string    `json:"room_number"`
        TotalUtilityBill float64   `json:"total_utility_bill"`
        Status           string    `json:"status"`
    }

    err := h.DB.Table("monthly_billing mb").
        Select(`
            mb.billing_id,
            DATE_FORMAT(mb.billing_month, '%Y-%m-%d') AS billing_month,
            r.room_number,
            mb.total_utility_bill,
            mb.status
        `).
        Joins("LEFT JOIN rooms r ON mb.room_id = r.room_id").
        Joins("LEFT JOIN leases l ON r.room_id = l.room_id").
        Where("l.tenant_id = ?", tenantID).
        Scan(&invoices).Error

    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusOK, invoices)
}


