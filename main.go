package main

import (
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/narasimha-swamy/event-trigger/config"
	_ "github.com/narasimha-swamy/event-trigger/docs"
	"github.com/narasimha-swamy/event-trigger/handlers"
	"github.com/narasimha-swamy/event-trigger/jobs"
	"github.com/narasimha-swamy/event-trigger/models"
	"github.com/narasimha-swamy/event-trigger/scheduler"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func main() {
	db := config.ConnectDB()
	db.AutoMigrate(&models.Trigger{}, &models.EventLog{})

	go scheduler.StartScheduler(db)
	go jobs.StartRetentionJob(db)

	router := gin.Default()

	// Serve static files
	router.Static("/static", "./web")
	router.LoadHTMLGlob("web/*.html")

	// Frontend routes
	router.GET("/web", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", nil)
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// Swagger
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// API routes
	api := router.Group("/api")
	{
		api.POST("/triggers", handlers.CreateTrigger(db))
		api.GET("/triggers", handlers.GetTriggers(db))
		api.PUT("/triggers/:id", handlers.UpdateTrigger(db))
		api.DELETE("/triggers/:id", handlers.DeleteTrigger(db))
		api.POST("/triggers/fire/:id", handlers.FireAPITrigger(db))
		api.GET("/events", handlers.GetEvents(db))
		api.POST("/triggers/:id/test", handlers.TestTrigger(db))
	}

	router.Run("0.0.0.0:" + port)
}
