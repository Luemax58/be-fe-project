package models

// User maps to the `users` table
type User struct {
	UserID       uint    `json:"user_id" gorm:"primaryKey;autoIncrement"`
	Username     string  `json:"username" gorm:"type:varchar(50);unique;not null"`
	PasswordHash string  `json:"password_hash" gorm:"column:password_hash;type:varchar(255);not null"`
	FullName     string  `json:"full_name" gorm:"type:varchar(100);not null"`
	Phone        *string `json:"phone" gorm:"type:varchar(15)"`

	Role   string `json:"role" gorm:"type:enum('owner','tenant');default:'tenant';not null"`
	RoomID *uint  `json:"room_id" gorm:"column:room_id"` // สำคัญมาก ใช้ใน frontend หลายหน้า

	// ความสัมพันธ์ (ถ้าต้องการ)
	// Room                *Room                `json:"-" gorm:"foreignKey:RoomID"` 
	// Leases              []Lease              `json:"-" gorm:"foreignKey:TenantID"`
	// Payments            []Payment            `json:"-" gorm:"foreignKey:TenantID"`
	// MaintenanceRequests []MaintenanceRequest `json:"-" gorm:"foreignKey:TenantID"`
	// Announcements       []Announcement       `json:"-" gorm:"foreignKey:UserID"`
}

// TableName explicitly tells GORM the table name
func (User) TableName() string {
	return "users"
}
