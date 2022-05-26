package main

import (
	"fmt"
	health "main/health"
	post "main/post"
	"time"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.New()
	r.Use(gin.LoggerWithFormatter(func(param gin.LogFormatterParams) string {
		return fmt.Sprintf("%s - [%s] \"%s %s %s %d %s \" %s\"\n",
				param.ClientIP,
				param.TimeStamp.Format(time.RFC1123),
				param.Method,
				param.Path,
				param.Request.Proto,
				param.StatusCode,
				param.Latency,
				param.ErrorMessage,
		)
	}))
	r.Use(gin.Recovery())
	r.GET("/healthz", health.HealthController)
	r.POST("/log", post.PostController())
	r.Run()
}
