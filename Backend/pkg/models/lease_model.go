package models

import "time"

// Lease maps to the `leases` table
type Lease struct {
    LeaseID         uint      `gorm:"primaryKey;autoIncrement"`
    RoomID          uint      `gorm:"not null"`
    TenantID        uint      `gorm:"not null"`
    StartDate       time.Time `gorm:"type:date;not null"`
    EndDate         time.Time `gorm:"type:date;not null"`
    SecurityDeposit float64   `gorm:"type:decimal(10,2);not null"`
    Status          string    `gorm:"type:enum('active','expired','terminated');default:'active';not null"`

    // --- Relationships (GORM) ---
    Room   Room `gorm:"foreignKey:RoomID"`
    Tenant User `gorm:"foreignKey:TenantID"`
}

// TableName explicitly tells GORM the table name
func (Lease) TableName() string {
    return "leases"
}