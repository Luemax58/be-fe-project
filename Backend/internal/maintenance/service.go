package maintenance

import (
	"errors"
	"time"

	"github.com/Luemax58/be-fe-project/pkg/models"
)

type IMaintenanceService interface {
	CreateRequest(roomID, tenantID uint, description string) (*models.MaintenanceRequest, error)
	ListRequests(filters map[string]interface{}) ([]models.MaintenanceRequest, error)
}

type maintenanceService struct {
	repo IMaintenanceRepository
}

func NewMaintenanceService(repo IMaintenanceRepository) IMaintenanceService {
	return &maintenanceService{repo: repo}
}

func (s *maintenanceService) CreateRequest(roomID, tenantID uint, description string) (*models.MaintenanceRequest, error) {
	// 1. ตรวจสอบ Active Lease
	lease, err := s.repo.FindActiveLease(roomID)
	if err != nil {
		return nil, errors.New("ไม่มี active lease เจอใน room นี้")
	}

	// 2. ตรวจสอบว่าเป็นผู้เช่าคนนี้จริงไหม
	if lease.TenantID != tenantID {
		return nil, errors.New("tenant ไม่ได้เป็นของ room นี้")
	}

	// 3. เช็ค Pending Request
	count, _ := s.repo.CountPendingRequests(roomID)
	if count > 0 {
		return nil, errors.New("มี pending request สำหรับ room นี้แล้ว")
	}

	// 4. สร้าง Request
	mr := models.MaintenanceRequest{
		RoomID:           roomID,
		TenantID:         tenantID,
		IssueDescription: description,
		RequestDate:      time.Now(),
		Status:           "pending",
	}

	if err := s.repo.Create(&mr); err != nil {
		return nil, err
	}

	return &mr, nil
}

func (s *maintenanceService) ListRequests(filters map[string]interface{}) ([]models.MaintenanceRequest, error) {
	return s.repo.FindAll(filters)
}
