package keyboard

import (
	"github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/image"
	"github.com/docker/docker/api/types/volume"
	"github.com/docker/docker/client"
	"context"
	"log"
)

func CreateMainMenuKeyboard() *tgbotapi.InlineKeyboardMarkup {
    return &tgbotapi.InlineKeyboardMarkup{
        InlineKeyboard: [][]tgbotapi.InlineKeyboardButton{
            {
                tgbotapi.NewInlineKeyboardButtonData("Containers", "containers"),
                tgbotapi.NewInlineKeyboardButtonData("Images", "images"),
                tgbotapi.NewInlineKeyboardButtonData("Volumes", "volumes"),
            },
        },
    }
}

func CreateContainersKeyboard() *tgbotapi.InlineKeyboardMarkup {
    apiClient, err := client.NewClientWithOpts(client.FromEnv)
    if err != nil {
        log.Panic(err)
    }
    defer apiClient.Close()

    containers, err := apiClient.ContainerList(context.Background(), container.ListOptions{All: true})
    if err != nil {
        log.Panic(err)
    }

    var rows [][]tgbotapi.InlineKeyboardButton
    for _, ctr := range containers {
        name := ctr.Names[0]
        if len(name) > 0 && name[0] == '/' {
            name = name[1:]
        }
        rows = append(rows, tgbotapi.NewInlineKeyboardRow(
            tgbotapi.NewInlineKeyboardButtonData(name, ctr.ID[:12]),
        ))
    }

    rows = append(rows, tgbotapi.NewInlineKeyboardRow(
        tgbotapi.NewInlineKeyboardButtonData("Back", "main_menu"),
    ))

    return &tgbotapi.InlineKeyboardMarkup{InlineKeyboard: rows}
}

func CreateImagesKeyboard() *tgbotapi.InlineKeyboardMarkup {
    apiClient, err := client.NewClientWithOpts(client.FromEnv)
    if err != nil {
        log.Panic(err)
    }
    defer apiClient.Close()

    images, err := apiClient.ImageList(context.Background(), image.ListOptions{All: true})
    if err != nil {
        log.Panic(err)
    }

    var rows [][]tgbotapi.InlineKeyboardButton
    for _, img := range images {
        rows = append(rows, tgbotapi.NewInlineKeyboardRow(
            tgbotapi.NewInlineKeyboardButtonData(img.RepoTags[0], img.ID[:12]),
        ))
    }

    rows = append(rows, tgbotapi.NewInlineKeyboardRow(
        tgbotapi.NewInlineKeyboardButtonData("Back", "main_menu"),
    ))

    return &tgbotapi.InlineKeyboardMarkup{InlineKeyboard: rows}
}

func CreateVolumesKeyboard() *tgbotapi.InlineKeyboardMarkup {
    apiClient, err := client.NewClientWithOpts(client.FromEnv)
    if err != nil {
        log.Panic(err)
    }
    defer apiClient.Close()

    volumes, err := apiClient.VolumeList(context.Background(), volume.ListOptions{})
    if err != nil {
        log.Panic(err)
    }

    var rows [][]tgbotapi.InlineKeyboardButton
    for _, vol := range volumes.Volumes {
        rows = append(rows, tgbotapi.NewInlineKeyboardRow(
            tgbotapi.NewInlineKeyboardButtonData(vol.Name, vol.Name),
        ))
    }

    rows = append(rows, tgbotapi.NewInlineKeyboardRow(
        tgbotapi.NewInlineKeyboardButtonData("Back", "main_menu"),
    ))

    return &tgbotapi.InlineKeyboardMarkup{InlineKeyboard: rows}
}

func CreateContainerActionsKeyboard(containerID string) *tgbotapi.InlineKeyboardMarkup {
    return &tgbotapi.InlineKeyboardMarkup{
        InlineKeyboard: [][]tgbotapi.InlineKeyboardButton{
            {
                tgbotapi.NewInlineKeyboardButtonData("Start", "start_"+containerID),
                tgbotapi.NewInlineKeyboardButtonData("Stop", "stop_"+containerID),
            },
            {
                tgbotapi.NewInlineKeyboardButtonData("Restart", "restart_"+containerID),
                tgbotapi.NewInlineKeyboardButtonData("Remove", "remove_"+containerID),
            },
            {
                tgbotapi.NewInlineKeyboardButtonData("Back", "containers"),
            },
        },
    }
}

func CreateBackKeyboard(parent string) *tgbotapi.InlineKeyboardMarkup {
    return &tgbotapi.InlineKeyboardMarkup{
        InlineKeyboard: [][]tgbotapi.InlineKeyboardButton{
            {
                tgbotapi.NewInlineKeyboardButtonData("Back", parent),
            },
        },
    }
}