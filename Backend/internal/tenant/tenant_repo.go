package tenant

import (
    "github.com/Luemax58/be-fe-project/pkg/models"
    "gorm.io/gorm"
)

type TenantRepository struct {
    db *gorm.DB
}

func NewTenantRepository(db *gorm.DB) *TenantRepository {
    return &TenantRepository{db}
}

// ดึง tenant ทั้งหมด
func (r *TenantRepository) GetAllTenants() ([]models.User, error) {
    var tenants []models.User
    err := r.db.Where("role = ?", "tenant").Find(&tenants).Error
    return tenants, err
}

// ดึงข้อมูล tenant รายบุคคล
func (r *TenantRepository) GetTenantByID(id uint) (*models.User, error) {
    var tenant models.User
    err := r.db.Where("user_id = ? AND role = ?", id, "tenant").First(&tenant).Error
    return &tenant, err
}
