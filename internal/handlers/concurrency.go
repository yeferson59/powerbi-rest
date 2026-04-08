package handlers

import (
	"fmt"
	"math"
	"math/rand"
	"net/http"
	"runtime"
	"sort"
	"strconv"
	"sync"
	"time"

	"github.com/labstack/echo/v5"
)

const (
	defaultConcurrencyN      = 10_000_000
	maxConcurrencyN          = 20_000_000
	defaultBenchmarkRuns     = 3
	maxBenchmarkRuns         = 20
	defaultParallelWorkers   = 2
	maxBenchmarkWorkers      = 64
	randomUpperBound         = 200
	parallelThreadMultiplier = 9973
)

type benchmarkStats struct {
	Strategy            string    `json:"strategy"`
	N                   int       `json:"n"`
	Runs                int       `json:"runs"`
	Workers             int       `json:"workers"`
	AvgMs               float64   `json:"avg_ms"`
	MinMs               float64   `json:"min_ms"`
	MaxMs               float64   `json:"max_ms"`
	P95Ms               float64   `json:"p95_ms"`
	StdDevMs            float64   `json:"stddev_ms"`
	TotalMs             float64   `json:"total_ms"`
	ThroughputOpsPerSec float64   `json:"throughput_ops_per_sec"`
	SamplesMs           []float64 `json:"samples_ms"`
}

type parallelComparison struct {
	Speedup          float64 `json:"speedup"`
	Efficiency       float64 `json:"efficiency"`
	ParallelFraction float64 `json:"parallel_fraction"`
}

type parallelMetricsInput struct {
	N               int `json:"n"`
	Runs            int `json:"runs"`
	ParallelWorkers int `json:"parallel_workers"`
	ThreadWorkers   int `json:"thread_workers"`
}

type parallelMetricsResponse struct {
	Input               parallelMetricsInput          `json:"input"`
	Sequential          benchmarkStats                `json:"sequential"`
	Parallel            benchmarkStats                `json:"parallel"`
	ParallelWithThreads benchmarkStats                `json:"parallel_with_threads"`
	Comparison          map[string]parallelComparison `json:"comparison"`
	TimeSeq             float64                       `json:"timeSeq"`
	TimePar1            float64                       `json:"timePar1"`
	TimePar2            float64                       `json:"timePar2"`
	SpeedUp             float64                       `json:"speedUp"`
	SpeedUp2            float64                       `json:"speedUp2"`
	PorcentParallel     float64                       `json:"porcentParallel"`
	PorcentParallel2    float64                       `json:"porcentParallel2"`
}

func (h *Handler) HandlerSequential(c *echo.Context) error {
	n, err := queryInt(c, "n", defaultConcurrencyN, 1, maxConcurrencyN)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]any{"error": err.Error()})
	}

	runs, err := queryInt(c, "runs", defaultBenchmarkRuns, 1, maxBenchmarkRuns)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]any{"error": err.Error()})
	}

	stats := benchmarkSequential(n, runs)

	c.Set("complexity", "O(n)")
	c.Set("n", n)

	return c.JSON(http.StatusOK, stats)
}

func (h *Handler) HandlerParallel(c *echo.Context) error {
	n, err := queryInt(c, "n", defaultConcurrencyN, 1, maxConcurrencyN)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]any{"error": err.Error()})
	}

	runs, err := queryInt(c, "runs", defaultBenchmarkRuns, 1, maxBenchmarkRuns)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]any{"error": err.Error()})
	}

	workers, err := queryInt(c, "workers", defaultParallelWorkers, 1, maxBenchmarkWorkers)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]any{"error": err.Error()})
	}

	stats := benchmarkParallel("parallel", n, runs, workers)

	c.Set("complexity", "O(n)")
	c.Set("n", n)

	return c.JSON(http.StatusOK, stats)
}

