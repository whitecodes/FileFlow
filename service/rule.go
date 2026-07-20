package service

import (
	"database/sql"
	"FileFlow/db"
	"FileFlow/model"
)

func CreateRule(r *model.Rule) (int64, error) {
	res, err := db.DB.Exec(
		`INSERT INTO rules (name, pattern, target_dir, rename_template, priority, enabled)
		 VALUES (?, ?, ?, ?, ?, ?)`,
		r.Name, r.Pattern, r.TargetDir, r.RenameTemplate, r.Priority, r.Enabled,
	)
	if err != nil {
		return 0, err
	}
	return res.LastInsertId()
}

func ListRules() ([]model.Rule, error) {
	rows, err := db.DB.Query(
		`SELECT id, name, pattern, target_dir, rename_template, priority, enabled, created_at, updated_at
		 FROM rules ORDER BY priority DESC, id ASC`,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var rules []model.Rule
	for rows.Next() {
		var r model.Rule
		if err := rows.Scan(&r.ID, &r.Name, &r.Pattern, &r.TargetDir,
			&r.RenameTemplate, &r.Priority, &r.Enabled, &r.CreatedAt, &r.UpdatedAt); err != nil {
			return nil, err
		}
		rules = append(rules, r)
	}
	return rules, rows.Err()
}

func GetRule(id int64) (*model.Rule, error) {
	var r model.Rule
	err := db.DB.QueryRow(
		`SELECT id, name, pattern, target_dir, rename_template, priority, enabled, created_at, updated_at
		 FROM rules WHERE id = ?`, id,
	).Scan(&r.ID, &r.Name, &r.Pattern, &r.TargetDir,
		&r.RenameTemplate, &r.Priority, &r.Enabled, &r.CreatedAt, &r.UpdatedAt)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &r, nil
}

func DeleteRule(id int64) error {
	_, err := db.DB.Exec(`DELETE FROM rules WHERE id = ?`, id)
	return err
}

func UpdateRule(id int64, r *model.Rule) error {
	res, err := db.DB.Exec(
		`UPDATE rules
		 SET name = ?, pattern = ?, target_dir = ?, rename_template = ?, priority = ?, enabled = ?,
		     updated_at = CURRENT_TIMESTAMP
		 WHERE id = ?`,
		r.Name, r.Pattern, r.TargetDir, r.RenameTemplate, r.Priority, r.Enabled, id,
	)
	if err != nil {
		return err
	}
	n, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if n == 0 {
		return sql.ErrNoRows
	}
	return nil
}
