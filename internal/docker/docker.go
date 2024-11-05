package docker

import (
	"context"
	"fmt"
	"strings"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/image"
	"github.com/docker/docker/api/types/volume"
	"github.com/docker/docker/client"
)

func IsContainerID(id string) bool {
	return len(id) == 12
}


func IsImageID(id string) bool {
	// Check if id is an image ID
	return len(id) == 12
}

func IsVolumeID(id string) bool {
	// Check if id is a volume ID (not a full check but can be enhanced)
	return len(id) > 0
}

func IsContainerAction(data string) bool {
	return strings.HasPrefix(data, "start_") || strings.HasPrefix(data, "stop_") || strings.HasPrefix(data, "restart_") || strings.HasPrefix(data, "remove_")
}

func ParseContainerAction(data string) (string, string) {
	parts := strings.Split(data, "_")
	if len(parts) == 2 {
		return parts[0], parts[1]
	}
	return "", ""
}
func GetContainerInfo(id string) string {
	apiClient, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		return fmt.Sprintf("Error: %v", err)
	}
	defer apiClient.Close()

	containers, err := apiClient.ContainerList(context.Background(), container.ListOptions{All: true})
	if err != nil {
		return fmt.Sprintf("Error: %v", err)
	}

	// Найдем контейнер по ID
	for _, ctr := range containers {
		if ctr.ID[:12] == id {
			name := ctr.Names[0]
			if len(name) > 0 && name[0] == '/' {
				name = name[1:]
			}
			return fmt.Sprintf("Container: %s\nImage: %s\nStatus: %s\n", name, ctr.Image, ctr.Status)
		}
	}

	return "Container not found"
}


func GetImageInfo(id string) string {
	apiClient, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		return fmt.Sprintf("Error: %v", err)
	}
	defer apiClient.Close()

	images, err := apiClient.ImageList(context.Background(), image.ListOptions{All: true})
	if err != nil {
		return fmt.Sprintf("Error: %v", err)
	}

	for _, img := range images {
		if img.ID[:12] == id {
			return fmt.Sprintf("Image: %s\nTags: %v\nContainers: %d\n", img.ID[:12], img.RepoTags, img.Containers)
		}
	}

	return "Image not found"
}

func GetVolumeInfo(id string) string {
	apiClient, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		return fmt.Sprintf("Error: %v", err)
	}
	defer apiClient.Close()

	volumes, err := apiClient.VolumeList(context.Background(), volume.ListOptions{})
	if err != nil {
		return fmt.Sprintf("Error: %v", err)
	}

	for _, vol := range volumes.Volumes {
		if vol.Name == id {
			return fmt.Sprintf("Volume: %s\nLabels: %v\n", vol.Name, vol.Labels)
		}
	}

	return "Volume not found"
}

func StartContainer(id string) string {
	apiClient, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		return fmt.Sprintf("Error: %v", err)
	}
	defer apiClient.Close()

	err = apiClient.ContainerStart(context.Background(), id, container.StartOptions{})
	if err != nil {
		return fmt.Sprintf("Error starting container: %v", err)
	}

	return "Container started successfully"
}

func StopContainer(id string) string {
	apiClient, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		return fmt.Sprintf("Error: %v", err)
	}
	defer apiClient.Close()

	err = apiClient.ContainerStop(context.Background(), id, container.StopOptions{})
	if err != nil {
		return fmt.Sprintf("Error stopping container: %v", err)
	}

	return "Container stopped successfully"
}

func RestartContainer(id string) string {
	apiClient, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		return fmt.Sprintf("Error: %v", err)
	}
	defer apiClient.Close()

	err = apiClient.ContainerRestart(context.Background(), id, container.StopOptions{})
	if err != nil {
		return fmt.Sprintf("Error restarting container: %v", err)
	}

	return "Container restarted successfully"
}

func RemoveContainer(id string) string {
	apiClient, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		return fmt.Sprintf("Error: %v", err)
	}
	defer apiClient.Close()

	err = apiClient.ContainerRemove(context.Background(), id, container.RemoveOptions{Force: true})
	if err != nil {
		return fmt.Sprintf("Error removing container: %v", err)
	}

	return "Container removed successfully"
}

func RemoveImage(id string) string {
	apiClient, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		return fmt.Sprintf("Error: %v", err)
	}
	defer apiClient.Close()

	// ImageRemove returns two values: removed images and error
	removedImages, err := apiClient.ImageRemove(context.Background(), id, image.RemoveOptions{Force: true})
	if err != nil {
		return fmt.Sprintf("Error removing image: %v", err)
	}

	return fmt.Sprintf("Image removed successfully: %v", removedImages)
}

func RemoveVolume(id string) string {
	apiClient, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		return fmt.Sprintf("Error: %v", err)
	}
	defer apiClient.Close()

	err = apiClient.VolumeRemove(context.Background(), id, true)
	if err != nil {
		return fmt.Sprintf("Error removing volume: %v", err)
	}

	return "Volume removed successfully"
}
