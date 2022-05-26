package main

import (
	"fmt"
	health "main/health"
	middleware "main/middleware"
	post "main/post"
	"os"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func main() {
	logger, _ := zap.NewProduction()	
	r := gin.New()
	r.Use(middleware.Logger(logger))
	r.Use(gin.Recovery())

	r.GET("/health", health.HealthController)
	r.POST("/log", post.PostController(logger))
	
	port := os.Getenv("PORT")
	logger.Info(fmt.Sprintf("Listening on port %s", port))
	err := r.Run(fmt.Sprintf(":%s", port))
	logger.Error(err.Error())
}
