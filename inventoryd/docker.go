package main

import (
	"context"
	"fmt"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"github.com/go-co-op/gocron"
	"go.uber.org/zap"
	"reflect"
	"sort"
	"strings"
	"time"
)

type DockerService struct {
	cli             *client.Client
	presenceService *PresenceService
	runningImages   []string
}

func NewDockerService(ps *PresenceService) (*DockerService, error) {
	cli, err := client.NewClientWithOpts()
	if err != nil {
		return nil, err
	}

	svc := &DockerService{
		cli:             cli,
		presenceService: ps,
	}

	return svc, nil
}

func (ds *DockerService) Run(cronExpression string) error {
	s := gocron.NewScheduler(time.UTC)

	ds.updateIfNeeded()

	_, err := s.Cron(cronExpression).Do(ds.updateIfNeeded)
	if err != nil {
		return err
	}

	s.StartBlocking()

	return nil
}

func (ds *DockerService) updateIfNeeded() {
	// poll running containers for whether a presence update is required
	updateNeeded, err := ds.updateNeeded()
	if err != nil {
		zap.S().Error(err)
		return
	}

	// update presence message
	if updateNeeded {
		// build map of image repos/tags
		updatedImages := make(map[string]string)
		for _, image := range ds.runningImages {
			// extract image repo+tag
			repo, err := getImageRepository(image)
			if err != nil {
				zap.S().Error(err)
				continue
			}
			tag, err := getImageTag(image)
			if err != nil {
				zap.S().Error(err)
				continue
			}

			// add image to updatedImages
			updatedImages[repo] = tag
		}

		// update presence message
		err = ds.presenceService.Update(updatedImages)
		if err != nil {
			zap.S().Error(err)
		}

		return
	}
}

func (ds *DockerService) updateNeeded() (bool, error) {
	// initialise context
	ctx, cancelFunc := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancelFunc()

	// list running images
	var runningImages []string
	containers, err := ds.cli.ContainerList(ctx, types.ContainerListOptions{})
	if err != nil {
		return false, err
	}

	for _, container := range containers {
		if strings.ToLower(container.State) == "running" {
			runningImages = append(runningImages, container.Image)
		}
	}

	sort.Strings(runningImages)
	sort.Strings(ds.runningImages)

	if !reflect.DeepEqual(runningImages, ds.runningImages) {
		ds.runningImages = runningImages
		return true, nil
	}

	return false, nil
}

func getImageTag(image string) (string, error) {
	// split image into component tokens
	tokens := strings.Split(image, ":")
	nTokens := len(tokens)

	// return err if invalid image format
	if nTokens < 2 {
		return "", fmt.Errorf("invalid container name: %s", image)
	}

	return tokens[nTokens-1], nil
}

func getImageRepository(image string) (string, error) {
	// split image into component tokens
	tokens := strings.Split(image, ":")
	nTokens := len(tokens)

	// return err if invalid image format
	if nTokens < 2 {
		return "", fmt.Errorf("invalid container name: %s", image)
	}

	return tokens[nTokens-2], nil
}