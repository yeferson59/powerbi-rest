package metrics

import (
	"context"
	"time"
)

type Record struct {
	ID           int       `json:"id"`
	RequestID    string    `json:"request_id"`
	Route        string    `json:"route"`
	Method       string    `json:"method"`
	Complexity   string    `json:"complexity"`
	NParam       int       `json:"n_param"`
	ResponseMs   float64   `json:"response_ms"`
	StatusCode   int       `json:"status_code"`
	ErrorMessage string    `json:"error_message"`
	CreatedAt    time.Time `json:"created_at"`
}

type CreateInput struct {
	Route       string
	Method      string
	Complexity  string
	NParam      int
	ResponseMs  float64
	StatusCode  int
	ErrorMsg    string
	RequestedAt time.Time
}

type Store interface {
	List(ctx context.Context) ([]Record, error)
	Create(ctx context.Context, metric CreateInput) error
}
