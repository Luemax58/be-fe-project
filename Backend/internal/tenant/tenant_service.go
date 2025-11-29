package tenant

import (
    "github.com/Luemax58/be-fe-project/pkg/models"
)

type TenantService struct {
    repo *TenantRepository
}

func NewTenantService(repo *TenantRepository) *TenantService {
    return &TenantService{repo}
}

func (s *TenantService) GetAllTenants() ([]models.User, error) {
    return s.repo.GetAllTenants()
}

func (s *TenantService) GetTenantByID(id uint) (*models.User, error) {
    return s.repo.GetTenantByID(id)
}
