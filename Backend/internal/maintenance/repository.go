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
}

type maintenanceRepository struct {
	db *gorm.DB
}

func NewMaintenanceRepository(db *gorm.DB) IMaintenanceRepository {
	return &maintenanceRepository{db: db}
}

func (r *maintenanceRepository) Create(req *models.MaintenanceRequest) error {
	return r.db.Create(req).Error
}

func (r *maintenanceRepository) FindAll(filters map[string]interface{}) ([]models.MaintenanceRequest, error) {
	var requests []models.MaintenanceRequest
	db := r.db.Model(&models.MaintenanceRequest{}).Preload("Room").Preload("Tenant")

	if val, ok := filters["room_id"]; ok && val != "" {
		db = db.Where("room_id = ?", val)
	}
	if val, ok := filters["tenant_id"]; ok && val != "" {
		db = db.Where("tenant_id = ?", val)
	}
	if val, ok := filters["status"]; ok && val != "" {
		db = db.Where("status = ?", val)
	}
	// เพิ่ม Filter วันที่ได้ตามต้องการ

	db = db.Order("request_date DESC")
	err := db.Find(&requests).Error
	return requests, err
}

func (r *maintenanceRepository) FindActiveLease(roomID uint) (*models.Lease, error) {
	var lease models.Lease
	err := r.db.Where("room_id = ? AND status = 'active'", roomID).First(&lease).Error
	return &lease, err
}

func (r *maintenanceRepository) CountPendingRequests(roomID uint) (int64, error) {
	var count int64
	r.db.Model(&models.MaintenanceRequest{}).
		Where("room_id = ? AND status = 'pending'", roomID).
		Count(&count)
	return count, nil
}
