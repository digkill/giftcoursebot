package models

import (
	"database/sql"
	"fmt"
	"github.com/sirupsen/logrus"
	"time"
)

type User struct {
	ChatID         int64          `db:"id"`
	StartDate      time.Time      `db:"start_date"`
	LastLessonSent int            `db:"last_lessonSent"`
	CreatedAt      sql.NullString `db:"created_at"`
	UpdatedAt      sql.NullString `db:"updated_at"`
}

type UserModel struct {
	DB *sql.DB
}

func (u *UserModel) GetAllUsers() []User {
	rows, err := u.DB.Query("SELECT chat_id, start_date FROM users")
	if err != nil {
		logrus.Error("GetAllUsers error:", err)
		return nil
	}
	defer rows.Close()

	var users []User
	for rows.Next() {
		var user User
		err := rows.Scan(&user.ChatID, &user.StartDate)
		if err != nil {
			logrus.Error("Scan user:", err)
			continue
		}
		users = append(users, user)
	}
	return users
}

func (u *UserModel) RegisterUser(chatID int64) {
	now := time.Now().Format("2006-01-02 15:04:05")
	_, err := u.DB.Exec("INSERT IGNORE INTO users (chat_id, start_date, last_lesson_sent) VALUES (?, ?, 0)", chatID, now)
	if err != nil {
		fmt.Println("RegisterUser error:", err)
	}
}

func (u *UserModel) UpdateLastLesson(chatID int64, lessonNum int) {
	u.DB.Exec("UPDATE users SET last_lesson_sent = ? WHERE chat_id = ?", lessonNum, chatID)
}

func (u *UserModel) SaveFeedback(chatID int64, feedback string) {
	now := time.Now().Format("2006-01-02 15:04:05")
	u.DB.Exec("INSERT INTO feedback (chat_id, feedback, created_at) VALUES (?, ?, ?)", chatID, feedback, now)
}
