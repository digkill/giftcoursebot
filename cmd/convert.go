package main

import (
	"fmt"
	"github.com/digkill/giftcoursebot/internal/helpers"
	"os"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("❌ Нет команды. Используй: convert")
		return
	}

	command := os.Args[1]

	switch command {
	case "convert":
		fmt.Println("🚀 Запускаем конвертацию GIF → MP4...")
		helpers.ConvertGifToMp4Folder()
	default:
		fmt.Println("❌ Неизвестная команда:", command)
	}
}
