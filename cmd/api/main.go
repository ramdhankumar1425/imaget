package main

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/ramdhankumar1425/imaget/internal/handler"
	"github.com/ramdhankumar1425/imaget/internal/infra"
	"github.com/ramdhankumar1425/imaget/internal/worker"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("error loading .env file")
	}

	infra.InitRedis()
	infra.InitImageKit()
	worker.InitPool()

	r := gin.Default()

	r.POST("/transform", handler.HandleTransform)
	r.GET("/transform/:id", handler.HandleGetResult)
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "8000"
	}

	r.Run(":" + port)
}
