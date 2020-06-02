package main

import (
	"fmt"
	"io"
	"log"
	"os"

	"github.com/gin-gonic/gin"
)

func setupRouter() *gin.Engine {
	// Disable Console Color
	gin.DisableConsoleColor()

	// Access log
	accessLog, err := os.Create("logs/access.log")
	if err != nil {
		fmt.Println(err)
	}
	gin.DefaultWriter = io.MultiWriter(accessLog)

	// Error log
	errorLog, err := os.Create("logs/error.log")
	if err != nil {
		fmt.Println(err)
	}
	gin.DefaultErrorWriter = io.MultiWriter(errorLog)

	r := gin.New()

	// Logger middleware write the logs to gin.DefaultWriter
	r.Use(gin.Logger())

	// Recovery middleware recovers from any panics and writes a 500 if there was one
	r.Use(gin.Recovery())

	// Ping
	r.GET("/ping", func(c *gin.Context) {
		c.String(200, "pong")
	})

	return r
}

func main() {
	r := setupRouter()
	err := r.Run(":8080")
	if err != nil {
		log.Fatal(err)
	}
}
