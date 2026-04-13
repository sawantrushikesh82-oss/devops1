package handlers

import (
	"strconv"

	"github.com/gin-gonic/gin"
)

func (h *Handler) ListLogs(c *gin.Context) {
	page, _ := strconv.Atoi(defaultString(c.Query("page"), "1"))
	limit, _ := strconv.Atoi(defaultString(c.Query("limit"), "20"))
	logs, err := h.Q.GetLogs(c, page, limit, c.Query("status"), c.Query("username"))
	if err != nil {
		errResp(c, 500, "INTERNAL_ERROR", "failed to list logs")
		return
	}
	ok(c, logs)
}

func defaultString(v, d string) string {
	if v == "" {
		return d
	}
	return v
}
