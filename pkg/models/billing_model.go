package models

import "time"

// MonthlyBilling maps to the `monthly_billing` table
type MonthlyBilling struct {
    BillingID         uint      `gorm:"primaryKey;autoIncrement"`
    RoomID            uint      `gorm:"not null"`
    BillingMonth      time.Time `gorm:"type:date;not null"`
    DueDate           *time.Time `gorm:"type:date"`
    WaterUnits        float64   `gorm:"type:decimal(10,2)"`
    ElectricityUnits  float64   `gorm:"type:decimal(10,2)"`
    WaterBill         float64   `gorm:"type:decimal(10,2)"`
    ElectricityBill   float64   `gorm:"type:decimal(10,2)"`
    TotalUtilityBill  float64   `gorm:"type:decimal(10,2)"`
    Status            string    `gorm:"type:enum('unpaid','paid','overdue');default:'unpaid';not null"`

    // --- Relationships (GORM) ---
    Room     Room `gorm:"foreignKey:RoomID"`
    Payments []Payment `gorm:"foreignKey:BillingID"`
}

// TableName explicitly tells GORM the table name
func (MonthlyBilling) TableName() string {
    return "monthly_billing"
}

// Payment maps to the `payments` table
type Payment struct {
    PaymentID     uint      `gorm:"primaryKey;autoIncrement"`
    BillingID     uint      `gorm:"not null"`
    TenantID      uint      `gorm:"not null"`
    AmountPaid    float64   `gorm:"type:decimal(10,2);not null"`
    PaymentDate   time.Time `gorm:"not null;default:CURRENT_TIMESTAMP"`
    PaymentMethod string    `gorm:"type:enum('cash','transfer');not null"`
    Notes         *string   `gorm:"type:varchar(255)"`

    // --- Relationships (GORM) ---
    MonthlyBilling MonthlyBilling `gorm:"foreignKey:BillingID"`
    Tenant         User           `gorm:"foreignKey:TenantID"`
}

// TableName explicitly tells GORM the table name
func (Payment) TableName() string {
    return "payments"
}