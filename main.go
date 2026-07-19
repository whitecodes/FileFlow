package main

import (
	"fmt"

	"FileFlow/config"
	"FileFlow/handler"

	"github.com/labstack/echo/v4"
)

func main() {
	cfg := config.Load()

	e := echo.New()

	e.POST("/api/webhook", handler.Webhook)

	e.GET("/api/rules", func(c echo.Context) error {
		return c.JSON(200, []interface{}{})
	})

	e.Logger.Fatal(e.Start(fmt.Sprintf(":%d", cfg.Port)))
}
