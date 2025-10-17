package models

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