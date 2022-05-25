package main

import (
	"fmt"
	"log"
	"time"

	"github.com/gin-gonic/gin"
)

type Log struct {
	UserId int `json:"user_id"`
	Total float64 `json:"total"`
	Title string `json:"title"`
	Metadata Metadata `json:"meta"`
	Completed bool `json:"completed"`
}

type Metadata struct {
	Logins []Login `json:"logins"`
	PhoneNumbers PhoneNumbers `json:"phone_numbers"`
}

type Login struct {
	Time string `json:"time"`
	IP string `json:"ip"`
}

type PhoneNumbers struct {
	Home string `json:"home"`
	Mobile string `json:"mobile"`
}

func main() {
	r := gin.Default()
	r.GET("/healthz", func(c *gin.Context) {
		c.JSON(200, "OK")
	})

	postChannel := make(chan Log)
	bufferIntervalChannel := make(chan int)


	r.POST("/log", func(c *gin.Context) {
		var logs []Log
		err := c.ShouldBindJSON(&logs)
		if err != nil {
			log.Println(err.Error())
		}
		for _, log := range logs {
			// process log
			go func() { postChannel <- log }()
		}

		c.JSON(202, nil)
	})

	r.Run() // listen and serve on 0.0.0.0:8080
}
