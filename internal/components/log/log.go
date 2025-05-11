package logger

import (
	"github.com/sirupsen/logrus"
	"os"
)

type Log struct {
	logger *logrus.Logger
}

func NewLog(logger *logrus.Logger) *Log {
	return &Log{
		logger: logger,
	}
}

func (l Log) Init() {
	// Создаем директорию, если её нет
	err := os.MkdirAll("logs", os.ModePerm)
	if err != nil {
		l.logger.Fatalf("не удалось создать папку logs: %v", err)
	}

	// Открываем (или создаём) файл для логов
	file, err := os.OpenFile("logs/bot.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		l.logger.Fatalf("не удалось открыть лог-файл: %v", err)
	}

	// Настраиваем вывод логов
	l.logger.SetOutput(file)
	l.logger.SetLevel(logrus.DebugLevel) // выводим всё, включая debug
	l.logger.SetFormatter(&logrus.TextFormatter{
		FullTimestamp: true,
	})
}
