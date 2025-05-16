package handlers

import (
	"github.com/digkill/giftcoursebot/internal/components/db"
	"github.com/digkill/giftcoursebot/internal/components/logger"
	"github.com/digkill/giftcoursebot/internal/helpers"
	"github.com/digkill/giftcoursebot/internal/models"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
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
👉 Первое задание уже на подходе. А если что — я всегда рядом!
*Удачи тебе! Пусть этот курс принесёт пользу и удовольствие!* 🌟
`

		// Получаем первый урок (day_number = 0)
		lesson := lessonModel.GetLessonByDay(0)
		if lesson == nil {
			bot.Send(tgbotapi.NewMessage(msg.Chat.ID, "Занятие не найдено. Обратитесь в поддержку."))
			return
		}

		helpers.SendLesson(bot, lg, lesson, msg.Chat.ID, welcomeMessage)

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
