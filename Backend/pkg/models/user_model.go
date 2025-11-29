package models

// User maps to the `users` table
type User struct {
	UserID       uint    `gorm:"primaryKey;autoIncrement"`
	Username     string  `gorm:"type:varchar(50);unique;not null"`
	PasswordHash string  `gorm:"column:password_hash;type:varchar(255);not null"` // เก็บ Hash ยาวๆ ตรงนี้
	FullName     string  `gorm:"type:varchar(100);not null"`
	Phone        *string `gorm:"type:varchar(15)"`
	Role         string  `gorm:"type:enum('owner','tenant');default:'tenant';not null"`

	// --- Relationships ---
	Room                *Room                `gorm:"foreignKey:TenantID"`
	Leases              []Lease              `gorm:"foreignKey:TenantID"`
	Payments            []Payment            `gorm:"foreignKey:TenantID"`
	MaintenanceRequests []MaintenanceRequest `gorm:"foreignKey:TenantID"`
	Announcements       []Announcement       `gorm:"foreignKey:UserID"`
}

// TableName explicitly tells GORM the table name
func (User) TableName() string {
	return "users"
}
