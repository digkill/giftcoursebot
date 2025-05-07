package main

import (
	"github.com/digkill/giftcoursebot/internal/components/handlers"
	"log"
	"os"

	"github.com/digkill/giftcoursebot/internal/components/db"
	"github.com/digkill/giftcoursebot/internal/components/scheduler"

	"github.com/digkill/giftcoursebot/internal/models"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func main() {
	dsn := os.Getenv("MYSQL_DSN")
	botToken := os.Getenv("TELEGRAM_BOT_TOKEN")

	if dsn == "" || botToken == "" {
		log.Fatal("MYSQL_DSN or TELEGRAM_BOT_TOKEN is not set")
	}

	db.Init(dsn)

	bot, err := tgbotapi.NewBotAPI(botToken)
	if err != nil {
		log.Panic(err)
	}

	userModel := &models.UserModel{DB: db.Conn}
	lessonModel := &models.LessonModel{DB: db.Conn}

	go scheduler.StartScheduler(bot, userModel, lessonModel)

	log.Println("Bot is running...")

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60
	updates := bot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message != nil {
			handlers.HandleMessage(bot, userModel, lessonModel, update.Message)
		} else if update.CallbackQuery != nil {
			handlers.HandleCallback(bot, update.CallbackQuery)
		}
	}
}
