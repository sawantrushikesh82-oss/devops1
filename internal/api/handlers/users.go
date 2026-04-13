package handlers

import (
	"errors"
	"net/http"

	"devops1/internal/models"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
	"golang.org/x/crypto/bcrypt"
)

func (h *Handler) ListUsers(c *gin.Context) {
	u, err := h.Q.ListUsers(c)
	if err != nil {
		errResp(c, 500, "INTERNAL_ERROR", "failed to list users")
		return
	}
	ok(c, u)
}

func (h *Handler) CreateUser(c *gin.Context) {
	var u models.User
	if err := c.ShouldBindJSON(&u); err != nil {
		errResp(c, 400, "VALIDATION_ERROR", "invalid payload")
		return
	}
	hash, _ := bcrypt.GenerateFromPassword([]byte(u.PasswordHash), bcrypt.DefaultCost)
	u.PasswordHash = string(hash)
	u.IsActive = true
	if err := h.Q.CreateUser(c, u); err != nil {
		errResp(c, 409, "CONFLICT", "username already exists")
		return
	}
	c.JSON(http.StatusCreated, gin.H{"message": "created"})
}

func (h *Handler) GetUser(c *gin.Context) {
	u, err := h.Q.GetUser(c, c.Param("id"))
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			errResp(c, 404, "NOT_FOUND", "user not found")
			return
		}
		errResp(c, 500, "INTERNAL_ERROR", "failed to get user")
		return
	}
	ok(c, u)
}

func (h *Handler) UpdateUser(c *gin.Context) {
	var u models.User
	if err := c.ShouldBindJSON(&u); err != nil {
		errResp(c, 400, "VALIDATION_ERROR", "invalid payload")
		return
	}
	if u.PasswordHash != "" {
		hash, _ := bcrypt.GenerateFromPassword([]byte(u.PasswordHash), bcrypt.DefaultCost)
		u.PasswordHash = string(hash)
	}
	if err := h.Q.UpdateUser(c, c.Param("id"), u); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			errResp(c, 404, "NOT_FOUND", "user not found")
			return
		}
		errResp(c, 500, "INTERNAL_ERROR", "failed to update user")
		return
	}
	ok(c, gin.H{"message": "updated"})
}

func (h *Handler) DeleteUser(c *gin.Context) {
	if err := h.Q.DeleteUser(c, c.Param("id")); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			errResp(c, 404, "NOT_FOUND", "user not found")
			return
		}
		errResp(c, 500, "INTERNAL_ERROR", "failed to delete user")
		return
	}
	ok(c, gin.H{"message": "deleted"})
}
