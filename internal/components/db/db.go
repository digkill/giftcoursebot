package db

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"os"
)

type DB struct {
	Conn *sql.DB
}

func InitDB() *DB {
	// Пример: "user:password@tcp(127.0.0.1:3306)/giftbot"
	dsn := os.Getenv("MYSQL_DSN")
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		panic(err)
	}

	db.Exec(`CREATE TABLE IF NOT EXISTS users (
        chat_id BIGINT PRIMARY KEY,
        start_date DATETIME,
        last_lesson_sent INT DEFAULT 0
    )`)

	db.Exec(`CREATE TABLE IF NOT EXISTS feedback (
        id INT AUTO_INCREMENT PRIMARY KEY,
        chat_id BIGINT,
        feedback TEXT,
        created_at DATETIME
    )`)

	return &DB{Conn: db}
}

func (db *DB) Close() {
	db.Conn.Close()
}
