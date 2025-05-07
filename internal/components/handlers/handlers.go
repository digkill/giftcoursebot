package handlers

import (
	"github.com/digkill/giftcoursebot/internal/components/db"
	"github.com/digkill/giftcoursebot/internal/models"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"time"
)

func HandleMessage(bot *tgbotapi.BotAPI, userModel *models.UserModel, lessonModel *models.LessonModel, msg *tgbotapi.Message) {
	switch msg.Text {
	case "/start":
		// Регистрируем пользователя
		userModel.DB.Exec("INSERT IGNORE INTO users (chat_id, start_date) VALUES (?, ?)", msg.Chat.ID, time.Now())

		welcomeMessage := `🎉 *Добро пожаловать в наш подарок-курс!* 🎁

Ты только что открыл дверь в маленькое путешествие знаний, вдохновения и радости! 💡

📚 *Каждый день* тебя ждёт новый урок, только полезное и с заботой о тебе.

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

		// Отправляем сообщение
		messageLesson := tgbotapi.NewMessage(msg.Chat.ID, "🎒 *Дружок, держи урок!* 📚")
		messageLesson.ParseMode = "Markdown"

		bot.Send(messageLesson)

		reply := tgbotapi.NewMessage(msg.Chat.ID, lesson.Content)
		reply.ReplyMarkup = FeedbackButtons()
		bot.Send(reply)

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
