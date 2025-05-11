package scheduler

import (
	"time"

	"github.com/digkill/giftcoursebot/internal/models"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/sirupsen/logrus"
)

func StartScheduler(bot *tgbotapi.BotAPI, userModel *models.UserModel, lessonModel *models.LessonModel) {
	ticker := time.NewTicker(1 * time.Hour)

	for {
		<-ticker.C

		users := userModel.GetAllUsers()
		for _, user := range users {
			daysSinceStart := int(time.Since(user.StartDate).Hours() / 24)
			daysSinceStart = 1
			// Получаем все уроки, которые пользователь уже получил
			sentLessonIDs := lessonModel.GetSentLessonIDs(user.ChatID)

			// Получаем урок, который соответствует текущему дню
			nextLesson := lessonModel.GetLessonByDay(daysSinceStart)

			if nextLesson == nil {
				continue // Все уроки пройдены
			}

			// Проверяем, отправлялся ли уже этот урок
			if contains(sentLessonIDs, int(nextLesson.ID)) {
				continue
			}

			// Отправляем урок
			msg := tgbotapi.NewMessage(user.ChatID, nextLesson.Content)
			msg.ReplyMarkup = FeedbackButtons()

			if _, err := bot.Send(msg); err != nil {
				logrus.Warnf("[Scheduler][StartScheduler] Error sending to %d: %v", user.ChatID, err)
				continue
			}

			// Сохраняем факт отправки
			lessonModel.MarkLessonSent(user.ChatID, int(nextLesson.ID))
		}
	}
}

func FeedbackButtons() tgbotapi.InlineKeyboardMarkup {
	return tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("👍", "feedback_good"),
			tgbotapi.NewInlineKeyboardButtonData("👎", "feedback_bad"),
		),
	)
}

func contains(slice []int, value int) bool {
	for _, v := range slice {
		if v == value {
			return true
		}
	}
	return false
}
