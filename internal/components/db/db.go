package db

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"github.com/sirupsen/logrus"
)

var Conn *sql.DB

func Init(dsn string) {
	var err error
	Conn, err = sql.Open("mysql", dsn)
	if err != nil {
		logrus.Fatal("Failed to connect to DB:", err)
	}

	if err := Conn.Ping(); err != nil {
		logrus.Fatal("DB unreachable:", err)
	}

	logrus.Info("Connected to database")
}
