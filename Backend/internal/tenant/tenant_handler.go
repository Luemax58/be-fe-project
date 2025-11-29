package tenant

import (
    "net/http"
    "strconv"

    "github.com/gin-gonic/gin"
)

type TenantHandler struct {
    service *TenantService
}

func NewTenantHandler(service *TenantService) *TenantHandler {
    return &TenantHandler{service}
}

// GET /tenants
func (h *TenantHandler) GetAllTenants(c *gin.Context) {
    tenants, err := h.service.GetAllTenants()
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    c.JSON(http.StatusOK, tenants)
}

// GET /tenants/:id
func (h *TenantHandler) GetTenantByID(c *gin.Context) {
    id, _ := strconv.Atoi(c.Param("id"))

    tenant, err := h.service.GetTenantByID(uint(id))
    if err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "tenant not found"})
        return
    }

    c.JSON(http.StatusOK, tenant)
}
