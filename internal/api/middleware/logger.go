package middleware

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
)

func RequestLogger(log zerolog.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		c.Next()
		log.Info().Str("component", "api").Str("method", c.Request.Method).Str("path", c.Request.URL.Path).Int("status", c.Writer.Status()).Int64("duration_ms", time.Since(start).Milliseconds()).Msg("request")
	}
}
