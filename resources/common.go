package resources

import "time"

type GetRequest struct {
	Key string `json:"key" binding:"required"`
}

type GetResponse struct {
	Value string `json:"value"`
	Found bool   `json:"found"`
}

type SetRequest struct {
	Key        string        `json:"key"`
	Value      interface{}   `json:"value"`
	Expiration time.Duration `json:"expiration"`
}

// ErrorResponse represents an error response payload
type ErrorResponse struct {
	Message string `json:"message"`
}

type GetAllResponse struct {
	Key        string      `json:"key"`
	Value      interface{} `json:"value"`
	Expiration time.Time   `json:"expiration"`
}
