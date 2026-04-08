package handlers

import (
	"math/rand"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v5"
	"github.com/yeferson59/powerbi-rest/internal/metrics"
)

type Handler struct {
	metricsStore metrics.Store
}

func New(metricsStore metrics.Store) *Handler {
	return &Handler{
		metricsStore: metricsStore,
	}
}

func (Handler) HandlerRoot(c *echo.Context) error {
	c.Response().Header().Set("Cache-Control", "no-store, no-cache, must-revalidate")
	c.Response().Header().Set("Pragma", "no-cache")
	c.Response().Header().Set("Expires", "0")

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

func (h *Handler) HandlerSummary(c *echo.Context) error {
	c.Response().Header().Set("Cache-Control", "no-store, no-cache, must-revalidate")
	c.Response().Header().Set("Pragma", "no-cache")
	c.Response().Header().Set("Expires", "0")

	records, err := h.metricsStore.List(c.Request().Context())
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]any{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, records)
}
