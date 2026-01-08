package main

import (
	"fmt"
	"io"
	"log"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

const maxAudioSize = 40 * 1024 * 1024 // 20 MB

func main() {
	botToken := os.Getenv("TELEGRAM_BOT_TOKEN")
	if botToken == "" {
		log.Fatal("TELEGRAM_BOT_TOKEN not set")
	}

	bot, err := tgbotapi.NewBotAPI(botToken)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Authorized as @%s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := bot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message == nil {
			continue
		}

		chatID := update.Message.Chat.ID
		text := strings.TrimSpace(update.Message.Text)

		if !isYouTubeURL(text) {
			bot.Send(tgbotapi.NewMessage(
				chatID,
				"❌ Пришли корректную ссылку на YouTube-видео",
			))
			continue
		}

		if err := checkVideoDuration(text); err != nil {
			bot.Send(tgbotapi.NewMessage(
				chatID,
				"❌ Видео слишком длинное (максимум 2 часа)",
			))
			continue
		}

		// одно сообщение, которое будем редактировать
		statusMsg, _ := bot.Send(
			tgbotapi.NewMessage(chatID, "⏳ Скачивание"),
		)

		filePath, err := downloadMP3(text)

		if err != nil {
			bot.Send(tgbotapi.NewEditMessageText(
				chatID,
				statusMsg.MessageID,
				"❌ Ошибка скачивания",
			))
			log.Println(err)
			continue
		}

		// финальный статус
		bot.Send(tgbotapi.NewEditMessageText(
			chatID,
			statusMsg.MessageID,
			"✅ Готово, отправляю файл...",
		))

		if err := sendMP3(bot, chatID, filePath); err != nil {
			log.Println(err)
			bot.Send(tgbotapi.NewMessage(
				chatID,
				"❌ Ошибка при отправке файла",
			))
		}

		os.Remove(filePath)
	}

}

func downloadMP3(url string) (string, error) {
	tmpDir := os.TempDir()

	cmd := exec.Command(
		"yt-dlp",
		"-x",
		"--audio-format", "mp3",
		"--audio-quality", "0",
		"-o", filepath.Join(tmpDir, "%(title)s.%(ext)s"),
		url,
	)

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err := cmd.Run()
	if err != nil {
		return "", err
	}

	// ищем последний mp3 в tmp
	files, err := filepath.Glob(filepath.Join(tmpDir, "*.mp3"))
	if err != nil || len(files) == 0 {
		return "", fmt.Errorf("mp3 not found")
	}

	// берём самый свежий
	latest := files[len(files)-1]
	return latest, nil
}

func isYouTubeURL(raw string) bool {
	u, err := url.Parse(raw)
	if err != nil {
		return false
	}

	host := strings.ToLower(u.Host)

	switch host {
	case "www.youtube.com", "youtube.com", "m.youtube.com":
		// обычное видео
		if u.Path == "/watch" {
			return u.Query().Get("v") != ""
		}

		// shorts
		if strings.HasPrefix(u.Path, "/shorts/") {
			id := strings.TrimPrefix(u.Path, "/shorts/")
			id = strings.Trim(id, "/")
			return id != ""
		}

		return false

	case "youtu.be":
		// короткая ссылка
		return strings.Trim(u.Path, "/") != ""

	default:
		return false
	}
}

func sendMP3(bot *tgbotapi.BotAPI, chatID int64, filePath string) error {
	info, err := os.Stat(filePath)
	if err != nil {
		return err
	}

	sizeMB := info.Size() / (1024 * 1024)
	log.Printf("MP3 size: %d MB", sizeMB)

	// создаём временный файл с коротким ASCII-именем
	tmpFile, err := os.CreateTemp("", "audio-*.mp3")
	if err != nil {
		return err
	}
	tmpPath := tmpFile.Name()
	tmpFile.Close()

	// копируем данные
	if err := copyFile(filePath, tmpPath); err != nil {
		return err
	}
	defer os.Remove(tmpPath)

	file := tgbotapi.FilePath(tmpPath)

	if info.Size() <= maxAudioSize {
		audio := tgbotapi.NewAudio(chatID, file)
		_, err := bot.Send(audio)
		return err
	}

	doc := tgbotapi.NewDocument(chatID, file)
	_, err = bot.Send(doc)
	return err
}

func copyFile(src, dst string) error {
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()

	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, in)
	return err
}
