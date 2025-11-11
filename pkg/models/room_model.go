package models

// Room maps to the `rooms` table
type Room struct {
    RoomID       uint    `gorm:"primaryKey;autoIncrement"`
    RoomNumber   string  `gorm:"type:varchar(10);unique;not null"`
    Floor        *int    `gorm:"type:int(3)"` // ใช้ pointer *int สำหรับ nullable
    Status       string  `gorm:"type:enum('available','occupied','maintenance');default:'available';not null"`
    TenantID     *uint   `gorm:"column:tenant_id"` // Nullable Foreign Key
    BaseRent     float64 `gorm:"type:decimal(10,2);not null"`
    FurnitureFee float64 `gorm:"type:decimal(10,2);not null"`

    // --- Relationships (GORM) ---
    // 1. Room belongs to one Tenant (User)
    Tenant User `gorm:"foreignKey:TenantID"`
    
    // 2. Room can have many Leases
    Leases []Lease `gorm:"foreignKey:RoomID"`
    
    // 3. Room can have many MonthlyBillings
    MonthlyBillings []MonthlyBilling `gorm:"foreignKey:RoomID"`
    
    // 4. Room can have many MaintenanceRequests
    MaintenanceRequests []MaintenanceRequest `gorm:"foreignKey:RoomID"`
}

// TableName explicitly tells GORM the table name
func (Room) TableName() string {
    return "rooms"
}