package handlers

import (
	"net/http"

	"devops1/internal/models"
	"github.com/gin-gonic/gin"
)

func (h *Handler) ListMappings(c *gin.Context) {
	m, err := h.Q.ListMappings(c)
	if err != nil {
		errResp(c, 500, "INTERNAL_ERROR", "failed to list mappings")
		return
	}
	ok(c, m)
}

func (h *Handler) CreateMapping(c *gin.Context) {
	var m models.Mapping
	if err := c.ShouldBindJSON(&m); err != nil {
		errResp(c, 400, "VALIDATION_ERROR", "invalid payload")
		return
	}
	m.IsActive = true
	if err := h.Q.CreateMapping(c, m); err != nil {
		errResp(c, 409, "CONFLICT", "mapping already exists")
		return
	}
	c.JSON(http.StatusCreated, gin.H{"message": "created"})
}

func (h *Handler) DeleteMapping(c *gin.Context) {
	if err := h.Q.DeleteMapping(c, c.Param("id")); err != nil {
		errResp(c, 404, "NOT_FOUND", "mapping not found")
		return
	}
	ok(c, gin.H{"message": "deleted"})
}
