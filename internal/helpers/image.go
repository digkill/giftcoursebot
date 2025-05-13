package helpers

import (
	"encoding/base64"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"io"
	"log"
	"mime"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

func EncodeImageToBase64(imageBytes []byte, fileMimeType string) (string, error) {

	// Кодируем в base64
	base64Str := base64.StdEncoding.EncodeToString(imageBytes)

	// Определяем MIME-тип по расширению
	mimeType := mime.TypeByExtension(fileMimeType)
	if mimeType == "" {
		mimeType = "application/octet-stream" // По умолчанию, если неизвестный тип
	}

	// Формируем data URL
	dataURL := fmt.Sprintf("data:%s;base64,%s", mimeType, base64Str)

	return dataURL, nil
}

func IsImageOrVideo(ext string) bool {
	switch strings.ToLower(ext) {
	case ".jpg", ".jpeg", ".png", ".gif", ".mp4", ".mov", ".avi":
		return true
	default:
		return false
	}
}

func SendPhoto(file *os.File, bot *tgbotapi.BotAPI, channelId int64, caption string) error {
	photoMsg := tgbotapi.NewPhoto(channelId, tgbotapi.FileReader{Name: file.Name(), Reader: file})
	photoMsg.Caption = caption

	_, err := bot.Send(photoMsg)
	return err
}

func SendMedia(file *os.File, bot *tgbotapi.BotAPI, channelId int64, caption string) error {
	ext := strings.ToLower(filepath.Ext(file.Name()))

	switch ext {
	case ".gif", ".mp4":
		// Отправляем как анимацию (GIF будет воспроизводиться)
		msg := tgbotapi.NewAnimation(channelId, tgbotapi.FileReader{Name: file.Name(), Reader: file})
		msg.Caption = caption
		_, err := bot.Send(msg)
		if err != nil {
			log.Println("Ошибка отправки анимации:", err)
		}
		return err
	case ".png", ".jpg", ".jpeg":
		// Отправляем как фото
		msg := tgbotapi.NewPhoto(channelId, tgbotapi.FileReader{Name: file.Name(), Reader: file})
		msg.Caption = caption
		_, err := bot.Send(msg)
		if err != nil {
			log.Println("Ошибка отправки фото:", err)
		}
		return err
	default:
		// Отправляем как документ
		msg := tgbotapi.NewDocument(channelId, tgbotapi.FileReader{Name: file.Name(), Reader: file})
		msg.Caption = caption
		_, err := bot.Send(msg)
		if err != nil {
			log.Println("Ошибка отправки документа:", err)
		}
		return err
	}
}

func SendMediaSmart(file *os.File, bot *tgbotapi.BotAPI, channelId int64, caption string) error {
	// Считываем первые 512 байт для определения mime type
	header := make([]byte, 512)
	_, err := file.Read(header)
	if err != nil && err != io.EOF {
		log.Println("Ошибка чтения файла для определения MIME:", err)
		return err
	}

	// Определяем MIME type
	mimeType := http.DetectContentType(header)

	// Возвращаем указатель чтения файла в начало (чтобы Telegram смог прочитать весь файл)
	_, err = file.Seek(0, 0)
	if err != nil {
		log.Println("Ошибка возврата указателя файла:", err)
		return err
	}

	// Определяем способ отправки
	if strings.HasPrefix(mimeType, "image/") {
		if strings.HasSuffix(mimeType, "gif") {
			// gif как анимацию
			msg := tgbotapi.NewAnimation(channelId, tgbotapi.FileReader{Name: file.Name(), Reader: file})
			msg.Caption = caption
			_, err = bot.Send(msg)
			return err
		} else {
			// остальное как фото
			msg := tgbotapi.NewPhoto(channelId, tgbotapi.FileReader{Name: file.Name(), Reader: file})
			msg.Caption = caption
			_, err = bot.Send(msg)
			return err
		}
	} else if strings.HasPrefix(mimeType, "video/") {
		// видео (mp4) — через анимацию (или можешь через SendVideo, если хочешь)
		msg := tgbotapi.NewAnimation(channelId, tgbotapi.FileReader{Name: file.Name(), Reader: file})
		msg.Caption = caption
		_, err = bot.Send(msg)
		return err
	} else {
		// все остальное — как документ
		msg := tgbotapi.NewDocument(channelId, tgbotapi.FileReader{Name: file.Name(), Reader: file})
		msg.Caption = caption
		_, err = bot.Send(msg)
		return err
	}
}
