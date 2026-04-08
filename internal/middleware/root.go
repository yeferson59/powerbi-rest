package middleware

import "github.com/yeferson59/powerbi-rest/internal/metrics"

type Middleware struct {
	metricsStore metrics.Store
}

func New(metricsStore metrics.Store) *Middleware {
	return &Middleware{
		metricsStore: metricsStore,
	}
}
