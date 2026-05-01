package api

import (
	"crypto/rsa"

	"devops1/internal/api/handlers"
	"devops1/internal/api/middleware"
	"devops1/internal/config"
	"devops1/internal/db"
	syncpkg "devops1/internal/sync"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
)

func NewServer(cfg *config.Config, q *db.Queries, engine *syncpkg.Engine, logger zerolog.Logger, pubKey *rsa.PublicKey, privKey *rsa.PrivateKey) *gin.Engine {
	r := gin.New()
	r.Use(middleware.RequestLogger(logger), gin.Recovery())
	h := handlers.New(cfg, q, engine, logger, pubKey, privKey)
	r.GET("/health", h.HealthCheck)
	r.POST("/auth/login", h.Login)
	r.POST("/auth/refresh", h.RefreshToken)
	authed := r.Group("/")
	authed.Use(middleware.JWTAuth(pubKey))
	{
		authed.GET("/users", h.ListUsers)
		authed.POST("/users", h.CreateUser)
		authed.GET("/users/:id", h.GetUser)
		authed.PUT("/users/:id", h.UpdateUser)
		authed.DELETE("/users/:id", h.DeleteUser)
		authed.GET("/mappings", h.ListMappings)
		authed.POST("/mappings", h.CreateMapping)
		authed.DELETE("/mappings/:id", h.DeleteMapping)
		authed.GET("/logs", h.ListLogs)
		authed.POST("/sync/trigger", h.TriggerSync)
		authed.GET("/sync/status", h.GetSyncStatus)
	}
	return r
}
