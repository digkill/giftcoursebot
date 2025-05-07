package models

import (
	"database/sql"
	"github.com/sirupsen/logrus"
	"time"
)

type Feedback struct {
	ID        int
	ChatID    int64
	Feedback  string
	CreatedAt time.Time
}

type FeedbackModel struct {
	DB *sql.DB
}

// Save сохраняет отзыв пользователя
func (f *FeedbackModel) Save(chatID int64, feedbackText string) error {
	_, err := f.DB.Exec(
		"INSERT INTO user_feedback (chat_id, feedback, created_at) VALUES (?, ?, ?)",
		chatID, feedbackText, time.Now(),
	)
	if err != nil {
		logrus.Errorf("FeedbackModel.Save error: %v", err)
	}
	return err
}

// GetAll возвращает все отзывы
func (f *FeedbackModel) GetAll() ([]Feedback, error) {
	rows, err := f.DB.Query("SELECT id, chat_id, feedback, created_at FROM user_feedback ORDER BY created_at DESC")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var feedbacks []Feedback
	for rows.Next() {
		var fb Feedback
		if err := rows.Scan(&fb.ID, &fb.ChatID, &fb.Feedback, &fb.CreatedAt); err == nil {
			feedbacks = append(feedbacks, fb)
		}
	}
	return feedbacks, nil
}
