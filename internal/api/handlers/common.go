package handlers

import (
	"crypto/rsa"
	"net/http"
	"time"

	"devops1/internal/config"
	"devops1/internal/db"
	"devops1/internal/models"
	syncpkg "devops1/internal/sync"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/rs/zerolog"
)

type Handler struct {
	Cfg    *config.Config
	Q      *db.Queries
	Engine *syncpkg.Engine
	Logger zerolog.Logger
	Pub    *rsa.PublicKey
	Priv   *rsa.PrivateKey
}

func New(cfg *config.Config, q *db.Queries, engine *syncpkg.Engine, logger zerolog.Logger, pub *rsa.PublicKey, priv *rsa.PrivateKey) *Handler {
	return &Handler{Cfg: cfg, Q: q, Engine: engine, Logger: logger, Pub: pub, Priv: priv}
}

func errResp(c *gin.Context, status int, code, msg string) {
	c.JSON(status, gin.H{"error": msg, "code": code})
}

func tokenFor(user models.User, priv *rsa.PrivateKey, ttl time.Duration) (string, error) {
	now := time.Now()
	claims := jwt.MapClaims{"sub": user.ID, "role": user.Role, "iat": now.Unix(), "exp": now.Add(ttl).Unix()}
	return jwt.NewWithClaims(jwt.SigningMethodRS256, claims).SignedString(priv)
}

func ok(c *gin.Context, data any) { c.JSON(http.StatusOK, data) }
