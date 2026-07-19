package handler

import (
	"log"
	"net/http"

	"github.com/labstack/echo/v4"
)

type WebhookRequest struct {
	FileName string `json:"file_name"`
	Event    string `json:"event"`
}

func Webhook(c echo.Context) error {
	var req WebhookRequest
	if err := c.Bind(&req); err != nil {
		log.Printf("[webhook] bind error: %v", err)
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "invalid request body",
		})
	}

	log.Printf("[webhook] event=%s file=%s", req.Event, req.FileName)

	return c.JSON(http.StatusOK, map[string]string{
		"status": "ok",
	})
}
