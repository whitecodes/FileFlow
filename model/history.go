package model

import "time"

type History struct {
	ID        int64     `json:"id"`
	FileName  string    `json:"file_name"`
	Event     string    `json:"event,omitempty"`
	RuleName  string    `json:"rule_name,omitempty"`
	SrcPath   string    `json:"src_path,omitempty"`
	DstPath   string    `json:"dst_path,omitempty"`
	Status    string    `json:"status"`       // matched / no_match / error
	ErrorMsg  string    `json:"error_msg,omitempty"`
	CreatedAt time.Time `json:"created_at"`
}
