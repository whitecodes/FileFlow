package handler

import (
	"log"
	"net/http"

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
		_ = service.RecordHistory(req.FileName, req.Event, "", "", "", "received", "")

		// Try each search directory
		for _, dir := range cfg.SearchDirs {
			result, err := service.ProcessFile(dir, req.FileName)
			if err != nil {
				log.Printf("[webhook] process error in %s: %v", dir, err)
				_ = service.RecordHistory(req.FileName, req.Event, "", "", "", "error", err.Error())
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

		_ = service.RecordHistory(req.FileName, req.Event, "", "", "", "no_match", "")
		return c.JSON(http.StatusOK, map[string]string{
			"status": "ok",
		})
	}
}
