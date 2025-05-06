package main

import (
	"github.com/digkill/giftcoursebot/internal/components/db"
	"github.com/digkill/giftcoursebot/internal/components/handlers"
	"github.com/digkill/giftcoursebot/internal/components/scheduler"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
	"log"
	"os"
)

func main() {
	logrus.SetLevel(logrus.DebugLevel)

	if err := godotenv.Load(); err != nil {
		logrus.Warnf("load env failed: %v", err)
	}

	var token = os.Getenv("YOUR_TELEGRAM_BOT_TOKEN")

	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		log.Panic(err)
	}

	db := db.InitDB()
	defer db.Close()

	go scheduler.StartScheduler(bot, db)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := bot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message != nil {
			handlers.HandleMessage(bot, db, update.Message)
		} else if update.CallbackQuery != nil {
			handlers.HandleCallback(bot, db, update.CallbackQuery)
		}
	}
}
