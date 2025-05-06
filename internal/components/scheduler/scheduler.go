package scheduler

import (
	"github.com/digkill/giftcoursebot/internal/components/db"
	"github.com/digkill/giftcoursebot/internal/models"
	"github.com/sirupsen/logrus"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func StartScheduler(bot *tgbotapi.BotAPI, user *models.User, lesson *models.Lesson) {
	ticker := time.NewTicker(1 * time.Hour)
	for {
		<-ticker.C
		users := user.GetAllUsers()
		for _, user := range users {
			days := int(time.Since(user.UserModel.StartDate).Hours() / 24)
			if days >= 0 && days < len(lessons) && user.LastLessonSent < days {
				msg := tgbotapi.NewMessage(user.ChatID, lessons[days])
				msg.ReplyMarkup = FeedbackButtons()
				_, err := bot.Send(msg)
				if err != nil {
					logrus.Warnf("[Scheduler][StartScheduler] Error: %v", err)
				}
				db.UpdateLastLesson(user.ChatID, days)
			}
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
