package handler

import (
	"embed"
	"io/fs"
	"net/http"

	"github.com/labstack/echo/v4"
)

//go:embed views/*
var staticFS embed.FS

func RegisterStatic(e *echo.Echo) {
	sub, _ := fs.Sub(staticFS, "views")
	h := http.FileServer(http.FS(sub))

	// Remove trailing slash redirect for /ui prefix
	e.Pre(func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			p := c.Request().URL.Path
			if len(p) >= 3 && p[:3] == "/ui" {
				// Let the file server handle routing directly
				return next(c)
			}
			return next(c)
		}
	})

	e.GET("/ui", func(c echo.Context) error {
		r := c.Request()
		r.URL.Path = "/"
		h.ServeHTTP(c.Response().Writer, r)
		return nil
	})

	e.GET("/ui/*", func(c echo.Context) error {
		r := c.Request()
		r.URL.Path = "/" + c.Param("*")
		h.ServeHTTP(c.Response().Writer, r)
		return nil
	})
}
