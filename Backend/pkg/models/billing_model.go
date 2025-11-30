package models

import "time"

// MonthlyBilling maps to the `monthly_billing` table
type MonthlyBilling struct {
	BillingID        uint       `gorm:"primaryKey;autoIncrement" json:"billing_id"`
	RoomID           uint       `gorm:"not null" json:"room_id"`
	BillingMonth     time.Time  `gorm:"type:date;not null" json:"billing_month"`
	DueDate          *time.Time `gorm:"type:date" json:"due_date"`
	WaterUnits       float64    `gorm:"type:decimal(10,2)" json:"water_units"`
	ElectricityUnits float64    `gorm:"type:decimal(10,2)" json:"electricity_units"`
	WaterBill        float64    `gorm:"type:decimal(10,2)" json:"water_bill"`
	ElectricityBill  float64    `gorm:"type:decimal(10,2)" json:"electricity_bill"`
	TotalUtilityBill float64    `gorm:"type:decimal(10,2)" json:"total_utility_bill"`
	Status           string     `gorm:"type:enum('unpaid','paid','overdue');default:'unpaid';not null" json:"status"`

	// --- Relationships (GORM) ---
	Room     Room      `gorm:"foreignKey:RoomID" json:"room"`
	Payments []Payment `gorm:"foreignKey:BillingID" json:"payments"`
}

// TableName explicitly tells GORM the table name
func (MonthlyBilling) TableName() string {
	return "monthly_billing"
}

// Payment maps to the `payments` table
type Payment struct {
	PaymentID     uint      `gorm:"primaryKey;autoIncrement" json:"payment_id"`
	BillingID     uint      `gorm:"not null" json:"billing_id"`
	TenantID      uint      `gorm:"not null" json:"tenant_id"`
	AmountPaid    float64   `gorm:"type:decimal(10,2);not null" json:"amount_paid"`
	PaymentDate   time.Time `gorm:"not null;default:CURRENT_TIMESTAMP" json:"payment_date"`
	PaymentMethod string    `gorm:"type:enum('cash','transfer');not null" json:"payment_method"`
	Notes         *string   `gorm:"type:varchar(255)" json:"notes"`

	// --- Relationships (GORM) ---
	MonthlyBilling MonthlyBilling `gorm:"foreignKey:BillingID" json:"billing"`
	Tenant         User           `gorm:"foreignKey:TenantID" json:"tenant"`
}

// TableName explicitly tells GORM the table name
func (Payment) TableName() string {
	return "payments"
}
