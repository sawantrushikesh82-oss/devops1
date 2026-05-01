package handlers

import (
	"net"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func (h *Handler) HealthCheck(c *gin.Context) {
	dbOK := h.Q.Pool.Ping(c) == nil
	extOK := checkPort(h.Cfg.SFTP.ExternalAddress)
	intOK := checkPort(h.Cfg.SFTP.InternalAddress)
	status := http.StatusOK
	if !dbOK || !extOK || !intOK {
		status = http.StatusServiceUnavailable
	}
	c.JSON(status, gin.H{"api": "ok", "db": dbOK, "ext_sftp": extOK, "int_sftp": intOK})
}

func checkPort(addr string) bool {
	conn, err := net.DialTimeout("tcp", addr, 2*time.Second)
	if err != nil {
		return false
	}
	_ = conn.Close()
	return true
}
