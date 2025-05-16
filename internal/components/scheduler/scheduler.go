package scheduler

import (
	"fmt"
	logger "github.com/digkill/giftcoursebot/internal/components/logger"
	"github.com/digkill/giftcoursebot/internal/helpers"
	"time"

	"github.com/digkill/giftcoursebot/internal/models"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func StartScheduler(bot *tgbotapi.BotAPI, userModel *models.UserModel, lessonModel *models.LessonModel, lg *logger.Log) {
	ticker := time.NewTicker(1 * time.Hour)

	for {
		<-ticker.C

		users := userModel.GetAllUsers()
		for _, user := range users {
			lg.Logger.Infof("Проверяем пользователя %d", user.ChatID)
			daysSinceStart := int(time.Since(user.StartDate).Hours() / 24)

			lg.Logger.Infof("Прошло дней с начала: %d", daysSinceStart)
			// Получаем все уроки, которые пользователь уже получил
			sentLessonIDs := lessonModel.GetSentLessonIDs(user.ChatID)

			// Получаем урок, который соответствует текущему дню
			nextLesson := lessonModel.GetLessonByDay(daysSinceStart)

			if nextLesson == nil {
				lg.Logger.Warnf("Урок на день %d не найден", daysSinceStart)
				continue // Все уроки пройдены
			}

			// Проверяем, отправлялся ли уже этот урок
			if contains(sentLessonIDs, int(nextLesson.ID)) {
				lg.Logger.Warnf("Урок %d уже отправлен пользователю %d", nextLesson.ID, user.ChatID)
				continue
			}

			helpers.SendLesson(bot, lg, nextLesson, user.ChatID, fmt.Sprintf("🤗 Держи новое занятие! № %d", daysSinceStart))

			// Сохраняем факт отправки
			lg.Logger.Infof("Урок %d отправлен пользователю %d", nextLesson.ID, user.ChatID)

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
