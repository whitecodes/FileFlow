package db

import (
	"database/sql"
	"log"

	_ "modernc.org/sqlite"
)

var DB *sql.DB

func Init(dsn string) {
	var err error
	DB, err = sql.Open("sqlite", dsn)
	if err != nil {
		log.Fatalf("[db] open error: %v", err)
	}
	if err = DB.Ping(); err != nil {
		log.Fatalf("[db] ping error: %v", err)
	}
	log.Println("[db] connected")
	migrate()
}

func migrate() {
	queries := []string{
		`CREATE TABLE IF NOT EXISTS rules (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			name TEXT NOT NULL,
			pattern TEXT NOT NULL,
			target_dir TEXT NOT NULL,
			rename_template TEXT NOT NULL DEFAULT '{title}.{ext}',
			priority INTEGER NOT NULL DEFAULT 0,
			enabled INTEGER NOT NULL DEFAULT 1,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
		)`,
		`CREATE TABLE IF NOT EXISTS history (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			file_name TEXT NOT NULL,
			event TEXT DEFAULT '',
			rule_name TEXT DEFAULT '',
			src_path TEXT DEFAULT '',
			dst_path TEXT DEFAULT '',
			status TEXT NOT NULL DEFAULT '',
			error_msg TEXT DEFAULT '',
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP
		)`,
	}
	for _, q := range queries {
		if _, err := DB.Exec(q); err != nil {
			log.Fatalf("[db] migrate error: %v", err)
		}
	}
	log.Println("[db] migration done")
}

func Close() {
	if DB != nil {
		DB.Close()
	}
}
