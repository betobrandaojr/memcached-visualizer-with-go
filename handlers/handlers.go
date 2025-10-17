package handlers

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"memcached-management/models"
	"memcached-management/services"
)

type Handler struct {
	memcachedService *services.MemcachedService
	logger           *logrus.Logger
}

func NewHandler(memcachedService *services.MemcachedService, logger *logrus.Logger) *Handler {
	return &Handler{
		memcachedService: memcachedService,
		logger:           logger,
	}
}

func (h *Handler) ServeIndex(c *gin.Context) {
	c.File("web/index.html")
}

func (h *Handler) HandleConnect(c *gin.Context) {
	var req models.ConnectRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.WithError(err).Error("Invalid request data")
		c.JSON(http.StatusBadRequest, models.ConnectResponse{Success: false, Error: "Invalid data"})
		return
	}

	if req.URL == "" {
		c.JSON(http.StatusBadRequest, models.ConnectResponse{Success: false, Error: "URL is required"})
		return
	}

	if err := h.memcachedService.Connect(req.URL); err != nil {
		h.logger.WithError(err).WithField("url", req.URL).Error("Failed to connect to Memcached")
		c.JSON(http.StatusInternalServerError, models.ConnectResponse{Success: false, Error: "Unable to connect: " + err.Error()})
		return
	}

	h.logger.WithField("url", req.URL).Info("Successfully connected to Memcached")
	c.JSON(http.StatusOK, models.ConnectResponse{Success: true, Message: "Connection successful!"})
}

func (h *Handler) HandleSet(c *gin.Context) {
	var req models.ItemRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.WithError(err).Error("Invalid request data")
		c.JSON(http.StatusBadRequest, models.ItemResponse{Success: false, Error: "Invalid data"})
		return
	}

	if err := h.memcachedService.Set(req.Key, req.Value); err != nil {
		h.logger.WithError(err).WithField("key", req.Key).Error("Failed to set item")
		c.JSON(http.StatusInternalServerError, models.ItemResponse{Success: false, Error: "Error saving: " + err.Error()})
		return
	}

	h.logger.WithField("key", req.Key).Info("Item saved successfully")
	c.JSON(http.StatusOK, models.ItemResponse{Success: true, Message: "Item saved successfully!"})
}

func (h *Handler) HandleGet(c *gin.Context) {
	var req models.ItemRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.WithError(err).Error("Invalid request data")
		c.JSON(http.StatusBadRequest, models.ItemResponse{Success: false, Error: "Invalid data"})
		return
	}

	item, err := h.memcachedService.Get(req.Key)
	if err != nil {
		h.logger.WithError(err).WithField("key", req.Key).Warn("Item not found")
		c.JSON(http.StatusNotFound, models.ItemResponse{Success: false, Error: "Item not found: " + err.Error()})
		return
	}

	items := []models.Item{{Key: item.Key, Value: string(item.Value)}}
	c.JSON(http.StatusOK, models.ItemResponse{Success: true, Items: items})
}

func (h *Handler) HandleGetMultiple(c *gin.Context) {
	var req models.ItemRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.WithError(err).Error("Invalid request data")
		c.JSON(http.StatusBadRequest, models.ItemResponse{Success: false, Error: "Invalid data"})
		return
	}

	items, err := h.memcachedService.GetMultiple(req.Keys)
	if err != nil {
		h.logger.WithError(err).WithField("keys", req.Keys).Warn("Failed to get multiple items")
		c.JSON(http.StatusNotFound, models.ItemResponse{Success: false, Error: err.Error()})
		return
	}

	var responseItems []models.Item
	for _, item := range items {
		responseItems = append(responseItems, models.Item{Key: item.Key, Value: string(item.Value)})
	}

	c.JSON(http.StatusOK, models.ItemResponse{Success: true, Items: responseItems})
}

func (h *Handler) HandleDelete(c *gin.Context) {
	var req models.ItemRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.WithError(err).Error("Invalid request data")
		c.JSON(http.StatusBadRequest, models.ItemResponse{Success: false, Error: "Invalid data"})
		return
	}

	if err := h.memcachedService.Delete(req.Key); err != nil {
		h.logger.WithError(err).WithField("key", req.Key).Error("Failed to delete item")
		c.JSON(http.StatusInternalServerError, models.ItemResponse{Success: false, Error: "Error deleting: " + err.Error()})
		return
	}

	h.logger.WithField("key", req.Key).Info("Item deleted successfully")
	c.JSON(http.StatusOK, models.ItemResponse{Success: true, Message: "Item deleted successfully!"})
}

func (h *Handler) HandleFlush(c *gin.Context) {
	if err := h.memcachedService.FlushAll(); err != nil {
		h.logger.WithError(err).Error("Failed to flush cache")
		c.JSON(http.StatusInternalServerError, models.ItemResponse{Success: false, Error: "Error flushing cache: " + err.Error()})
		return
	}

	h.logger.Info("Cache flushed successfully")
	c.JSON(http.StatusOK, models.ItemResponse{Success: true, Message: "All items cleared successfully!"})
}

func (h *Handler) HandleListKeys(c *gin.Context) {
	keys, err := h.memcachedService.GetAllKeys()
	if err != nil {
		h.logger.WithError(err).Error("Failed to list keys")
		c.JSON(http.StatusInternalServerError, models.ItemResponse{Success: false, Error: "Error listing keys: " + err.Error()})
		return
	}

	var items []models.Item
	for _, key := range keys {
		items = append(items, models.Item{Key: key})
	}

	h.logger.WithField("count", len(keys)).Info("Keys listed successfully")
	c.JSON(http.StatusOK, models.ItemResponse{Success: true, Items: items, Message: fmt.Sprintf("Found %d keys", len(keys))})
}