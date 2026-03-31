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
	apiRoutes := map[string]string{
		"/o1":     "O(1)",
		"/on":     "O(n)",
		"/onlogn": "O(n log n)",
		"/on2":    "O(n²)",
		"/o2n":    "O(2^n)",
	}

	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c *echo.Context) error {
			start := time.Now()

			err := next(c)

			path := c.Path()

			if path == "/" || path == "/summary" {
				return err
			}

			elapsed := float64(time.Since(start).Microseconds()) / 1000.0

			statusCode := 200
			if resp, ok := c.Response().(*echo.Response); ok {
				statusCode = resp.Status
				if statusCode == 0 {
					statusCode = 200
				}
			}

			complexity := ""
			n := 0

			if val, ok := c.Get("complexity").(string); ok {
				complexity = val
			} else if apiComplexity, exists := apiRoutes[path]; exists {
				complexity = apiComplexity
			}

			if val, ok := c.Get("n").(int); ok {
				n = val
			}

			metric := Metric{
				Route:       path,
				Method:      c.Request().Method,
				Complexity:  complexity,
				NParam:      n,
				ResponseMs:  elapsed,
				StatusCode:  statusCode,
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
