package main

import (
	"github.com/digkill/giftcoursebot/internal/components/db"
	"github.com/digkill/giftcoursebot/internal/components/scheduler"
	"github.com/digkill/giftcoursebot/internal/models"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log"
	"os"
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

	log.Println("Bot is running... Press Ctrl+C to stop")
	select {} // block forever
}
