package helpers

import (
	"encoding/base64"
	"errors"
	"fmt"
	"github.com/digkill/giftcoursebot/internal/components/logger"
	"github.com/digkill/giftcoursebot/internal/models"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"io"
	"log"
	"mime"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
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
		err := SafeSend(bot, msg)
		if err != nil {
			log.Println("Ошибка отправки анимации:", err)
		}
		return err
	case ".png", ".jpg", ".jpeg":
		// Отправляем как фото
		msg := tgbotapi.NewPhoto(channelId, tgbotapi.FileReader{Name: file.Name(), Reader: file})
		msg.Caption = caption
		err := SafeSend(bot, msg)
		if err != nil {
			log.Println("Ошибка отправки фото:", err)
		}
		return err
	default:
		// Отправляем как документ
		msg := tgbotapi.NewDocument(channelId, tgbotapi.FileReader{Name: file.Name(), Reader: file})
		msg.Caption = caption
		err := SafeSend(bot, msg)
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
			err = SafeSend(bot, msg)
			return err
		} else {
			// остальное как фото
			msg := tgbotapi.NewPhoto(channelId, tgbotapi.FileReader{Name: file.Name(), Reader: file})
			msg.Caption = caption
			err = SafeSend(bot, msg)
			return err
		}
	} else if strings.HasPrefix(mimeType, "video/") {
		// видео (mp4) — через анимацию (или можешь через SendVideo, если хочешь)
		msg := tgbotapi.NewAnimation(channelId, tgbotapi.FileReader{Name: file.Name(), Reader: file})
		msg.Caption = caption
		err = SafeSend(bot, msg)
		return err
	} else {
		// все остальное — как документ
		msg := tgbotapi.NewDocument(channelId, tgbotapi.FileReader{Name: file.Name(), Reader: file})
		msg.Caption = caption
		err = SafeSend(bot, msg)
		return err
	}
}

// Универсальный отправитель с автоматическим retry при 429
func SafeSend(bot *tgbotapi.BotAPI, msg tgbotapi.Chattable) error {
	for {
		_, err := bot.Send(msg)
		if err == nil {
			return nil
		}

		var apiErr *tgbotapi.Error
		if errors.As(err, &apiErr) {
			// Корректная проверка RetryAfter
			if apiErr.ResponseParameters.RetryAfter > 0 {
				log.Printf("Telegram API rate limit hit. Retrying after %d seconds...\n", apiErr.ResponseParameters.RetryAfter)
				time.Sleep(time.Duration(apiErr.ResponseParameters.RetryAfter) * time.Second)
				continue
			}

			log.Printf("Telegram API error: %v\n", apiErr.Message)
			return apiErr
		}

		log.Printf("Unknown error sending message: %v\n", err)
		return err
	}
}

type MediaMeta struct {
	Width    int
	Height   int
	Duration int
}

func CreateInputMedia(fileName string, data []byte, caption string, isFirst bool, meta *MediaMeta) interface{} {
	ext := strings.ToLower(filepath.Ext(fileName))

	// Создаем FileBytes для файла
	requestFile := tgbotapi.FileBytes{
		Name:  fileName,
		Bytes: data,
	}

	switch ext {
	case ".gif":
		media := tgbotapi.NewInputMediaAnimation(requestFile)
		if meta != nil {
			media.Width = meta.Width
			media.Height = meta.Height
			media.Duration = meta.Duration
		}
		if isFirst {
			media.Caption = caption
			media.ParseMode = "Markdown"
		}
		return media

	case ".jpg", ".jpeg", ".png", ".webp":
		media := tgbotapi.NewInputMediaPhoto(requestFile)
		if isFirst {
			media.Caption = caption
			media.ParseMode = "Markdown"
		}
		return media

	case ".mp4":
		media := tgbotapi.NewInputMediaVideo(requestFile)
		if meta != nil {
			media.Width = meta.Width
			media.Height = meta.Height
			media.Duration = meta.Duration
		}
		if isFirst {
			media.Caption = caption
			media.ParseMode = "Markdown"
		}
		return media

	default:
		media := tgbotapi.NewInputMediaDocument(requestFile)
		if isFirst {
			media.Caption = caption
			media.ParseMode = "Markdown"
		}
		return media
	}
}

func SendLesson(bot *tgbotapi.BotAPI, lg *logger.Log, lesson *models.Lesson, chatId int64, message string) {

	// Отправляем урок
	msgTitle := tgbotapi.NewMessage(chatId, message+"  \n\n"+lesson.Title+"  \n\n"+lesson.Content+"  \n\n"+lesson.Link+"  \n\n"+"*Самостоятельные задания:*")
	msgTitle.ParseMode = "Markdown"
	msgTitle.ReplyMarkup = FeedbackButtons()
	_, err := bot.Send(msgTitle)
	if err != nil {
		lg.Logger.Warnf("failed to send start content: %v", err)
	}

	imageDir := "./assets/images/"

	imagePath1 := filepath.Join(imageDir, lesson.Image)
	imagePath2 := filepath.Join(imageDir, lesson.Image2)

	data1, err := os.ReadFile(imagePath1)
	if err != nil {
		lg.Logger.Warnf("failed to read %s: %v", imagePath1, err)
		return
	}

	data2, err := os.ReadFile(imagePath2)
	if err != nil {
		lg.Logger.Warnf("failed to read %s: %v", imagePath2, err)
		return
	}

	// Общая подпись (будет только у первого элемента)
	combinedCaption := fmt.Sprintf("📸 %s\n\n🖼️ %s", lesson.Caption, lesson.Caption2)

	// Метаданные для GIF (по желанию)
	//meta := &MediaMeta{Width: 480, Height: 320, Duration: 10}

	// Формируем список media
	mediaGroup := []interface{}{
		CreateInputMedia(lesson.Image, data1, combinedCaption, true, nil),
		CreateInputMedia(lesson.Image2, data2, "", false, nil),
	}
	_, err = bot.SendMediaGroup(tgbotapi.NewMediaGroup(chatId, mediaGroup))
	if err != nil {
		log.Printf("Failed to send media group: %v", err)
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

func ConvertGifToMp4(inputPath, outputPath string) error {
	cmd := exec.Command("ffmpeg", "-y", "-i", inputPath, "-movflags", "+faststart", "-pix_fmt", "yuv420p", outputPath)
	out, err := cmd.CombinedOutput()
	if err != nil {
		log.Printf("ffmpeg error: %v\nOutput: %s", err, out)
		return err
	}
	return nil
}

func ConvertGifToMp4Folder() {
	inputDir := "./assets/images"

	err := filepath.Walk(inputDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !info.IsDir() && strings.HasSuffix(strings.ToLower(info.Name()), ".gif") {
			outputPath := strings.TrimSuffix(path, filepath.Ext(path)) + ".mp4"
			log.Printf("Конвертируем: %s → %s\n", path, outputPath)
			if err := ConvertGifToMp4(path, outputPath); err != nil {
				log.Printf("Ошибка при конвертации %s: %v", path, err)
			}
		}
		return nil
	})

	if err != nil {
		log.Fatalf("Ошибка обхода директории: %v", err)
	}

	log.Println("Конвертация всех GIF завершена.")
}
