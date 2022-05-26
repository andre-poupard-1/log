package health

import (
	"github.com/gin-gonic/gin"
)

func HealthController(c *gin.Context) {
	c.JSON(200, "OK")
}