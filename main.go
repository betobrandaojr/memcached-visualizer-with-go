package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/bradfitz/gomemcache/memcache"
)

var mc *memcache.Client

type ConnectRequest struct {
	URL string `json:"url"`
}

type ConnectResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message,omitempty"`
	Error   string `json:"error,omitempty"`
}

type ItemRequest struct {
	Key   string   `json:"key"`
	Keys  []string `json:"keys"`
	Value string   `json:"value"`
}

type ItemResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message,omitempty"`
	Error   string `json:"error,omitempty"`
	Items   []Item `json:"items,omitempty"`
}

type Item struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

func main() {
	http.HandleFunc("/", serveIndex)
	http.HandleFunc("/connect", handleConnect)
	http.HandleFunc("/set", handleSet)
	http.HandleFunc("/get", handleGet)
	http.HandleFunc("/getMultiple", handleGetMultiple)
	http.HandleFunc("/delete", handleDelete)

	fmt.Println("Server running on http://localhost:5000")
	if err := http.ListenAndServe(":5000", nil); err != nil {
		fmt.Printf("Server failed to start: %v\n", err)
	}
}

func serveIndex(w http.ResponseWriter, r *http.Request) {
	// Validate path to prevent directory traversal
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}
	http.ServeFile(w, r, "templates/index.html")
}

func handleConnect(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	var req ConnectRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response := ConnectResponse{Success: false, Error: "Invalid data"}
		json.NewEncoder(w).Encode(response)
		return
	}

	if req.URL == "" {
		response := ConnectResponse{Success: false, Error: "URL is required"}
		json.NewEncoder(w).Encode(response)
		return
	}

	// Validate and parse URL
	host := strings.TrimSpace(req.URL)
	if host == "" {
		response := ConnectResponse{Success: false, Error: "URL is required"}
		json.NewEncoder(w).Encode(response)
		return
	}
	
	// Basic validation to prevent SSRF
	if strings.Contains(host, "://") {
		response := ConnectResponse{Success: false, Error: "Invalid URL format"}
		json.NewEncoder(w).Encode(response)
		return
	}
	
	if !strings.Contains(host, ":") {
		host += ":11211"
	}

	// Test connection
	mc = memcache.New(host)
	mc.Timeout = 5 * time.Second

	if err := mc.Ping(); err != nil {
		response := ConnectResponse{Success: false, Error: fmt.Sprintf("Unable to connect: %v", err)}
		json.NewEncoder(w).Encode(response)
		return
	}

	response := ConnectResponse{Success: true, Message: "Connection successful!"}
	json.NewEncoder(w).Encode(response)
}

func handleSet(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	if mc == nil {
		response := ItemResponse{Success: false, Error: "Not connected to Memcached"}
		json.NewEncoder(w).Encode(response)
		return
	}

	var req ItemRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response := ItemResponse{Success: false, Error: "Invalid data"}
		json.NewEncoder(w).Encode(response)
		return
	}

	// Validate and sanitize input
	req.Key = strings.TrimSpace(req.Key)
	req.Value = strings.TrimSpace(req.Value)
	
	if req.Key == "" || req.Value == "" {
		response := ItemResponse{Success: false, Error: "Key and value are required"}
		json.NewEncoder(w).Encode(response)
		return
	}
	
	// Validate key length (memcached limit is 250 bytes)
	if len(req.Key) > 250 {
		response := ItemResponse{Success: false, Error: "Key too long (max 250 characters)"}
		json.NewEncoder(w).Encode(response)
		return
	}

	item := &memcache.Item{Key: req.Key, Value: []byte(req.Value)}
	if err := mc.Set(item); err != nil {
		response := ItemResponse{Success: false, Error: fmt.Sprintf("Error saving: %v", err)}
		json.NewEncoder(w).Encode(response)
		return
	}

	response := ItemResponse{Success: true, Message: "Item saved successfully!"}
	json.NewEncoder(w).Encode(response)
}

func handleGet(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	if mc == nil {
		response := ItemResponse{Success: false, Error: "Not connected to Memcached"}
		json.NewEncoder(w).Encode(response)
		return
	}

	var req ItemRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response := ItemResponse{Success: false, Error: "Invalid data"}
		json.NewEncoder(w).Encode(response)
		return
	}

	if req.Key == "" {
		response := ItemResponse{Success: false, Error: "Key is required"}
		json.NewEncoder(w).Encode(response)
		return
	}

	item, err := mc.Get(req.Key)
	if err != nil {
		response := ItemResponse{Success: false, Error: fmt.Sprintf("Item not found: %v", err)}
		json.NewEncoder(w).Encode(response)
		return
	}

	items := []Item{{Key: item.Key, Value: string(item.Value)}}
	response := ItemResponse{Success: true, Items: items}
	json.NewEncoder(w).Encode(response)
}

func handleDelete(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	if mc == nil {
		response := ItemResponse{Success: false, Error: "Not connected to Memcached"}
		json.NewEncoder(w).Encode(response)
		return
	}

	var req ItemRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response := ItemResponse{Success: false, Error: "Invalid data"}
		json.NewEncoder(w).Encode(response)
		return
	}

	if req.Key == "" {
		response := ItemResponse{Success: false, Error: "Key is required"}
		json.NewEncoder(w).Encode(response)
		return
	}

	if err := mc.Delete(req.Key); err != nil {
		response := ItemResponse{Success: false, Error: fmt.Sprintf("Error deleting: %v", err)}
		json.NewEncoder(w).Encode(response)
		return
	}

	response := ItemResponse{Success: true, Message: "Item deleted successfully!"}
	json.NewEncoder(w).Encode(response)
}

func handleGetMultiple(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	if mc == nil {
		response := ItemResponse{Success: false, Error: "Not connected to Memcached"}
		json.NewEncoder(w).Encode(response)
		return
	}

	var req ItemRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response := ItemResponse{Success: false, Error: "Invalid data"}
		json.NewEncoder(w).Encode(response)
		return
	}

	if len(req.Keys) == 0 {
		response := ItemResponse{Success: false, Error: "At least one key is required"}
		json.NewEncoder(w).Encode(response)
		return
	}

	items := []Item{}
	for _, key := range req.Keys {
		if item, err := mc.Get(key); err == nil {
			items = append(items, Item{Key: item.Key, Value: string(item.Value)})
		}
	}

	if len(items) == 0 {
		response := ItemResponse{Success: false, Error: "No items found"}
		json.NewEncoder(w).Encode(response)
		return
	}

	response := ItemResponse{Success: true, Items: items}
	json.NewEncoder(w).Encode(response)
}
