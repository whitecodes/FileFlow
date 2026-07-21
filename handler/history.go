package handler

import (
	"net/http"

	"FileFlow/model"
	"FileFlow/service"

	"github.com/labstack/echo/v4"
)

func ListHistory(c echo.Context) error {
	records, err := service.ListHistory(100)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "failed to list history"})
	}
	if records == nil {
		records = []model.History{}
	}
	return c.JSON(http.StatusOK, records)
}
