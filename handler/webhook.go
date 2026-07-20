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

		// Try each search directory
		for _, dir := range cfg.SearchDirs {
			result, err := service.ProcessFile(dir, req.FileName)
			if err != nil {
				log.Printf("[webhook] process error in %s: %v", dir, err)
				continue
			}
			if result != nil {
				log.Printf("[webhook] processed: %s -> %s (rule=%q)", result.SrcPath, result.DstPath, result.Rule.Name)
				return c.JSON(http.StatusOK, map[string]string{
					"status":    "ok",
					"src_path":  result.SrcPath,
					"dst_path":  result.DstPath,
					"rule_name": result.Rule.Name,
				})
			}
		}

		return c.JSON(http.StatusOK, map[string]string{
			"status": "ok",
		})
	}
}
