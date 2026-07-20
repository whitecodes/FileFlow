package handler

import (
	"log"
	"net/http"
	"strconv"

	"FileFlow/model"
	"FileFlow/service"

	"github.com/labstack/echo/v4"
)

type createRuleRequest struct {
	Name           string `json:"name"`
	Pattern        string `json:"pattern"`
	TargetDir      string `json:"target_dir"`
	RenameTemplate string `json:"rename_template"`
	Priority       int    `json:"priority"`
	Enabled        bool   `json:"enabled"`
}

func CreateRule(c echo.Context) error {
	var req createRuleRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid request body"})
	}
	if req.Name == "" || req.Pattern == "" || req.TargetDir == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "name, pattern and target_dir are required"})
	}
	if req.RenameTemplate == "" {
		req.RenameTemplate = "{title}.{ext}"
	}

	rule := &model.Rule{
		Name:           req.Name,
		Pattern:        req.Pattern,
		TargetDir:      req.TargetDir,
		RenameTemplate: req.RenameTemplate,
		Priority:       req.Priority,
		Enabled:        req.Enabled,
	}

	id, err := service.CreateRule(rule)
	if err != nil {
		log.Printf("[rule] create error: %v", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "failed to create rule"})
	}

	created, err := service.GetRule(id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "rule created but failed to fetch"})
	}

	return c.JSON(http.StatusCreated, created)
}

func ListRules(c echo.Context) error {
	rules, err := service.ListRules()
	if err != nil {
		log.Printf("[rule] list error: %v", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "failed to list rules"})
	}
	if rules == nil {
		rules = []model.Rule{}
	}
	return c.JSON(http.StatusOK, rules)
}

func GetRule(c echo.Context) error {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid id"})
	}

	rule, err := service.GetRule(id)
	if err != nil {
		log.Printf("[rule] get error: %v", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "failed to get rule"})
	}
	if rule == nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "rule not found"})
	}
	return c.JSON(http.StatusOK, rule)
}

func DeleteRule(c echo.Context) error {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid id"})
	}

	// Check existence first
	rule, err := service.GetRule(id)
	if err != nil {
		log.Printf("[rule] get error: %v", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "failed to get rule"})
	}
	if rule == nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "rule not found"})
	}

	if err := service.DeleteRule(id); err != nil {
		log.Printf("[rule] delete error: %v", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "failed to delete rule"})
	}

	return c.JSON(http.StatusOK, map[string]string{"status": "deleted"})
}
