package handlers

import (
	"github.com/digkill/giftcoursebot/internal/components/db"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func HandleMessage(bot *tgbotapi.BotAPI, db *db.DB, msg *tgbotapi.Message) {
	switch msg.Text {
	case "/start":
		db.RegisterUser(msg.Chat.ID)
		lesson := lessons[0]
		reply := tgbotapi.NewMessage(msg.Chat.ID, lesson)
		reply.ReplyMarkup = FeedbackButtons()
		bot.Send(reply)
	}
}

func HandleCallback(bot *tgbotapi.BotAPI, db *DB, cb *tgbotapi.CallbackQuery) {
	if cb.Data == "feedback_good" || cb.Data == "feedback_bad" {
		db.SaveFeedback(cb.Message.Chat.ID, cb.Data)
		reply := tgbotapi.NewMessage(cb.Message.Chat.ID, "Спасибо за ваш отзыв!")
		bot.Send(reply)
	}
}
