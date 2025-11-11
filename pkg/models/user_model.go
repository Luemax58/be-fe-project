package models

// User maps to the `users` table
type User struct {
	UserID       uint    `gorm:"primaryKey;autoIncrement"`
	Username     string  `gorm:"type:varchar(50);unique;not null"`
	PasswordHash string  `gorm:"column:password_hash;type:varchar(255);not null"`
	PasswordSalt string  `gorm:"column:password_salt;type:varchar(255);not null"`
	FullName     string  `gorm:"type:varchar(100);not null"`
	Phone        *string `gorm:"type:varchar(15)"` // ใช้ pointer *string สำหรับ field ที่ nullable
	Role         string  `gorm:"type:enum('owner','tenant');default:'tenant';not null"`

	// --- Relationships (GORM) ---
	// 1. User (tenant) can have one Room (from rooms.tenant_id)
	Room *Room `gorm:"foreignKey:TenantID"`

	// 2. User (tenant) can have many Leases
	Leases []Lease `gorm:"foreignKey:TenantID"`

	// 3. User (tenant) can have many Payments
	Payments []Payment `gorm:"foreignKey:TenantID"`

	// 4. User (tenant) can have many MaintenanceRequests
	MaintenanceRequests []MaintenanceRequest `gorm:"foreignKey:TenantID"`

	// 5. User (owner) can have many Announcements
	Announcements []Announcement `gorm:"foreignKey:UserID"`
}

// TableName explicitly tells GORM the table name
func (User) TableName() string {
	return "users"
}
