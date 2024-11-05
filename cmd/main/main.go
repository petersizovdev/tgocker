package main

import (
	"log"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/petersizovdev/tgocker/internal/docker"
	"github.com/petersizovdev/tgocker/internal/keyboard"
)

func main() {
	bot, err := tgbotapi.NewBotAPI("7340918018:AAEaowo-v0Dbrm9FTS3fzIOCRbAXg6E2b6c")
	if err != nil {
		log.Panic(err)
	}

	bot.Debug = true

	log.Printf("Authorized on account %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := bot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message != nil && update.Message.IsCommand() {
			handleCommand(bot, update.Message)
		} else if update.CallbackQuery != nil {
			handleCallback(bot, update.CallbackQuery)
		}
	}
}

func handleCommand(bot *tgbotapi.BotAPI, message *tgbotapi.Message) {
	msg := tgbotapi.NewMessage(message.Chat.ID, "")

	switch message.Command() {
	case "start":
		msg.Text = "Hello! I am a Docker management bot. Please select a command:"
		msg.ReplyMarkup = keyboard.CreateMainMenuKeyboard()
	default:
		msg.Text = "I don't know that command"
	}

	if _, err := bot.Send(msg); err != nil {
		log.Panic(err)
	}
}

func handleCallback(bot *tgbotapi.BotAPI, callback *tgbotapi.CallbackQuery) {
	callbackData := callback.Data
	chatID := callback.Message.Chat.ID
	messageID := callback.Message.MessageID

	// Создайте новый объект для редактирования сообщения
	editMsg := tgbotapi.NewEditMessageText(chatID, messageID, "")

	// Обработка команд меню и выбора контейнеров, образов и томов
	switch callbackData {
	case "main_menu":
		editMsg.Text = "Please select a command:"
		editMsg.ReplyMarkup = keyboard.CreateMainMenuKeyboard()
	case "containers":
		editMsg.Text = "Select a container:"
		editMsg.ReplyMarkup = keyboard.CreateContainersKeyboard()
	case "images":
		editMsg.Text = "Select an image:"
		editMsg.ReplyMarkup = keyboard.CreateImagesKeyboard()
	case "volumes":
		editMsg.Text = "Select a volume:"
		editMsg.ReplyMarkup = keyboard.CreateVolumesKeyboard()
	default:
		// Проверяем, является ли это действием с контейнером
		if docker.IsContainerAction(callbackData) {
			// Извлекаем действие и ID контейнера
			action, containerID := docker.ParseContainerAction(callbackData)
			handleContainerAction(bot, &editMsg, action, containerID)
		} else if docker.IsContainerID(callbackData) {
			// Если это просто выбор контейнера, показываем его информацию
			editMsg.Text = docker.GetContainerInfo(callbackData)
			editMsg.ReplyMarkup = keyboard.CreateContainerActionsKeyboard(callbackData)
		} else if docker.IsImageID(callbackData) {
			// Для образа: показываем информацию об образе
			editMsg.Text = docker.GetImageInfo(callbackData)
			editMsg.ReplyMarkup = keyboard.CreateBackKeyboard("images")
		} else if docker.IsVolumeID(callbackData) {
			// Для тома: показываем информацию о томе
			editMsg.Text = docker.GetVolumeInfo(callbackData)
			editMsg.ReplyMarkup = keyboard.CreateBackKeyboard("volumes")
		} else {
			editMsg.Text = "Unknown command"
		}
	}

	// Проверяем, отличается ли новый текст или кнопки от старого
	if editMsg.Text != callback.Message.Text || !compareInlineKeyboard(editMsg.ReplyMarkup, callback.Message.ReplyMarkup) {
		// Отправляем только если есть изменения
		if _, err := bot.Send(editMsg); err != nil {
			log.Panic(err)
		}
	}
}

// Сравнение InlineKeyboardMarkup
func compareInlineKeyboard(newMarkup, oldMarkup *tgbotapi.InlineKeyboardMarkup) bool {
	if newMarkup == nil && oldMarkup == nil {
		return true
	}
	if newMarkup == nil || oldMarkup == nil {
		return false
	}
	if len(newMarkup.InlineKeyboard) != len(oldMarkup.InlineKeyboard) {
		return false
	}
	for i := range newMarkup.InlineKeyboard {
		if len(newMarkup.InlineKeyboard[i]) != len(oldMarkup.InlineKeyboard[i]) {
			return false
		}
		for j := range newMarkup.InlineKeyboard[i] {
			if newMarkup.InlineKeyboard[i][j].Text != oldMarkup.InlineKeyboard[i][j].Text ||
				newMarkup.InlineKeyboard[i][j].CallbackData != oldMarkup.InlineKeyboard[i][j].CallbackData {
				return false
			}
		}
	}
	return true
}



func handleContainerAction(bot *tgbotapi.BotAPI, editMsg *tgbotapi.EditMessageTextConfig, action string, containerID string) {
	var result string

	// Выполнение действия с контейнером
	switch action {
	case "start":
		result = docker.StartContainer(containerID)
	case "stop":
		result = docker.StopContainer(containerID)
	case "restart":
		result = docker.RestartContainer(containerID)
	case "remove":
		result = docker.RemoveContainer(containerID)
	default:
		result = "Unknown action"
	}

	// Обновляем текст сообщения с результатом действия
	editMsg.Text = result

	// Обновляем кнопки для дальнейших действий
	editMsg.ReplyMarkup = keyboard.CreateContainerActionsKeyboard(containerID)

	// Отправляем обновленное сообщение
	if _, err := bot.Send(editMsg); err != nil {
		log.Panic(err)
	}
}
