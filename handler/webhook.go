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
			result, err := service.ProcessFile(dir, req.FileName)
			if err != nil {
				lastErr = err
				continue
			}
			if result != nil {
				log.Printf("[webhook] processed: %s -> %s (rule=%q)", result.SrcPath, result.DstPath, result.Rule.Name)
				_ = service.RecordHistory(req.FileName, req.Event, result.Rule.Name, result.SrcPath, result.DstPath, "matched", "")
				return c.JSON(http.StatusOK, map[string]string{
					"status":    "ok",
					"src_path":  result.SrcPath,
					"dst_path":  result.DstPath,
					"rule_name": result.Rule.Name,
				})
			}
		}

		// Determine the real outcome
		if lastErr != nil {
			msg := lastErr.Error()
			if strings.Contains(msg, "file not found") {
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
