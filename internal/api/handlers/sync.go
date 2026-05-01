package handlers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func (h *Handler) TriggerSync(c *gin.Context) {
	go func() { _ = h.Engine.RunOnce(c) }()
	c.JSON(http.StatusAccepted, gin.H{"message": "sync triggered"})
}

func (h *Handler) GetSyncStatus(c *gin.Context) {
	last, running := h.Engine.Status()
	next := last.Add(time.Duration(h.Cfg.Sync.IntervalSeconds) * time.Second)
	ok(c, gin.H{"last_run": last, "next_run": next, "running": running})
}
