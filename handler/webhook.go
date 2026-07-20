package handler

import (
	"log"
	"net/http"

	"FileFlow/service"

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

	rule, err := service.MatchRule(req.FileName)
	if err != nil {
		log.Printf("[webhook] match error: %v", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "failed to match rule",
		})
	}
	if rule != nil {
		log.Printf("[webhook] matched rule=%q (id=%d) pattern=%q target=%q",
			rule.Name, rule.ID, rule.Pattern, rule.TargetDir)
	}

	return c.JSON(http.StatusOK, map[string]string{
		"status": "ok",
	})
}
