package main

import (
	"memcached-management/config"
	"memcached-management/handlers"
	"memcached-management/services"

	"github.com/gin-gonic/gin"
)

func main() {
	logger := config.SetupLogger()
	memcachedService := services.NewMemcachedService()
	handler := handlers.NewHandler(memcachedService, logger)

	r := gin.Default()

	r.GET("/", handler.ServeIndex)
	r.POST("/connect", handler.HandleConnect)
	r.POST("/set", handler.HandleSet)
	r.POST("/get", handler.HandleGet)
	r.POST("/getMultiple", handler.HandleGetMultiple)
	r.POST("/delete", handler.HandleDelete)
	r.POST("/flush", handler.HandleFlush)
	r.POST("/listKeys", handler.HandleListKeys)

	logger.Info("Server starting on http://localhost:5000")
	r.Run(":5000")
}
