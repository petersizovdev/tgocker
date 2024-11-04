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

	editMsg := tgbotapi.NewEditMessageText(chatID, messageID, "")

	switch callbackData {
	case "main_menu":
		editMsg.Text = "Please select a command:"
		mainMenuKeyboard := keyboard.CreateMainMenuKeyboard()
		editMsg.ReplyMarkup = &mainMenuKeyboard
	case "containers":
		editMsg.Text = "Select a container:"
		containersKeyboard := keyboard.CreateContainersKeyboard()
		editMsg.ReplyMarkup = &containersKeyboard
	case "images":
		editMsg.Text = "Select an image:"
		imagesKeyboard := keyboard.CreateImagesKeyboard()
		editMsg.ReplyMarkup = &imagesKeyboard
	case "volumes":
		editMsg.Text = "Select a volume:"
		volumesKeyboard := keyboard.CreateVolumesKeyboard()
		editMsg.ReplyMarkup = &volumesKeyboard
	default:
		if docker.IsContainerID(callbackData) {
			editMsg.Text = docker.GetContainerInfo(callbackData)
			backKeyboard := keyboard.CreateBackKeyboard("containers")
			editMsg.ReplyMarkup = &backKeyboard
		} else if docker.IsImageID(callbackData) {
			editMsg.Text = docker.GetImageInfo(callbackData)
			backKeyboard := keyboard.CreateBackKeyboard("images")
			editMsg.ReplyMarkup = &backKeyboard
		} else if docker.IsVolumeID(callbackData) {
			editMsg.Text = docker.GetVolumeInfo(callbackData)
			backKeyboard := keyboard.CreateBackKeyboard("volumes")
			editMsg.ReplyMarkup = &backKeyboard
		} else {
			editMsg.Text = "Unknown command"
		}
	}

	if _, err := bot.Send(editMsg); err != nil {
		log.Panic(err)
	}
}
