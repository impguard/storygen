package main

import (
	"strconv"

	"github.com/gin-gonic/gin"
)

type status string

const (
	success status = "SUCCESS"
	failure        = "FAILURE"
)

type response struct {
	Status status      `json:"status"`
	Body   interface{} `json:"body,omitempty"`
}

var cache = Cache{}

func main() {
	r := gin.Default()

	r.GET("/print", func(c *gin.Context) {
		c.JSON(200, response{Status: success, Body: cache})
	})

	r.GET("/generate", func(c *gin.Context) {
		lengthQuery := c.DefaultQuery("length", "100")
		length, err := strconv.Atoi(lengthQuery)
		if err != nil {
			length = 100
		}

		story := cache.Generate(length)
		c.JSON(200, response{Status: success, Body: story})
	})

	r.POST("/train", func(c *gin.Context) {
		cache.Train(c.Request.Body)
		c.JSON(200, response{Status: success})
	})

	r.Run()
}
