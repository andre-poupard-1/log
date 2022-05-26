package main

import (
	"github.com/gin-gonic/gin"
	post "main/post"
	health "main/health"
)

func main() {
	r := gin.Default()
	r.GET("/healthz", health.HealthController)
	r.POST("/log", post.PostController())
	r.Run()
}