func (h *Handler) HandlerParallelWithThreads(c *echo.Context) error {
	n, err := queryInt(c, "n", defaultConcurrencyN, 1, maxConcurrencyN)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]any{"error": err.Error()})
	}

	runs, err := queryInt(c, "runs", defaultBenchmarkRuns, 1, maxBenchmarkRuns)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]any{"error": err.Error()})
	}

	workers, err := queryInt(c, "workers", defaultThreadWorkers(), 1, maxBenchmarkWorkers)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]any{"error": err.Error()})
	}

	stats := benchmarkParallel("parallel_with_threads", n, runs, workers)

	c.Set("complexity", "O(n)")
	c.Set("n", n)

	return c.JSON(http.StatusOK, stats)
}

func (h *Handler) HandlerParallelMetrics(c *echo.Context) error {
	n, err := queryInt(c, "n", defaultConcurrencyN, 1, maxConcurrencyN)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]any{"error": err.Error()})
	}

	runs, err := queryInt(c, "runs", defaultBenchmarkRuns, 1, maxBenchmarkRuns)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]any{"error": err.Error()})
	}

	parallelWorkers, err := queryInt(c, "parallel_workers", defaultParallelWorkers, 1, maxBenchmarkWorkers)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]any{"error": err.Error()})
	}

	threadWorkers, err := queryInt(c, "thread_workers", defaultThreadWorkers(), 1, maxBenchmarkWorkers)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]any{"error": err.Error()})
	}

	sequentialStats := benchmarkSequential(n, runs)
	parallelStats := benchmarkParallel("parallel", n, runs, parallelWorkers)
	threadStats := benchmarkParallel("parallel_with_threads", n, runs, threadWorkers)

	parallelCmp := compareBenchmarks(sequentialStats, parallelStats)
	threadCmp := compareBenchmarks(sequentialStats, threadStats)

	c.Set("complexity", "O(n)")
	c.Set("n", n)

	return c.JSON(http.StatusOK, parallelMetricsResponse{
		Input: parallelMetricsInput{
			N:               n,
			Runs:            runs,
			ParallelWorkers: parallelWorkers,
			ThreadWorkers:   threadWorkers,
		},
		Sequential:          sequentialStats,
		Parallel:            parallelStats,
		ParallelWithThreads: threadStats,
		Comparison: map[string]parallelComparison{
			"parallel":              parallelCmp,
			"parallel_with_threads": threadCmp,
		},
		TimeSeq:          sequentialStats.AvgMs / 1000.0,
		TimePar1:         parallelStats.AvgMs / 1000.0,
		TimePar2:         threadStats.AvgMs / 1000.0,
		SpeedUp:          parallelCmp.Speedup,
		SpeedUp2:         threadCmp.Speedup,
		PorcentParallel:  parallelCmp.ParallelFraction,
		PorcentParallel2: threadCmp.ParallelFraction,
	})
}

func benchmarkSequential(n, runs int) benchmarkStats {
	samples := make([]time.Duration, 0, runs)
	for range runs {
		samples = append(samples, runSequential(n))
	}

	return calculateStats("sequential", n, runs, 1, samples)
}

func benchmarkParallel(strategy string, n, runs, workers int) benchmarkStats {
	workers = normalizeWorkers(n, workers)

	samples := make([]time.Duration, 0, runs)
	for range runs {
		samples = append(samples, runParallel(n, workers))
	}

	return calculateStats(strategy, n, runs, workers, samples)
}

func runSequential(n int) time.Duration {
	data := make([]int, n)
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))

	start := time.Now()
	for i := range data {
		data[i] = rng.Intn(randomUpperBound)
	}

	return time.Since(start)
}

func runParallel(n, workers int) time.Duration {
	if workers <= 1 {
		return runSequential(n)
	}

	data := make([]int, n)
	workers = normalizeWorkers(n, workers)

	chunkSize := (n + workers - 1) / workers
	var wg sync.WaitGroup
	start := time.Now()

	for worker := range workers {
		startIdx := worker * chunkSize
		if startIdx >= n {
			break
		}

		endIdx := startIdx + chunkSize
		endIdx = min(endIdx, n)

		seed := time.Now().UnixNano() + int64(worker*parallelThreadMultiplier)
		from, to := startIdx, endIdx
		wg.Go(func() {
			rng := rand.New(rand.NewSource(seed))
			for i := from; i < to; i++ {
				data[i] = rng.Intn(randomUpperBound)
			}
		})
	}

	wg.Wait()
	return time.Since(start)
}

