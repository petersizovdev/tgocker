package docker

import (
	"context"
	"fmt"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/image"
	"github.com/docker/docker/api/types/volume"
	"github.com/docker/docker/client"
)

func IsContainerID(id string) bool {
	// Check if id is a container ID
	return len(id) == 12
}

func IsImageID(id string) bool {
	// Check if id is an image ID
	return len(id) == 12
}

func IsVolumeID(id string) bool {
	// Check if id is a volume ID
	return true
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
