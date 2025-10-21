package services

import (
	"bufio"
	"fmt"
	"net"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/bradfitz/gomemcache/memcache"
)

type MemcachedService struct {
	client *memcache.Client
	host   string
}

func NewMemcachedService() *MemcachedService {
	return &MemcachedService{}
}

func (s *MemcachedService) Connect(url string) error {
	host := strings.TrimSpace(url)
	if host == "" {
		return fmt.Errorf("URL is required")
	}
	
	if strings.Contains(host, "://") {
		return fmt.Errorf("invalid URL format")
	}
	
	if !strings.Contains(host, ":") {
		host += ":11211"
	}

	// Normalize common Docker hostnames to localhost
	if strings.HasPrefix(host, "memcached:") {
		host = strings.Replace(host, "memcached:", "localhost:", 1)
	}

	s.host = host
	s.client = memcache.New(host)
	s.client.Timeout = 5 * time.Second

	return s.client.Ping()
}

func (s *MemcachedService) IsConnected() bool {
	return s.client != nil
}

func (s *MemcachedService) Set(key, value string) error {
	if s.client == nil {
		return fmt.Errorf("not connected to Memcached")
	}
	
	key = strings.TrimSpace(key)
	value = strings.TrimSpace(value)
	
	if key == "" || value == "" {
		return fmt.Errorf("key and value are required")
	}
	
	if len(key) > 250 {
		return fmt.Errorf("key too long (max 250 characters)")
	}

	item := &memcache.Item{Key: key, Value: []byte(value)}
	return s.client.Set(item)
}

func (s *MemcachedService) Get(key string) (*memcache.Item, error) {
	if s.client == nil {
		return nil, fmt.Errorf("not connected to Memcached")
	}
	
	if key == "" {
		return nil, fmt.Errorf("key is required")
	}

	return s.client.Get(key)
}

func (s *MemcachedService) GetMultiple(keys []string) ([]memcache.Item, error) {
	if s.client == nil {
		return nil, fmt.Errorf("not connected to Memcached")
	}
	
	if len(keys) == 0 {
		return nil, fmt.Errorf("at least one key is required")
	}

	var items []memcache.Item
	for _, key := range keys {
		if item, err := s.client.Get(key); err == nil {
			items = append(items, *item)
		}
	}

	if len(items) == 0 {
		return nil, fmt.Errorf("no items found")
	}

	return items, nil
}

func (s *MemcachedService) Delete(key string) error {
	if s.client == nil {
		return fmt.Errorf("not connected to Memcached")
	}
	
	if key == "" {
		return fmt.Errorf("key is required")
	}

	return s.client.Delete(key)
}

func (s *MemcachedService) FlushAll() error {
	if s.client == nil {
		return fmt.Errorf("not connected to Memcached")
	}

	return s.client.FlushAll()
}

func (s *MemcachedService) GetAllKeys() ([]string, error) {
	if s.client == nil {
		return nil, fmt.Errorf("not connected to Memcached")
	}

	conn, err := net.Dial("tcp", s.host)
	if err != nil {
		return nil, fmt.Errorf("failed to connect: %v", err)
	}
	defer conn.Close()

	// Get slab IDs
	fmt.Fprintf(conn, "stats items\r\n")
	scanner := bufio.NewScanner(conn)
	slabIDs := make(map[int]bool)
	
	for scanner.Scan() {
		line := scanner.Text()
		if line == "END" {
			break
		}
		re := regexp.MustCompile(`STAT items:(\d+):number`)
		matches := re.FindStringSubmatch(line)
		if len(matches) > 1 {
			if slabID, err := strconv.Atoi(matches[1]); err == nil {
				slabIDs[slabID] = true
			}
		}
	}

	// Get keys from each slab
	var keys []string
	for slabID := range slabIDs {
		fmt.Fprintf(conn, "stats cachedump %d 0\r\n", slabID)
		for scanner.Scan() {
			line := scanner.Text()
			if line == "END" {
				break
			}
			re := regexp.MustCompile(`ITEM ([^\s]+)`)
			matches := re.FindStringSubmatch(line)
			if len(matches) > 1 {
				keys = append(keys, matches[1])
			}
		}
	}

	return keys, nil
}