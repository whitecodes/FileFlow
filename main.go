package main

import (
	"fmt"

	"FileFlow/config"
	"FileFlow/db"
	"FileFlow/handler"

	"github.com/labstack/echo/v4"
)

func main() {
	cfg := config.Load()

	db.Init(cfg.DBPath)
	defer db.Close()

	e := echo.New()

	handler.RegisterStatic(e)

	e.GET("/api/history", handler.ListHistory)

	e.POST("/api/webhook", handler.Webhook(cfg))

	e.GET("/api/rules", handler.ListRules)
	e.POST("/api/rules", handler.CreateRule)
	e.GET("/api/rules/:id", handler.GetRule)
	e.PUT("/api/rules/:id", handler.UpdateRule)
	e.DELETE("/api/rules/:id", handler.DeleteRule)

	e.Logger.Fatal(e.Start(fmt.Sprintf(":%d", cfg.Port)))
}
