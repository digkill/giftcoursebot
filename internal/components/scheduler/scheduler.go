package scheduler

import (
	logger "github.com/digkill/giftcoursebot/internal/components/logger"
	"github.com/digkill/giftcoursebot/internal/helpers"
	"os"
	"path/filepath"
	"time"

	"github.com/digkill/giftcoursebot/internal/models"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func StartScheduler(bot *tgbotapi.BotAPI, userModel *models.UserModel, lessonModel *models.LessonModel, lg *logger.Log) {
	ticker := time.NewTicker(1 * time.Minute)
	daysSinceStart := 0
	for {
		<-ticker.C

		users := userModel.GetAllUsers()
		for _, user := range users {
			lg.Logger.Infof("Проверяем пользователя %d", user.ChatID)
			// daysSinceStart := int(time.Since(user.StartDate).Hours() / 24)

			daysSinceStart = daysSinceStart + 1

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

			// Отправляем урок
			msgTitle := tgbotapi.NewMessage(user.ChatID, "🎓 "+nextLesson.Title)
			msgTitle.ParseMode = "Markdown"
			bot.Send(msgTitle)

			msgContent := tgbotapi.NewMessage(user.ChatID, nextLesson.Content)
			msgContent.ParseMode = "Markdown"
			bot.Send(msgContent)

			imageDir := "./assets/images/"
			imageOutputPath := filepath.Join(imageDir, nextLesson.Image)

			fileImage, err := os.Open(imageOutputPath)
			if err != nil {
				lg.Logger.Warnf("failed to open file: %w", err)

			}
			defer fileImage.Close()

			err = helpers.SendMedia(fileImage, bot, user.ChatID, nextLesson.Caption)
			if err != nil {
				lg.Logger.Warnf("failed to send photo: %w", err)
			}

			image2OutputPath := filepath.Join(imageDir, nextLesson.Image2)

			fileImage2, err := os.Open(image2OutputPath)
			if err != nil {
				lg.Logger.Warnf("failed to open file2: %w", err)

			}
			defer fileImage2.Close()

			err = helpers.SendMedia(fileImage2, bot, user.ChatID, nextLesson.Caption2)
			if err != nil {
				lg.Logger.Warnf("failed to send photo2: %w", err)
			}

			reply := tgbotapi.NewMessage(user.ChatID, nextLesson.Link)
			// reply.ParseMode = "Markdown"
			reply.ReplyMarkup = FeedbackButtons()
			_, err = bot.Send(reply)
			if err != nil {
				lg.Logger.Warnf("failed to send link: %w", err)
			}

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
