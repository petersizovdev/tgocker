package main

import (
	"context"
	"fmt"
	"log"

	"github.com/docker/docker/client"
)

func main() {
	// Инициализация Docker клиента
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		log.Fatalf("Не удалось инициализировать Docker клиент: %v", err)
	}

	// Получение списка контейнеров
	containers, err := cli.ContainerList(context.Background(), types.ContainerListOptions{})
	if err != nil {
		log.Fatalf("Не удалось получить список контейнеров: %v", err)
	}

	// Вывод списка контейнеров в лог
	for _, container := range containers {
		fmt.Printf("ID: %s, Image: %s, Status: %s\n", container.ID[:10], container.Image, container.Status)
	}
}