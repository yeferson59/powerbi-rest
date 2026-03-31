package middleware

import (
	"github.com/labstack/echo/v5"
	"github.com/labstack/echo/v5/middleware"
)

func (Middleware) Cors() echo.MiddlewareFunc {
	return middleware.CORS("*")
}
