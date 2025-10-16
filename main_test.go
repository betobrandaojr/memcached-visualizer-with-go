package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	
	"github.com/bradfitz/gomemcache/memcache"
)

func TestServeIndex(t *testing.T) {
	req, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(serveIndex)
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}
}

func TestHandleConnect_InvalidMethod(t *testing.T) {
	req, err := http.NewRequest("GET", "/connect", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(handleConnect)
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusMethodNotAllowed {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusMethodNotAllowed)
	}
}

func TestHandleConnect_EmptyURL(t *testing.T) {
	connectReq := ConnectRequest{URL: ""}
	jsonData, _ := json.Marshal(connectReq)

	req, err := http.NewRequest("POST", "/connect", bytes.NewBuffer(jsonData))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(handleConnect)
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	var response ConnectResponse
	err = json.Unmarshal(rr.Body.Bytes(), &response)
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

func TestHandleSet_NotConnected(t *testing.T) {
	mc = nil // Ensure not connected

	itemReq := ItemRequest{Key: "test", Value: "value"}
	jsonData, _ := json.Marshal(itemReq)

	req, err := http.NewRequest("POST", "/set", bytes.NewBuffer(jsonData))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(handleSet)
	handler.ServeHTTP(rr, req)

	var response ItemResponse
	err = json.Unmarshal(rr.Body.Bytes(), &response)
	if err != nil {
		t.Fatal(err)
	}

	if response.Success {
		t.Error("Expected success to be false when not connected")
	}

	if response.Error != "Not connected to Memcached" {
		t.Errorf("Expected error 'Not connected to Memcached', got '%s'", response.Error)
	}
}

func TestHandleSet_EmptyKeyValue(t *testing.T) {
	// Mock connection to test validation
	mc = &memcache.Client{} // Mock connection
	
	itemReq := ItemRequest{Key: "", Value: ""}
	jsonData, _ := json.Marshal(itemReq)

	req, err := http.NewRequest("POST", "/set", bytes.NewBuffer(jsonData))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(handleSet)
	handler.ServeHTTP(rr, req)

	var response ItemResponse
	err = json.Unmarshal(rr.Body.Bytes(), &response)
	if err != nil {
		t.Fatal(err)
	}

	if response.Success {
		t.Error("Expected success to be false for empty key/value")
	}

	if response.Error != "Key and value are required" {
		t.Errorf("Expected error 'Key and value are required', got '%s'", response.Error)
	}
	
	mc = nil // Reset
}

func TestHandleGet_NotConnected(t *testing.T) {
	mc = nil // Ensure not connected

	itemReq := ItemRequest{Key: "test"}
	jsonData, _ := json.Marshal(itemReq)

	req, err := http.NewRequest("POST", "/get", bytes.NewBuffer(jsonData))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(handleGet)
	handler.ServeHTTP(rr, req)

	var response ItemResponse
	err = json.Unmarshal(rr.Body.Bytes(), &response)
	if err != nil {
		t.Fatal(err)
	}

	if response.Success {
		t.Error("Expected success to be false when not connected")
	}

	if response.Error != "Not connected to Memcached" {
		t.Errorf("Expected error 'Not connected to Memcached', got '%s'", response.Error)
	}
}

func TestHandleDelete_EmptyKey(t *testing.T) {
	// Mock connection to test validation
	mc = &memcache.Client{} // Mock connection
	
	itemReq := ItemRequest{Key: ""}
	jsonData, _ := json.Marshal(itemReq)

	req, err := http.NewRequest("POST", "/delete", bytes.NewBuffer(jsonData))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(handleDelete)
	handler.ServeHTTP(rr, req)

	var response ItemResponse
	err = json.Unmarshal(rr.Body.Bytes(), &response)
	if err != nil {
		t.Fatal(err)
	}

	if response.Success {
		t.Error("Expected success to be false for empty key")
	}

	if response.Error != "Key is required" {
		t.Errorf("Expected error 'Key is required', got '%s'", response.Error)
	}
	
	mc = nil // Reset
}

func TestHandleGetMultiple_EmptyKeys(t *testing.T) {
	// Mock connection to test validation
	mc = &memcache.Client{} // Mock connection
	
	itemReq := ItemRequest{Keys: []string{}}
	jsonData, _ := json.Marshal(itemReq)

	req, err := http.NewRequest("POST", "/getMultiple", bytes.NewBuffer(jsonData))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(handleGetMultiple)
	handler.ServeHTTP(rr, req)

	var response ItemResponse
	err = json.Unmarshal(rr.Body.Bytes(), &response)
	if err != nil {
		t.Fatal(err)
	}

	if response.Success {
		t.Error("Expected success to be false for empty keys")
	}

	if response.Error != "At least one key is required" {
		t.Errorf("Expected error 'At least one key is required', got '%s'", response.Error)
	}
	
	mc = nil // Reset
}