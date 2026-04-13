package handlers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

func (h *Handler) Login(c *gin.Context) {
	var req struct{ Username, Password string }
	if err := c.ShouldBindJSON(&req); err != nil {
		errResp(c, http.StatusBadRequest, "VALIDATION_ERROR", "invalid payload")
		return
	}
	u, err := h.Q.AuthPassword(c, req.Username)
	if err != nil || bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(req.Password)) != nil {
		errResp(c, http.StatusUnauthorized, "UNAUTHORIZED", "invalid credentials")
		return
	}
	access, _ := tokenFor(*u, h.Priv, 15*time.Minute)
	refresh, _ := tokenFor(*u, h.Priv, 7*24*time.Hour)
	_, _ = h.Q.Pool.Exec(c, `INSERT INTO sessions(id,user_id,refresh_token,expires_at) VALUES($1,$2,$3,$4)`, uuid.NewString(), u.ID, refresh, time.Now().Add(7*24*time.Hour))
	c.SetCookie("refresh_token", refresh, 604800, "/", "", false, true)
	c.JSON(http.StatusOK, gin.H{"access_token": access, "refresh_token": refresh})
}

func (h *Handler) RefreshToken(c *gin.Context) {
	cookie, _ := c.Cookie("refresh_token")
	if cookie == "" {
		errResp(c, http.StatusUnauthorized, "UNAUTHORIZED", "missing refresh token")
		return
	}
	var uid string
	err := h.Q.Pool.QueryRow(c, `SELECT user_id FROM sessions WHERE refresh_token=$1 AND expires_at>NOW()`, cookie).Scan(&uid)
	if err != nil {
		errResp(c, http.StatusUnauthorized, "UNAUTHORIZED", "invalid refresh token")
		return
	}
	u, err := h.Q.GetUser(c, uid)
	if err != nil {
		errResp(c, http.StatusUnauthorized, "UNAUTHORIZED", "invalid refresh token")
		return
	}
	access, _ := tokenFor(*u, h.Priv, 15*time.Minute)
	c.JSON(http.StatusOK, gin.H{"access_token": access})
}
