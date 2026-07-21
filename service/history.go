package service

import (
	"FileFlow/db"
	"FileFlow/model"
)

const maxHistoryRecords = 100

// RecordHistory 写入一条处理记录，并自动裁剪超过 100 条的旧记录。
func RecordHistory(fileName, event, ruleName, srcPath, dstPath, status, errorMsg string) error {
	_, err := db.DB.Exec(
		`INSERT INTO history (file_name, event, rule_name, src_path, dst_path, status, error_msg)
		 VALUES (?, ?, ?, ?, ?, ?, ?)`,
		fileName, event, ruleName, srcPath, dstPath, status, errorMsg,
	)
	if err != nil {
		return err
	}
	// Trim old records, keep only the latest maxHistoryRecords
	_, err = db.DB.Exec(
		`DELETE FROM history WHERE id NOT IN (
			SELECT id FROM history ORDER BY id DESC LIMIT ?
		)`, maxHistoryRecords,
	)
	return err
}

// ListHistory 返回最近 limit 条处理记录，按时间倒序。
func ListHistory(limit int) ([]model.History, error) {
	rows, err := db.DB.Query(
		`SELECT id, file_name, event, rule_name, src_path, dst_path, status, error_msg, created_at
		 FROM history ORDER BY id DESC LIMIT ?`, limit,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var records []model.History
	for rows.Next() {
		var h model.History
		if err := rows.Scan(&h.ID, &h.FileName, &h.Event, &h.RuleName,
			&h.SrcPath, &h.DstPath, &h.Status, &h.ErrorMsg, &h.CreatedAt); err != nil {
			return nil, err
		}
		records = append(records, h)
	}
	return records, rows.Err()
}
