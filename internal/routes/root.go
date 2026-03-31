package routes

import (
	"github.com/labstack/echo/v5"
	"github.com/yeferson59/powerbi-rest/internal/handlers"
	"github.com/yeferson59/powerbi-rest/internal/middleware"
)

type Route struct {
	echo       *echo.Echo
	handler    handlers.Handler
	middleware middleware.Middleware
}

func New(echo *echo.Echo, handler handlers.Handler, middleware middleware.Middleware) Route {
	return Route{
		echo:       echo,
		handler:    handler,
		middleware: middleware,
	}
}

func (r Route) Init() error {
	r.echo.Use(r.middleware.RequestID())
	r.echo.Use(r.middleware.Cors())
	r.echo.Use(r.middleware.Logger())
	r.echo.Use(r.middleware.Error())
	r.echo.Use(r.middleware.Metrics())

	r.echo.GET("/", r.handler.HandlerRoot)
	r.echo.GET("/o1", r.handler.HandlerO1)
	r.echo.GET("/on", r.handler.HandlerOn)
	r.echo.GET("/onlogn", r.handler.HandlerONLogN)
	r.echo.GET("/on2", r.handler.HandlerON2)
	r.echo.GET("/o2n", r.handler.HandlerO2N)
	r.echo.GET("/summary", r.handler.HandlerSummary)

	return nil
}
