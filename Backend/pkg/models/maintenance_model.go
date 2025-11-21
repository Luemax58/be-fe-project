package models

import "time"

// MaintenanceRequest maps to the `maintenance_requests` table
type MaintenanceRequest struct {
    RequestID        uint      `gorm:"primaryKey;autoIncrement"`
    RoomID           uint      `gorm:"not null"`
    TenantID         uint      `gorm:"not null"`
    IssueDescription string    `gorm:"type:text;not null"`
    RequestDate      time.Time `gorm:"not null;default:CURRENT_TIMESTAMP"`
    Status           string    `gorm:"type:enum('pending','in_progress','completed');default:'pending';not null"`
    RepairCost       float64   `gorm:"type:decimal(10,2);default:0.00"`

    // --- Relationships (GORM) ---
    Room   Room `gorm:"foreignKey:RoomID"`
    Tenant User `gorm:"foreignKey:TenantID"`
}

// TableName explicitly tells GORM the table name
func (MaintenanceRequest) TableName() string {
    return "maintenance_requests"
}