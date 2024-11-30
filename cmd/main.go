package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os/exec"
	"strings"
)

const botToken = "7340918018:AAEaowo-v0Dbrm9FTS3fzIOCRbAXg6E2b6c"
const telegramAPI = "https://api.telegram.org/bot" + botToken

type Update struct {
	UpdateID int      `json:"update_id"`
	Message  *Message `json:"message"`
}

type Message struct {
	MessageID int    `json:"message_id"`
	Chat      Chat   `json:"chat"`
	Text      string `json:"text"`
}

type Chat struct {
	ID int64 `json:"id"`
}

type SendMessageRequest struct {
	ChatID int64  `json:"chat_id"`
	Text   string `json:"text"`
}

func main() {
	lastUpdateID := 0

	for {
		updates, err := getUpdates(lastUpdateID)
		if err != nil {
			log.Printf("Ошибка получения обновлений: %v\n", err)
			continue
		}

		for _, update := range updates {
			if update.Message != nil {
				handleMessage(update.Message)
				lastUpdateID = update.UpdateID + 1
			}
		}
	}
}

func getUpdates(offset int) ([]Update, error) {
	url := fmt.Sprintf("%s/getUpdates?offset=%d", telegramAPI, offset)
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var result struct {
		Ok     bool     `json:"ok"`
		Result []Update `json:"result"`
	}

	err = json.Unmarshal(body, &result)
	if err != nil {
		return nil, err
	}

	return result.Result, nil
}

func handleMessage(message *Message) {
	command := strings.ToLower(message.Text)

	switch {
	case command == "/containers":
		listContainers(message.Chat.ID)
	case strings.HasPrefix(command, "/start "):
		containerID := strings.TrimSpace(strings.TrimPrefix(command, "/start "))
		startContainer(message.Chat.ID, containerID)
	case strings.HasPrefix(command, "/stop "):
		containerID := strings.TrimSpace(strings.TrimPrefix(command, "/stop "))
		stopContainer(message.Chat.ID, containerID)
	default:
		sendMessage(message.Chat.ID, "Неизвестная команда. Доступные команды:\n/containers - список контейнеров\n/start <ID> - запустить контейнер\n/stop <ID> - остановить контейнер")
	}
}

func listContainers(chatID int64) {
	output, err := exec.Command("docker", "ps", "-a", "--format", "{{.ID}} {{.Image}} {{.Status}}").Output()
	if err != nil {
		sendMessage(chatID, fmt.Sprintf("Ошибка получения списка контейнеров: %v", err))
		return
	}

	if len(output) == 0 {
		sendMessage(chatID, "Контейнеры не найдены.")
	} else {
		sendMessage(chatID, string(output))
	}
}

func startContainer(chatID int64, containerID string) {
	err := exec.Command("docker", "start", containerID).Run()
	if err != nil {
		sendMessage(chatID, fmt.Sprintf("Ошибка запуска контейнера %s: %v", containerID, err))
	} else {
		sendMessage(chatID, fmt.Sprintf("Контейнер %s запущен.", containerID))
	}
}

func stopContainer(chatID int64, containerID string) {
	err := exec.Command("docker", "stop", containerID).Run()
	if err != nil {
		sendMessage(chatID, fmt.Sprintf("Ошибка остановки контейнера %s: %v", containerID, err))
	} else {
		sendMessage(chatID, fmt.Sprintf("Контейнер %s остановлен.", containerID))
	}
}

func sendMessage(chatID int64, text string) {
	url := fmt.Sprintf("%s/sendMessage", telegramAPI)

	message := SendMessageRequest{
		ChatID: chatID,
		Text:   text,
	}

	data, err := json.Marshal(message)
	if err != nil {
		log.Printf("Ошибка сериализации сообщения: %v\n", err)
		return
	}

	resp, err := http.Post(url, "application/json", bytes.NewBuffer(data))
	if err != nil {
		log.Printf("Ошибка отправки сообщения: %v\n", err)
		return
	}
	defer resp.Body.Close()
}
