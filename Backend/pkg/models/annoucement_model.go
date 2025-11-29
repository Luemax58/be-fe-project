package models

import (
	"time"
)

// Announcement struct
type Announcement struct {
	AnnouncementID uint   `gorm:"primaryKey;column:announcement_id"`
	UserID         uint   `gorm:"column:user_id;not null"`
	Title          string `gorm:"type:varchar(255);not null"`
	Content        string `gorm:"type:text;not null"`

	// ✅ แก้บรรทัดนี้ครับ (ใส่ precision (3) ให้ครบทั้ง Type และ Default)
	CreatedAt time.Time `gorm:"type:datetime(3);default:CURRENT_TIMESTAMP(3)"`

	// Relationships
	User User `gorm:"foreignKey:UserID"`
}

// TableName
func (Announcement) TableName() string {
	return "announcements"
}
