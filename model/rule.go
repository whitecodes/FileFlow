package model

import "time"

type Rule struct {
	ID             int64     `json:"id"`
	Name           string    `json:"name"`
	Pattern        string    `json:"pattern"`
	TargetDir      string    `json:"target_dir"`
	RenameTemplate string    `json:"rename_template"`
	Priority       int       `json:"priority"`
	Enabled        bool      `json:"enabled"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}
