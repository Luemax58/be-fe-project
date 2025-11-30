package maintenance

import (
	"github.com/Luemax58/be-fe-project/pkg/models"
	"gorm.io/gorm"
)

type IMaintenanceRepository interface {
	Create(req *models.MaintenanceRequest) error
	FindAll(filters map[string]interface{}) ([]models.MaintenanceRequest, error)
	FindActiveLease(roomID uint) (*models.Lease, error)
	CountPendingRequests(roomID uint) (int64, error)
	UpdateStatus(id string, status string) error
}

type MaintenanceRepository struct {
	db *gorm.DB
}

func NewMaintenanceRepository(db *gorm.DB) IMaintenanceRepository {
	return &MaintenanceRepository{db: db}
}

func (r *MaintenanceRepository) Create(req *models.MaintenanceRequest) error {
	return r.db.Create(req).Error
}

func (r *MaintenanceRepository) FindAll(filters map[string]interface{}) ([]models.MaintenanceRequest, error) {
	var results []models.MaintenanceRequest

	db := r.db.Model(&models.MaintenanceRequest{})

	if v, ok := filters["status"]; ok {
		db = db.Where("status = ?", v)
	}

	if err := db.Order("request_date DESC").Find(&results).Error; err != nil {
		return nil, err
	}

	return results, nil
}

func (r *MaintenanceRepository) FindActiveLease(roomID uint) (*models.Lease, error) {
	var lease models.Lease
	err := r.db.Where("room_id = ? AND status = ?", roomID, "active").First(&lease).Error
	if err != nil {
		return nil, err
	}
	return &lease, nil
}

func (r *MaintenanceRepository) CountPendingRequests(roomID uint) (int64, error) {
	var count int64
	err := r.db.Model(&models.MaintenanceRequest{}).
		Where("room_id = ? AND status = ?", roomID, "pending").
		Count(&count).Error
	return count, err
}

func (r *MaintenanceRepository) UpdateStatus(id string, status string) error {
	return r.db.Model(&models.MaintenanceRequest{}).
		Where("request_id = ?", id).
		Update("status", status).Error
}
