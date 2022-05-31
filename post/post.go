package post

import (
	config "main/config"
	"net/http"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func HandlePost(incomingEvent *chan Post, logger *zap.Logger) (func (c *gin.Context)) {
	return func(c *gin.Context) {
		var log Post
		err := c.ShouldBindJSON(&log)
		if err != nil {
			logger.Error(err.Error())
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"error": "JSON body validation failed",
				"message": err.Error(),
			})
			return
		}
	
		go func() { *incomingEvent <- log	}()
		
		// return 202 to show async processing
		c.JSON(202, nil)
	}
}

func PostController(logger *zap.Logger) func(c *gin.Context) {
	incoming := make(chan Post)
	interval := make(chan bool)
	buffer := PostBackBuffer[Post]{
		Incoming:    &incoming,
		Interval:    &interval,
		Buffer: make([]Post, 0),
		BufferLimit: config.GetConfig().BatchSize,
		PostbackUrl: config.GetConfig().PostEndpoint,
		Logger: logger,
	}

	go buffer.waitForEvents()
	
	return HandlePost(&incoming, logger)
}