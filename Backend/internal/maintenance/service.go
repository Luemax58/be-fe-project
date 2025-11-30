package maintenance

import (
	"errors"
	"time"

	"github.com/Luemax58/be-fe-project/pkg/models"
)

type IMaintenanceService interface {
	CreateRequest(roomID, tenantID uint, description string) (*models.MaintenanceRequest, error)
	ListRequests(filters map[string]interface{}) ([]models.MaintenanceRequest, error)
	UpdateStatus(id string, status string) error
}

type maintenanceService struct {
	repo IMaintenanceRepository
}

func NewMaintenanceService(repo IMaintenanceRepository) IMaintenanceService {
	return &maintenanceService{repo: repo}
}

func (s *maintenanceService) CreateRequest(roomID, tenantID uint, description string) (*models.MaintenanceRequest, error) {

	lease, err := s.repo.FindActiveLease(roomID)
	if err != nil {
		return nil, errors.New("ไม่มี active lease เจอใน room นี้")
	}

	if lease.TenantID != tenantID {
		return nil, errors.New("tenant ไม่ได้เป็นของ room นี้")
	}

	count, _ := s.repo.CountPendingRequests(roomID)
	if count > 0 {
		return nil, errors.New("มี pending request สำหรับ room นี้แล้ว")
	}

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

func (s *maintenanceService) UpdateStatus(id string, status string) error {
	return s.repo.UpdateStatus(id, status)
}