func normalizeWorkers(n, workers int) int {
	if workers < 1 {
		return 1
	}
	if workers > n {
		return n
	}

	return workers
}

func calculateStats(strategy string, n, runs, workers int, samples []time.Duration) benchmarkStats {
	samplesMs := make([]float64, len(samples))
	totalMs := 0.0
	minMs := math.MaxFloat64
	maxMs := 0.0

	for i, duration := range samples {
		sampleMs := float64(duration.Microseconds()) / 1000.0
		samplesMs[i] = sampleMs
		totalMs += sampleMs

		if sampleMs < minMs {
			minMs = sampleMs
		}
		if sampleMs > maxMs {
			maxMs = sampleMs
		}
	}

	avgMs := 0.0
	if runs > 0 {
		avgMs = totalMs / float64(runs)
	}

	variance := 0.0
	for _, sampleMs := range samplesMs {
		delta := sampleMs - avgMs
		variance += delta * delta
	}

	stdDevMs := 0.0
	if runs > 0 {
		stdDevMs = math.Sqrt(variance / float64(runs))
	}

	sortedSamples := append([]float64(nil), samplesMs...)
	sort.Float64s(sortedSamples)
	p95Ms := nearestRankPercentile(sortedSamples, 95)

	throughputOpsPerSec := 0.0
	if avgMs > 0 {
		throughputOpsPerSec = float64(n) / (avgMs / 1000.0)
	}

	return benchmarkStats{
		Strategy:            strategy,
		N:                   n,
		Runs:                runs,
		Workers:             workers,
		AvgMs:               avgMs,
		MinMs:               minMs,
		MaxMs:               maxMs,
		P95Ms:               p95Ms,
		StdDevMs:            stdDevMs,
		TotalMs:             totalMs,
		ThroughputOpsPerSec: throughputOpsPerSec,
		SamplesMs:           samplesMs,
	}
}

func nearestRankPercentile(sorted []float64, percentile float64) float64 {
	if len(sorted) == 0 {
		return 0
	}

	rank := int(math.Ceil((percentile / 100.0) * float64(len(sorted))))
	rank = max(rank, 1)
	rank = min(rank, len(sorted))

	return sorted[rank-1]
}

func compareBenchmarks(sequential, parallel benchmarkStats) parallelComparison {
	if parallel.AvgMs <= 0 || parallel.Workers <= 0 {
		return parallelComparison{}
	}

	speedup := sequential.AvgMs / parallel.AvgMs
	efficiency := speedup / float64(parallel.Workers)

	parallelFraction := 0.0
	if parallel.Workers > 1 && speedup > 0 {
		parallelFraction = (1 - (1 / speedup)) / (1 - (1 / float64(parallel.Workers)))
		if parallelFraction < 0 {
			parallelFraction = 0
		}
		if parallelFraction > 1 {
			parallelFraction = 1
		}
	}

	return parallelComparison{
		Speedup:          speedup,
		Efficiency:       efficiency,
		ParallelFraction: parallelFraction,
	}
}

func queryInt(c *echo.Context, key string, defaultValue, min, max int) (int, error) {
	raw := c.QueryParam(key)
	if raw == "" {
		return defaultValue, nil
	}

	value, err := strconv.Atoi(raw)
	if err != nil {
		return 0, fmt.Errorf("%s must be an integer between %d and %d", key, min, max)
	}
	if value < min || value > max {
		return 0, fmt.Errorf("%s must be between %d and %d", key, min, max)
	}

	return value, nil
}

func defaultThreadWorkers() int {
	workers := runtime.NumCPU()
	if workers < 2 {
		return 2
	}
	if workers > maxBenchmarkWorkers {
		return maxBenchmarkWorkers
	}

	return workers
}
