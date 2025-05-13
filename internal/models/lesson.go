package models

import (
	"database/sql"
	"github.com/sirupsen/logrus"
)

type Lesson struct {
	ID        int64          `db:"id"`
	Title     string         `db:"title"`
	Content   string         `db:"content"`
	Image     string         `db:"image"`
	Image2    string         `db:"image2"`
	Caption   string         `db:"caption"`
	Caption2  string         `db:"caption2"`
	Link      string         `db:"link"`
	Slug      sql.NullString `db:"slug"`
	CreatedAt sql.NullString `db:"created_at"`
	UpdatedAt sql.NullString `db:"updated_at"`
}

type LessonModel struct {
	DB *sql.DB
}

func (l *LessonModel) GetLessonByDay(day int) *Lesson {
	row := l.DB.QueryRow("SELECT id, title, content, image, image2, caption, caption2, slug, link, created_at, updated_at FROM lessons WHERE day_number = ? LIMIT 1", day)
	var lesson Lesson
	err := row.Scan(&lesson.ID, &lesson.Title, &lesson.Content, &lesson.Image, &lesson.Image2, &lesson.Caption, &lesson.Caption2, &lesson.Slug, &lesson.Link, &lesson.CreatedAt, &lesson.UpdatedAt)
	if err != nil {
		logrus.Warn("Lesson not found for day:", day)
		return nil
	}
	return &lesson
}

func (l *LessonModel) GetSentLessonIDs(chatID int64) []int {
	rows, err := l.DB.Query("SELECT lesson_id FROM user_lessons WHERE user_id = ?", chatID)
	if err != nil {
		logrus.Error("GetSentLessonIDs error:", err)
		return nil
	}
	defer rows.Close()

	var ids []int
	for rows.Next() {
		var id int
		rows.Scan(&id)
		ids = append(ids, id)
	}
	return ids
}

func (l *LessonModel) MarkLessonSent(chatID int64, lessonID int) {
	_, err := l.DB.Exec("INSERT INTO user_lessons (user_id, lesson_id, sent_at) VALUES (?, ?, NOW())", chatID, lessonID)
	if err != nil {
		logrus.Error("MarkLessonSent error:", err)
	}
}
