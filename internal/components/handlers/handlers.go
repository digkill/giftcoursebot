package handlers

import (
	"github.com/digkill/giftcoursebot/internal/components/db"
	"github.com/digkill/giftcoursebot/internal/components/logger"
	"github.com/digkill/giftcoursebot/internal/helpers"
	"github.com/digkill/giftcoursebot/internal/models"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"os"
	"path/filepath"
	"time"
)

func HandleMessage(bot *tgbotapi.BotAPI, userModel *models.UserModel, lessonModel *models.LessonModel, msg *tgbotapi.Message, lg *logger.Log) {
	switch msg.Text {
	case "/start":
		// Регистрируем пользователя
		userModel.DB.Exec("INSERT IGNORE INTO users (chat_id, start_date) VALUES (?, ?)", msg.Chat.ID, time.Now())

		welcomeMessage := `🎉 *Добро пожаловать в наш подарок-курс!* 🎁
Ты только что открыл дверь в маленькое путешествие знаний, вдохновения и радости! 💡
📚 *Каждый день* тебя ждёт новое занятие, только полезное и с заботой о тебе.
👉 Первый урок уже на подходе. А если что — я всегда рядом!
*Удачи тебе! Пусть этот курс принесёт пользу и удовольствие!* 🌟
`

		message := tgbotapi.NewMessage(msg.Chat.ID, welcomeMessage)
		message.ParseMode = "Markdown"
		bot.Send(message)

		// Получаем первый урок (day_number = 0)
		lesson := lessonModel.GetLessonByDay(0)
		if lesson == nil {
			bot.Send(tgbotapi.NewMessage(msg.Chat.ID, "Урок не найден. Обратитесь в поддержку."))
			return
		}

		// Отправляем урок
		msgTitle := tgbotapi.NewMessage(msg.Chat.ID, lesson.Title)
		msgTitle.ParseMode = "Markdown"
		bot.Send(msgTitle)

		msgContent := tgbotapi.NewMessage(msg.Chat.ID, lesson.Content)
		msgContent.ParseMode = "Markdown"
		bot.Send(msgContent)

		imageDir := "./assets/images/"
		imageOutputPath := filepath.Join(imageDir, lesson.Image)

		fileImage, err := os.Open(imageOutputPath)
		if err != nil {
			lg.Logger.Warnf("failed to open file: %w", err)

		}
		defer fileImage.Close()

		err = helpers.SendMedia(fileImage, bot, msg.Chat.ID, lesson.Caption)
		if err != nil {
			lg.Logger.Warnf("failed to send photo: %w", err)
		}

		image2OutputPath := filepath.Join(imageDir, lesson.Image2)

		fileImage2, err := os.Open(image2OutputPath)
		if err != nil {
			lg.Logger.Warnf("failed to open file2: %w", err)

		}
		defer fileImage2.Close()

		err = helpers.SendMedia(fileImage2, bot, msg.Chat.ID, lesson.Caption2)
		if err != nil {
			lg.Logger.Warnf("failed to send photo2: %w", err)
		}

		reply := tgbotapi.NewMessage(msg.Chat.ID, lesson.Link)
		// reply.ParseMode = "Markdown"
		reply.ReplyMarkup = FeedbackButtons()
		_, err = bot.Send(reply)
		if err != nil {
			lg.Logger.Warnf("failed to send link: %w", err)
		}

		// Помечаем, что урок отправлен
		lessonModel.MarkLessonSent(msg.Chat.ID, int(lesson.ID))
	}
}

func HandleCallback(bot *tgbotapi.BotAPI, cb *tgbotapi.CallbackQuery) {
	if cb.Data == "feedback_good" || cb.Data == "feedback_bad" {
		// Сохраняем отзыв
		db.Conn.Exec("INSERT INTO user_feedback (chat_id, feedback, created_at) VALUES (?, ?, ?)",
			cb.Message.Chat.ID, cb.Data, time.Now())

		// Отвечаем пользователю
		reply := tgbotapi.NewMessage(cb.Message.Chat.ID, "Спасибо за ваш отзыв!")
		bot.Send(reply)
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
