package handler

import (
	"log"
	"net/http"
	"strings"

	"FileFlow/config"
	"FileFlow/service"

	"github.com/labstack/echo/v4"
)

type WebhookRequest struct {
	FileName string `json:"file_name"`
	Event    string `json:"event"`
}

func Webhook(cfg *config.Config) echo.HandlerFunc {
	return func(c echo.Context) error {
		var req WebhookRequest
		if err := c.Bind(&req); err != nil {
			log.Printf("[webhook] bind error: %v", err)
			return c.JSON(http.StatusBadRequest, map[string]string{
				"error": "invalid request body",
			})
		}

		log.Printf("[webhook] event=%s file=%s", req.Event, req.FileName)

		var lastErr error
		for _, dir := range cfg.SearchDirs {
			results, err := service.ProcessFile(dir, req.FileName)
			if err != nil {
				lastErr = err
				continue
			}
			if len(results) > 0 {
				for _, r := range results {
					log.Printf("[webhook] processed: %s -> %s (rule=%q)", r.SrcPath, r.DstPath, r.Rule.Name)
					_ = service.RecordHistory(req.FileName, req.Event, r.Rule.Name, r.SrcPath, r.DstPath, "matched", "")
				}
				return c.JSON(http.StatusOK, map[string]string{
					"status": "ok",
				})
			}
		}

		if lastErr != nil {
			msg := lastErr.Error()
			if strings.Contains(msg, "not found") {
				_ = service.RecordHistory(req.FileName, req.Event, "", "", "", "not_found", "")
				return c.JSON(http.StatusOK, map[string]string{"status": "ok"})
			}
			_ = service.RecordHistory(req.FileName, req.Event, "", "", "", "error", msg)
			return c.JSON(http.StatusOK, map[string]string{"status": "ok"})
		}

		_ = service.RecordHistory(req.FileName, req.Event, "", "", "", "no_match", "")
		return c.JSON(http.StatusOK, map[string]string{
			"status": "ok",
		})
	}
}
