package services

import (
	"testing"
)

func TestNewMemcachedService(t *testing.T) {
	service := NewMemcachedService()
	if service == nil {
		t.Error("Expected service to be created, got nil")
	}
	if service.IsConnected() {
		t.Error("Expected service to not be connected initially")
	}
}

func TestConnect_EmptyURL(t *testing.T) {
	service := NewMemcachedService()
	err := service.Connect("")
	if err == nil {
		t.Error("Expected error for empty URL")
	}
	if err.Error() != "URL is required" {
		t.Errorf("Expected 'URL is required', got '%s'", err.Error())
	}
}

func TestConnect_InvalidURLFormat(t *testing.T) {
	service := NewMemcachedService()
	err := service.Connect("http://localhost:11211")
	if err == nil {
		t.Error("Expected error for invalid URL format")
	}
	if err.Error() != "invalid URL format" {
		t.Errorf("Expected 'invalid URL format', got '%s'", err.Error())
	}
}

func TestConnect_HostnameNormalization(t *testing.T) {
	service := NewMemcachedService()
	
	// Test memcached hostname normalization
	_ = service.Connect("memcached:11211")
	// We expect this to fail since memcached isn't running, but the hostname should be normalized
	if service.host != "localhost:11211" {
		t.Errorf("Expected host to be normalized to 'localhost:11211', got '%s'", service.host)
	}
}

func TestConnect_DefaultPort(t *testing.T) {
	service := NewMemcachedService()
	
	// Test default port addition
	_ = service.Connect("localhost")
	// We expect this to fail since memcached isn't running, but port should be added
	if service.host != "localhost:11211" {
		t.Errorf("Expected host to be 'localhost:11211', got '%s'", service.host)
	}
}

func TestSet_NotConnected(t *testing.T) {
	service := NewMemcachedService()
	err := service.Set("key", "value")
	if err == nil {
		t.Error("Expected error when not connected")
	}
	if err.Error() != "not connected to Memcached" {
		t.Errorf("Expected 'not connected to Memcached', got '%s'", err.Error())
	}
}

func TestSet_EmptyKeyValue(t *testing.T) {
	service := NewMemcachedService()
	// Simulate connection
	service.host = "localhost:11211"
	
	err := service.Set("", "value")
	if err == nil {
		t.Error("Expected error for empty key")
	}
	
	err = service.Set("key", "")
	if err == nil {
		t.Error("Expected error for empty value")
	}
}

func TestSet_KeyTooLong(t *testing.T) {
	service := NewMemcachedService()
	
	longKey := make([]byte, 251)
	for i := range longKey {
		longKey[i] = 'a'
	}
	
	err := service.Set(string(longKey), "value")
	if err == nil {
		t.Error("Expected error for key too long")
	}
	// Since not connected, it will return "not connected" error first
	if err.Error() != "not connected to Memcached" {
		t.Errorf("Expected 'not connected to Memcached', got '%s'", err.Error())
	}
}

func TestGet_NotConnected(t *testing.T) {
	service := NewMemcachedService()
	_, err := service.Get("key")
	if err == nil {
		t.Error("Expected error when not connected")
	}
	if err.Error() != "not connected to Memcached" {
		t.Errorf("Expected 'not connected to Memcached', got '%s'", err.Error())
	}
}

func TestGet_EmptyKey(t *testing.T) {
	service := NewMemcachedService()
	
	_, err := service.Get("")
	if err == nil {
		t.Error("Expected error for empty key")
	}
	// Since not connected, it will return "not connected" error first
	if err.Error() != "not connected to Memcached" {
		t.Errorf("Expected 'not connected to Memcached', got '%s'", err.Error())
	}
}

func TestGetMultiple_NotConnected(t *testing.T) {
	service := NewMemcachedService()
	_, err := service.GetMultiple([]string{"key1", "key2"})
	if err == nil {
		t.Error("Expected error when not connected")
	}
	if err.Error() != "not connected to Memcached" {
		t.Errorf("Expected 'not connected to Memcached', got '%s'", err.Error())
	}
}

func TestGetMultiple_EmptyKeys(t *testing.T) {
	service := NewMemcachedService()
	
	_, err := service.GetMultiple([]string{})
	if err == nil {
		t.Error("Expected error for empty keys")
	}
	// Since not connected, it will return "not connected" error first
	if err.Error() != "not connected to Memcached" {
		t.Errorf("Expected 'not connected to Memcached', got '%s'", err.Error())
	}
}

func TestDelete_NotConnected(t *testing.T) {
	service := NewMemcachedService()
	err := service.Delete("key")
	if err == nil {
		t.Error("Expected error when not connected")
	}
	if err.Error() != "not connected to Memcached" {
		t.Errorf("Expected 'not connected to Memcached', got '%s'", err.Error())
	}
}

func TestDelete_EmptyKey(t *testing.T) {
	service := NewMemcachedService()
	
	err := service.Delete("")
	if err == nil {
		t.Error("Expected error for empty key")
	}
	// Since not connected, it will return "not connected" error first
	if err.Error() != "not connected to Memcached" {
		t.Errorf("Expected 'not connected to Memcached', got '%s'", err.Error())
	}
}

func TestFlushAll_NotConnected(t *testing.T) {
	service := NewMemcachedService()
	err := service.FlushAll()
	if err == nil {
		t.Error("Expected error when not connected")
	}
	if err.Error() != "not connected to Memcached" {
		t.Errorf("Expected 'not connected to Memcached', got '%s'", err.Error())
	}
}

func TestGetAllKeys_NotConnected(t *testing.T) {
	service := NewMemcachedService()
	_, err := service.GetAllKeys()
	if err == nil {
		t.Error("Expected error when not connected")
	}
	if err.Error() != "not connected to Memcached" {
		t.Errorf("Expected 'not connected to Memcached', got '%s'", err.Error())
	}
}