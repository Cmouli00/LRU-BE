package main

import (
	"lru/controller"
	"time"

	"lru/service"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	cache := service.NewLRUCache(10)

	router := gin.Default()

	// Configure CORS
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	handler := controller.NewCacheHandler(cache)

	router.GET("/cache/get/:key", handler.Get)
	router.POST("/cache/set", handler.Set)
	router.DELETE("/cache/delete/:key", handler.Delete)
	router.GET("/cache/getall", handler.GetAll)

	router.Run(":8080")
}
