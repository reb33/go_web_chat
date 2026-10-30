package httpserver

import (
	"log/slog"
	"net/http"
	"text/template"
	"web_chat/internal/transport/httpserver/handlers"
	"web_chat/internal/transport/httpserver/ports"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	slogecho "github.com/samber/slog-echo"
)

func NewRouter(
	logger *slog.Logger,
	authUseCase ports.AuthUseCase,
) *echo.Echo {
	e := echo.New()
	e.Renderer = &echo.TemplateRenderer{
		Template: template.Must(template.ParseGlob("web/templates/*.html")),
	}

	e.Use(slogecho.New(logger))
	e.Use(middleware.Recover())

	e.GET("/ping", func(c echo.Context) error { return c.String(http.StatusOK, "pong") })

	e.GET("/", func(c echo.Context) error {
		return c.Render(http.StatusOK, "index", nil)
	})

	e.POST("/login", handlers.Login(authUseCase, logger))

	return e
}
