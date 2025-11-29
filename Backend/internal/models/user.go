package models

type User struct {
	UserID       uint   `json:"user_id" gorm:"primaryKey"`
	Username     string `json:"username"`
	PasswordHash string `json:"password_hash"`
	FullName     string `json:"full_name"`
	Phone        string `json:"phone"`
	Role         string `json:"role"`
}
