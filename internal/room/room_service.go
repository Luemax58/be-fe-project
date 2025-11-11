package room

import (
	"context" // <--- (1) เพิ่ม
	"github.com/Luemax58/be-fe-project/pkg/models"
)

// IRoomService (อัปเกรด "เมนู")
type IRoomService interface {
	GetAllRooms(ctx context.Context) ([]models.Room, error) // <--- (2) เพิ่ม ctx
}

// ----------------------------------------------------

// roomService (เหมือนเดิม)
type roomService struct {
	roomRepo IRoomRepository
}

// NewRoomService (เหมือนเดิม)
func NewRoomService(repo IRoomRepository) IRoomService {
	return &roomService{roomRepo: repo}
}

// ----------------------------------------------------

// GetAllRooms (อัปเกรด Logic)
func (s *roomService) GetAllRooms(ctx context.Context) ([]models.Room, error) { // <--- (3) เพิ่ม ctx
	
	// (4) ส่ง ctx "ทะลุ" ต่อไปให้ Repo
	rooms, err := s.roomRepo.GetAllRooms(ctx)
	if err != nil {
		return nil, err
	}

	return rooms, nil
}