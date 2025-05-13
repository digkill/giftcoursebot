package logger

import (
	"github.com/sirupsen/logrus"
	"os"
)

type Log struct {
	Logger *logrus.Logger
}

func NewLog(logger *logrus.Logger) *Log {
	return &Log{
		Logger: logger,
	}
}

func (l Log) Init() {
	// Создаем директорию, если её нет
	err := os.MkdirAll("logs", os.ModePerm)
	if err != nil {
		l.Logger.Fatalf("Не удалось создать папку logs: %v", err)
	}

	// Открываем (или создаём) файл для логов
	file, err := os.OpenFile("logs/bot.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		l.Logger.Fatalf("Не удалось открыть лог-файл: %v", err)
	}

	// Настраиваем вывод логов
	l.Logger.SetOutput(file)
	l.Logger.SetLevel(logrus.DebugLevel) // выводим всё, включая debug
	l.Logger.SetFormatter(&logrus.TextFormatter{
		FullTimestamp: true,
	})
}
