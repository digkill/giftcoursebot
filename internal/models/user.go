package models

import (
	"database/sql"
	"fmt"
	"time"
)

type UserModel struct {
	ChatID         int64          `db:"id"`
	StartDate      time.Time      `db:"start_date"`
	LastLessonSent int            `db:"last_lessonSent"`
	CreatedAt      sql.NullString `db:"created_at"`
	UpdatedAt      sql.NullString `db:"updated_at"`
}

type User struct {
	db *sql.DB
}

func (user *User) RegisterUser(chatID int64) {
	now := time.Now().Format("2006-01-02 15:04:05")
	_, err := user.db.Exec("INSERT IGNORE INTO users (chat_id, start_date, last_lesson_sent) VALUES (?, ?, 0)", chatID, now)
	if err != nil {
		fmt.Println("RegisterUser error:", err)
	}
}

func (user *User) GetAllUsers() []UserModel {
	rows, err := user.db.Query("SELECT chat_id, start_date, last_lesson_sent FROM users")
	if err != nil {
		fmt.Println("GetAllUsers error:", err)
		return nil
	}
	defer rows.Close()

	var users []User
	for rows.Next() {
		var u User
		var start time.Time
		rows.Scan(&u.ChatID, &start, &u.LastLessonSent)
		u.StartDate = start
		users = append(users, u)
	}
	return users
}

func (user *User) UpdateLastLesson(chatID int64, lessonNum int) {
	user.db.Exec("UPDATE users SET last_lesson_sent = ? WHERE chat_id = ?", lessonNum, chatID)
}

func (user *User) SaveFeedback(chatID int64, feedback string) {
	now := time.Now().Format("2006-01-02 15:04:05")
	user.db.Exec("INSERT INTO feedback (chat_id, feedback, created_at) VALUES (?, ?, ?)", chatID, feedback, now)
}
