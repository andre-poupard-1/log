package logger

import (
	"time"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

const (
	INCOMING_REQUEST = "INCOMING_REQUEST"
	FAILED_FINAL_BACKOFF_REQUEST = "FAILED_FINAL_BACKOFF_REQUEST"
	OUTGOING_POST_REQUEST = "OUTGOING_POST_REQUEST"
)

func Logger(logger *zap.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		beforeReq := time.Now()
		c.Next()
		latency := time.Since(beforeReq)
		logger.Info(INCOMING_REQUEST,
			zap.String("method", c.Request.Method),
			zap.String("path", c.Request.URL.Path),
			zap.Int64("nanosecond_latency", latency.Nanoseconds()),
		)
	}
}