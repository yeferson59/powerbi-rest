package middleware

import (
	"github.com/labstack/echo/v5"
	"github.com/labstack/echo/v5/middleware"
)

func (Middleware) Logger() echo.MiddlewareFunc {
	return middleware.RequestLogger()
}
