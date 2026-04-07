package handlers

import (
	"fmt"
	"math/rand"
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/labstack/echo/v5"
)

type Handler struct {
	db *pgxpool.Pool
}

func New(db *pgxpool.Pool) Handler {
	return Handler{
		db: db,
	}
}

func (Handler) HandlerRoot(c *echo.Context) error {
	return c.File("dashboard.html")
}

func (Handler) HandlerO1(c *echo.Context) error {
	c.Set("complexity", "O(1)")
	c.Set("n", 1)

	lookup := map[string]int{"go": 1, "fiber": 2, "powerbi": 3}
	result := lookup["fiber"]

	return c.JSON(http.StatusOK, map[string]any{"result": result, "complexity": "O(1)"})
}

func (Handler) HandlerOn(c *echo.Context) error {
	n, err := strconv.Atoi(c.QueryParamOr("n", "100"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]any{"error": "invalid n"})
	}

	c.Set("complexity", "O(n)")
	c.Set("n", n)

	sum := 0
	for i := range n {
		sum += i
	}

	return c.JSON(http.StatusOK, map[string]any{"result": sum, "n": n, "complexity": "O(n)"})
}

func generateRandomSlice(n int) []int {
	data := make([]int, n)

	for i := range data {
		data[i] = rand.Intn(300)
	}

	return data
}

func mergeSort(data []int) []int {
	if len(data) <= 1 {
		return data
	}

	mid := len(data) / 2
	left := mergeSort(data[:mid])
	right := mergeSort(data[mid:])

	return merge(left, right)
}

func merge(left, right []int) []int {
	result := make([]int, 0, len(left)+len(right))
	i, j := 0, 0

	for i < len(left) && j < len(right) {
		if left[i] < right[j] {
			result = append(result, left[i])
			i++
		} else {
			result = append(result, right[j])
			j++
		}
	}

	result = append(result, left[i:]...)
	result = append(result, right[j:]...)
	return result
}

func (Handler) HandlerONLogN(c *echo.Context) error {
	n, err := strconv.Atoi(c.QueryParamOr("n", "1000"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]any{"error": "invalid n"})
	}

	c.Set("complexity", "O(n log n)")
	c.Set("n", n)

	data := generateRandomSlice(n)
	sorted := mergeSort(data)
	return c.JSON(http.StatusOK, map[string]any{"sorted_len": len(sorted), "n": n, "complexity": "O(n log n)"})
}

func (Handler) HandlerON2(c *echo.Context) error {
	n, err := strconv.Atoi(c.QueryParamOr("n", "500"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]any{"error": "invalid n"})
	}
	c.Set("complexity", "O(n²)")
	c.Set("n", n)

	data := generateRandomSlice(n)
	for i := range len(data) {
		for j := 0; j < len(data)-i-1; j++ {
			if data[j] > data[j+1] {
				data[j], data[j+1] = data[j+1], data[j]
			}
		}
	}
	return c.JSON(http.StatusOK, map[string]any{"sorted_len": len(data), "n": n, "complexity": "O(n²)"})
}

func fibRecursive(n int) int {
	if n <= 1 {
		return n
	}

	return fibRecursive(n-1) + fibRecursive(n-2)
}

func (Handler) HandlerO2N(c *echo.Context) error {
	n, err := strconv.Atoi(c.QueryParamOr("n", "20"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]any{"error": "invalid n"})
	}

	if n > 35 {
		n = 35
	}

	c.Set("complexity", "O(2^n)")
	c.Set("n", n)

	result := fibRecursive(n)

	return c.JSON(http.StatusOK, map[string]any{"result": result, "n": n, "complexity": "O(2^n)"})
}

type Metric struct {
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

func (h *Handler) HandlerSummary(c *echo.Context) error {
	var metrics []Metric

	rows, err := h.db.Query(c.Request().Context(), "SELECT id, request_id, route, method, complexity, n_param, response_ms, status_code, error_message, created_at FROM metrics")
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]any{"error": err.Error()})
	}
	defer rows.Close()

	for rows.Next() {
		var m Metric
		err := rows.Scan(&m.ID, &m.RequestID, &m.Route, &m.Method, &m.Complexity, &m.NParam, &m.ResponseMs, &m.StatusCode, &m.ErrorMessage, &m.CreatedAt)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]any{"error": err.Error()})
		}

		metrics = append(metrics, m)
	}

	if err = rows.Err(); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]any{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, metrics)
}

const N = 10_000_000

func Sequential() float64 {
	arr := make([]int, N)

	t0 := time.Now()

	for i := range N {
		arr[i] = rand.Intn(200)
	}

	t1 := time.Since(t0)

	return t1.Seconds()
}

func Parallel() float64 {
	var wg sync.WaitGroup
	arr := make([]int, N)

	tp0 := time.Now()

	wg.Go(func() {
		for i := range N {
			arr[i] = rand.Intn(200)
		}
	})

	wg.Wait()

	tp1 := time.Since(tp0)

	fmt.Println(arr)

	return tp1.Seconds()
}

func ParallelWithThreads() float64 {
	var wg sync.WaitGroup
	totalThreads := 12
	chunckSize := N / totalThreads
	arr := make([]int, N)

	tp0 := time.Now()
	for i := range totalThreads {
		wg.Go(func() {
			start := i * chunckSize
			end := start + chunckSize
			if (i + 1) == totalThreads {
				end = N
			}

			for j := start; j < end; j++ {
				arr[j] = rand.Intn(200)
			}
		})
	}

	wg.Wait()

	tp1 := time.Since(tp0)

	return tp1.Seconds()
}

func (h *Handler) HandlerSequential(c *echo.Context) error {
	timeSeq := Sequential()

	c.Set("complexity", "O(n)")
	c.Set("n", N)

	return c.JSON(http.StatusOK, timeSeq)
}

func (h *Handler) HandlerParallel(c *echo.Context) error {
	timePar := Parallel()

	c.Set("complexity", "O(n)")
	c.Set("n", N)

	return c.JSON(http.StatusOK, timePar)
}

func (h *Handler) HandlerParallelWithThreads(c *echo.Context) error {
	timePar2 := ParallelWithThreads()

	c.Set("complexity", "O(n)")
	c.Set("n", N)

	return c.JSON(http.StatusOK, timePar2)
}

func (h *Handler) HandlerParallelMetrics(c *echo.Context) error {
	timeSeq := Sequential()

	timePar1 := Parallel()

	timePar2 := ParallelWithThreads()

	speedUp := timeSeq / timePar1
	speedUp2 := timeSeq / timePar2

	porcentParallel := (1 - (1 / speedUp)) / (1 - (1 / 2.0))
	porcentParallel2 := (1 - (1 / speedUp2)) / (1 - (1 / 12.0))

	return c.JSON(http.StatusOK, map[string]float64{
		"timeSeq":          timeSeq,
		"timePar1":         timePar1,
		"timePar2":         timePar2,
		"speedUp":          speedUp,
		"speedUp2":         speedUp2,
		"porcentParallel":  porcentParallel,
		"porcentParallel2": porcentParallel2,
	})
}
