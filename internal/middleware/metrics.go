package middleware

import (
	"context"
	"log"
	"time"

	"github.com/labstack/echo/v5"
)

type Metric struct {
	Route       string
	Method      string
	Complexity  string
	NParam      int
	ResponseMs  float64
	StatusCode  int
	ErrorMsg    string
	RequestedAt time.Time
}

func (m *Middleware) Metrics() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c *echo.Context) error {
			start := time.Now()

			err := next(c)

			elapsed := float64(time.Since(start).Microseconds()) / 1000.0

			var n int
			if val, ok := c.Get("n").(int); ok {
				n = val
			}
			var complexity string
			if val, ok := c.Get("complexity").(string); ok {
				complexity = val
			}

			metric := Metric{
				Route:       c.Path(),
				Method:      c.Request().Method,
				Complexity:  complexity,
				NParam:      n,
				ResponseMs:  elapsed,
				StatusCode:  0,
				RequestedAt: time.Now(),
			}

			if err != nil {
				metric.ErrorMsg = err.Error()
			}

			go func() {
				ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
				defer cancel()

				if _, errExec := m.db.Exec(ctx, "INSERT INTO metrics (route, method, complexity, n_param, response_ms, status_code, error_message, created_at) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)", metric.Route, metric.Method, metric.Complexity, metric.NParam, metric.ResponseMs, metric.StatusCode, metric.ErrorMsg, metric.RequestedAt); errExec != nil {
					log.Fatal(errExec)
				}
			}()

			return err
		}
	}
}
