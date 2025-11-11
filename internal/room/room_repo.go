package room

import (
	"context" // <--- (1) เพิ่ม
	"github.com/Luemax58/be-fe-project/pkg/models"
	"gorm.io/gorm"
)

// IRoomRepository (อัปเกรด "เมนู")
type IRoomRepository interface {
	GetAllRooms(ctx context.Context) ([]models.Room, error) // <--- (2) เพิ่ม ctx
	// TODO: GetRoomByID(ctx context.Context, id uint) (*models.Room, error)
	// TODO: CreateRoom(ctx context.Context, room *models.Room) error
}

// ----------------------------------------------------

// roomRepository (เหมือนเดิม)
type roomRepository struct {
	db *gorm.DB
}

// NewRoomRepository (เหมือนเดิม)
func NewRoomRepository(db *gorm.DB) IRoomRepository {
	return &roomRepository{db: db}
}

// ----------------------------------------------------

// GetAllRooms (อัปเกรด Logic)
func (r *roomRepository) GetAllRooms(ctx context.Context) ([]models.Room, error) { // <--- (3) เพิ่ม ctx
	var rooms []models.Room

	// (4) เพิ่ม .WithContext(ctx) เพื่อส่ง "ตัวจับเวลา" ให้ GORM
	if err := r.db.WithContext(ctx).Preload("Tenant").Find(&rooms).Error; err != nil {
		return nil, err
	}

	return rooms, nil
}