// @title Uptime Kuma Monitor API
// @version 1.0
// @description A Go REST API for managing Uptime Kuma monitors with MariaDB/MySQL backend
// @host localhost:8080
// @BasePath /api/v1
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and the token.
package main

import (
	"log"

	"zimbres/uptime-kuma-api/internal/config"
	"zimbres/uptime-kuma-api/internal/database"
	"zimbres/uptime-kuma-api/internal/handler"
	"zimbres/uptime-kuma-api/internal/middleware"
	"zimbres/uptime-kuma-api/internal/repository"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	_ "zimbres/uptime-kuma-api/docs"
)

func main() {
	cfg := config.LoadConfig()

	database.InitDB(cfg)
	defer database.CloseDB()

	monitorRepo := repository.NewMonitorRepository(database.DB)

	tagRepo := repository.NewTagRepository(database.DB)
	monitorHandler := handler.NewMonitorHandler(monitorRepo, tagRepo)
	tagHandler := handler.NewTagHandler(tagRepo)

	heartbeatRepo := repository.NewHeartbeatRepository(database.DB)
	heartbeatHandler := handler.NewHeartbeatHandler(heartbeatRepo)

	statsRepo := repository.NewStatsRepository(database.DB)
	statsHandler := handler.NewStatsHandler(statsRepo)

	monitorTagRepo := repository.NewMonitorTagRepository(database.DB)
	monitorTagHandler := handler.NewMonitorTagHandler(monitorTagRepo)

	maintenanceRepo := repository.NewMaintenanceRepository(database.DB)
	maintenanceHandler := handler.NewMaintenanceHandler(maintenanceRepo)

	monitorMaintenanceRepo := repository.NewMonitorMaintenanceRepository(database.DB)
	monitorMaintenanceHandler := handler.NewMonitorMaintenanceHandler(monitorMaintenanceRepo)

	r := gin.Default()

	r.Use(gin.Recovery())
	r.Use(gin.Logger())

	// Health check
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	// Swagger documentation
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// API routes
	api := r.Group("/api/v1")
	api.Use(middleware.AuthMiddleware(cfg))
	{
		monitors := api.Group("/monitors")
		{
			monitors.POST("", monitorHandler.CreateMonitor)
			monitors.GET("", monitorHandler.GetMonitors)
			monitors.GET("/:id", monitorHandler.GetMonitor)
			monitors.PUT("/:id", monitorHandler.UpdateMonitor)
			monitors.DELETE("/:id", monitorHandler.DeleteMonitor)
			monitors.POST("/:id/pause", monitorHandler.PauseMonitor)
			monitors.POST("/:id/resume", monitorHandler.ResumeMonitor)
			monitors.GET("/:id/heartbeat", heartbeatHandler.GetMonitorLastHeartbeat)
			monitors.GET("/:id/heartbeats", heartbeatHandler.GetMonitorHeartbeats)
			monitors.GET("/:id/stats", statsHandler.GetMonitorStats)
			monitors.POST("/:id/tags", monitorTagHandler.AddMonitorTag)
			monitors.GET("/:id/tags", monitorTagHandler.GetMonitorTags)
			monitors.DELETE("/:id/tags/:tagId", monitorTagHandler.DeleteMonitorTag)
			monitors.POST("/:id/maintenances", monitorMaintenanceHandler.AddMonitorMaintenance)
			monitors.GET("/:id/maintenances", monitorMaintenanceHandler.GetMonitorMaintenances)
			monitors.DELETE("/:id/maintenances/:maintenanceId", monitorMaintenanceHandler.DeleteMonitorMaintenance)
		}

		tags := api.Group("/tags")
		{
			tags.POST("", tagHandler.CreateTag)
			tags.GET("", tagHandler.GetTags)
			tags.GET("/:id", tagHandler.GetTag)
			tags.PUT("/:id", tagHandler.UpdateTag)
			tags.DELETE("/:id", tagHandler.DeleteTag)
		}

		maintenances := api.Group("/maintenances")
		{
			maintenances.POST("", maintenanceHandler.CreateMaintenance)
			maintenances.GET("", maintenanceHandler.GetMaintenances)
			maintenances.GET("/:id", maintenanceHandler.GetMaintenance)
			maintenances.PUT("/:id", maintenanceHandler.UpdateMaintenance)
			maintenances.DELETE("/:id", maintenanceHandler.DeleteMaintenance)
		}
	}

	port := cfg.ServerPort
	log.Printf("Server starting on port %s", port)
	log.Printf("Swagger UI available at http://localhost:%s/swagger/index.html", port)

	if err := r.Run(":" + port); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}