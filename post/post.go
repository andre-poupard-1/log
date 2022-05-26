package post

import (
	config "main/config"
	"github.com/gin-gonic/gin"
)

func HandlePost(incomingEvent *chan Post) (func (c *gin.Context)) {
	return func(c *gin.Context) {
		var logs []Post
		err := c.ShouldBindJSON(&logs)
		if err != nil {
	
		}
	
		go func() {
			for _, log := range logs {
				*incomingEvent <- log	
			}
		}()
	
		// return 202 to show async processing
		c.JSON(202, nil)
	}
}

func PostController() func(c *gin.Context) {
	incoming := make(chan Post)
	interval := make(chan int)
	buffer := PostBackBuffer[Post]{
		Incoming:    &incoming,
		Interval:    &interval,
		Buffer: make([]Post, 0),
		BufferLimit: config.GetConfig().BatchSize,
		PostbackUrl: config.GetConfig().PostEndpoint,
	}

	go buffer.waitForEvents()
	
	return HandlePost(&incoming)
}