package models

import "time"

type MaintenanceRequest struct {
    RequestID        uint      `json:"request_id" gorm:"primaryKey;autoIncrement"`
    RoomID           uint      `json:"room_id" gorm:"not null"`
    TenantID         uint      `json:"tenant_id" gorm:"not null"`
    IssueDescription string    `json:"issue_description" gorm:"type:text;not null"`
    RequestDate      time.Time `json:"request_date" gorm:"not null;default:CURRENT_TIMESTAMP"`
    Status           string    `json:"status" gorm:"type:enum('pending','in_progress','completed');default:'pending';not null"`
    RepairCost       float64   `json:"repair_cost" gorm:"type:decimal(10,2);default:0.00"`
}

func (MaintenanceRequest) TableName() string {
    return "maintenance_requests"
}
