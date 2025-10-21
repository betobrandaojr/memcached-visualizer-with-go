package integration

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"memcached-management/config"
	"memcached-management/handlers"
	"memcached-management/models"
	"memcached-management/services"
)

func setupRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	logger := config.SetupLogger()
	memcachedService := services.NewMemcachedService()
	handler := handlers.NewHandler(memcachedService, logger)

	r := gin.New()
	r.GET("/", handler.ServeIndex)
	r.POST("/connect", handler.HandleConnect)
	r.POST("/set", handler.HandleSet)
	r.POST("/get", handler.HandleGet)
	r.POST("/getMultiple", handler.HandleGetMultiple)
	r.POST("/delete", handler.HandleDelete)
	r.POST("/flush", handler.HandleFlush)
	r.POST("/listKeys", handler.HandleListKeys)

	return r
}

func TestServeIndex(t *testing.T) {
	router := setupRouter()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/", nil)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
	}
}

func TestHandleConnect_InvalidMethod(t *testing.T) {
	router := setupRouter()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/connect", nil)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("Expected status %d, got %d", http.StatusNotFound, w.Code)
	}
}

func TestHandleConnect_EmptyURL(t *testing.T) {
	router := setupRouter()

	connectReq := models.ConnectRequest{URL: ""}
	jsonData, _ := json.Marshal(connectReq)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/connect", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status %d, got %d", http.StatusBadRequest, w.Code)
	}

	var response models.ConnectResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	if err != nil {
		t.Fatal(err)
	}

	if response.Success {
		t.Error("Expected success to be false for empty URL")
	}

	if response.Error != "URL is required" {
		t.Errorf("Expected error 'URL is required', got '%s'", response.Error)
	}
}

func TestHandleConnect_InvalidJSON(t *testing.T) {
	router := setupRouter()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/connect", bytes.NewBuffer([]byte("invalid json")))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status %d, got %d", http.StatusBadRequest, w.Code)
	}
}

func TestHandleSet_NotConnected(t *testing.T) {
	router := setupRouter()

	itemReq := models.ItemRequest{Key: "test", Value: "value"}
	jsonData, _ := json.Marshal(itemReq)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/set", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	var response models.ItemResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	if err != nil {
		t.Fatal(err)
	}

	if response.Success {
		t.Error("Expected success to be false when not connected")
	}

	if response.Error != "Error saving: not connected to Memcached" {
		t.Errorf("Expected error 'Error saving: not connected to Memcached', got '%s'", response.Error)
	}
}

func TestHandleGet_NotConnected(t *testing.T) {
	router := setupRouter()

	itemReq := models.ItemRequest{Key: "test"}
	jsonData, _ := json.Marshal(itemReq)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/get", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	var response models.ItemResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	if err != nil {
		t.Fatal(err)
	}

	if response.Success {
		t.Error("Expected success to be false when not connected")
	}

	if response.Error != "Item not found: not connected to Memcached" {
		t.Errorf("Expected error 'Item not found: not connected to Memcached', got '%s'", response.Error)
	}
}

func TestHandleGetMultiple_NotConnected(t *testing.T) {
	router := setupRouter()

	itemReq := models.ItemRequest{Keys: []string{"key1", "key2"}}
	jsonData, _ := json.Marshal(itemReq)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/getMultiple", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	var response models.ItemResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	if err != nil {
		t.Fatal(err)
	}

	if response.Success {
		t.Error("Expected success to be false when not connected")
	}

	if response.Error != "not connected to Memcached" {
		t.Errorf("Expected error 'not connected to Memcached', got '%s'", response.Error)
	}
}

func TestHandleDelete_NotConnected(t *testing.T) {
	router := setupRouter()

	itemReq := models.ItemRequest{Key: "test"}
	jsonData, _ := json.Marshal(itemReq)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/delete", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	var response models.ItemResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	if err != nil {
		t.Fatal(err)
	}

	if response.Success {
		t.Error("Expected success to be false when not connected")
	}

	if response.Error != "Error deleting: not connected to Memcached" {
		t.Errorf("Expected error 'Error deleting: not connected to Memcached', got '%s'", response.Error)
	}
}

func TestHandleFlush_NotConnected(t *testing.T) {
	router := setupRouter()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/flush", bytes.NewBuffer([]byte("{}")))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	var response models.ItemResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	if err != nil {
		t.Fatal(err)
	}

	if response.Success {
		t.Error("Expected success to be false when not connected")
	}

	if response.Error != "Error flushing cache: not connected to Memcached" {
		t.Errorf("Expected error 'Error flushing cache: not connected to Memcached', got '%s'", response.Error)
	}
}

func TestHandleListKeys_NotConnected(t *testing.T) {
	router := setupRouter()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/listKeys", bytes.NewBuffer([]byte("{}")))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	var response models.ItemResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	if err != nil {
		t.Fatal(err)
	}

	if response.Success {
		t.Error("Expected success to be false when not connected")
	}

	if response.Error != "Error listing keys: not connected to Memcached" {
		t.Errorf("Expected error 'Error listing keys: not connected to Memcached', got '%s'", response.Error)
	}
}