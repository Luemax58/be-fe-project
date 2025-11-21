package models

import "time"

// Announcement maps to the `announcements` table
type Announcement struct {
    AnnouncementID uint      `gorm:"primaryKey;autoIncrement"`
    UserID         uint      `gorm:"not null"` // Refers to the 'owner'
    Title          string    `gorm:"type:varchar(255);not null"`
    Content        string    `gorm:"type:text;not null"`
    CreatedAt      time.Time `gorm:"not null;default:CURRENT_TIMESTAMP"`

    // --- Relationships (GORM) ---
    User User `gorm:"foreignKey:UserID"`
}

// TableName explicitly tells GORM the table name
func (Announcement) TableName() string {
    return "announcements"
}