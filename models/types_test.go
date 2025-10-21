package models

import (
	"encoding/json"
	"testing"
)

func TestConnectRequest_JSON(t *testing.T) {
	req := ConnectRequest{URL: "localhost:11211"}
	
	// Test marshaling
	data, err := json.Marshal(req)
	if err != nil {
		t.Errorf("Failed to marshal ConnectRequest: %v", err)
	}
	
	// Test unmarshaling
	var unmarshaled ConnectRequest
	err = json.Unmarshal(data, &unmarshaled)
	if err != nil {
		t.Errorf("Failed to unmarshal ConnectRequest: %v", err)
	}
	
	if unmarshaled.URL != req.URL {
		t.Errorf("Expected URL %s, got %s", req.URL, unmarshaled.URL)
	}
}

func TestConnectResponse_JSON(t *testing.T) {
	resp := ConnectResponse{
		Success: true,
		Message: "Connected successfully",
		Error:   "",
	}
	
	// Test marshaling
	data, err := json.Marshal(resp)
	if err != nil {
		t.Errorf("Failed to marshal ConnectResponse: %v", err)
	}
	
	// Test unmarshaling
	var unmarshaled ConnectResponse
	err = json.Unmarshal(data, &unmarshaled)
	if err != nil {
		t.Errorf("Failed to unmarshal ConnectResponse: %v", err)
	}
	
	if unmarshaled.Success != resp.Success {
		t.Errorf("Expected Success %v, got %v", resp.Success, unmarshaled.Success)
	}
	if unmarshaled.Message != resp.Message {
		t.Errorf("Expected Message %s, got %s", resp.Message, unmarshaled.Message)
	}
}

func TestItemRequest_JSON(t *testing.T) {
	req := ItemRequest{
		Key:   "test:key",
		Value: "test value",
		Keys:  []string{"key1", "key2", "key3"},
	}
	
	// Test marshaling
	data, err := json.Marshal(req)
	if err != nil {
		t.Errorf("Failed to marshal ItemRequest: %v", err)
	}
	
	// Test unmarshaling
	var unmarshaled ItemRequest
	err = json.Unmarshal(data, &unmarshaled)
	if err != nil {
		t.Errorf("Failed to unmarshal ItemRequest: %v", err)
	}
	
	if unmarshaled.Key != req.Key {
		t.Errorf("Expected Key %s, got %s", req.Key, unmarshaled.Key)
	}
	if unmarshaled.Value != req.Value {
		t.Errorf("Expected Value %s, got %s", req.Value, unmarshaled.Value)
	}
	if len(unmarshaled.Keys) != len(req.Keys) {
		t.Errorf("Expected %d keys, got %d", len(req.Keys), len(unmarshaled.Keys))
	}
}

func TestItemResponse_JSON(t *testing.T) {
	items := []Item{
		{Key: "key1", Value: "value1"},
		{Key: "key2", Value: "value2"},
	}
	
	resp := ItemResponse{
		Success: true,
		Message: "Items retrieved successfully",
		Error:   "",
		Items:   items,
	}
	
	// Test marshaling
	data, err := json.Marshal(resp)
	if err != nil {
		t.Errorf("Failed to marshal ItemResponse: %v", err)
	}
	
	// Test unmarshaling
	var unmarshaled ItemResponse
	err = json.Unmarshal(data, &unmarshaled)
	if err != nil {
		t.Errorf("Failed to unmarshal ItemResponse: %v", err)
	}
	
	if unmarshaled.Success != resp.Success {
		t.Errorf("Expected Success %v, got %v", resp.Success, unmarshaled.Success)
	}
	if len(unmarshaled.Items) != len(resp.Items) {
		t.Errorf("Expected %d items, got %d", len(resp.Items), len(unmarshaled.Items))
	}
	if unmarshaled.Items[0].Key != items[0].Key {
		t.Errorf("Expected first item key %s, got %s", items[0].Key, unmarshaled.Items[0].Key)
	}
}

func TestItem_JSON(t *testing.T) {
	item := Item{
		Key:   "test:key",
		Value: "test value",
	}
	
	// Test marshaling
	data, err := json.Marshal(item)
	if err != nil {
		t.Errorf("Failed to marshal Item: %v", err)
	}
	
	// Test unmarshaling
	var unmarshaled Item
	err = json.Unmarshal(data, &unmarshaled)
	if err != nil {
		t.Errorf("Failed to unmarshal Item: %v", err)
	}
	
	if unmarshaled.Key != item.Key {
		t.Errorf("Expected Key %s, got %s", item.Key, unmarshaled.Key)
	}
	if unmarshaled.Value != item.Value {
		t.Errorf("Expected Value %s, got %s", item.Value, unmarshaled.Value)
	}
}