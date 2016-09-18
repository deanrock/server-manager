package container

import (
	"../models"
	"../shared"
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"github.com/fsouza/go-dockerclient"
	"strings"
	"time"
)

type Container struct {
}

func GetHostConfig(account *models.Account, context *shared.SharedContext, dockerClient *docker.Client, app *models.App) (*docker.HostConfig, error) {
	var links []string

	if context != nil {
		apps := account.Apps()

		var images []models.Image
		context.PersistentDB.Find(&images)

		linkOtherApps := false
		if app != nil {
			var image models.Image
			if err := context.PersistentDB.Where("id = ?", app.Image_id).First(&image).Error; err != nil {
				return nil, errors.New("image doesnt exist")
			}

			if image.Type != "database" {
				linkOtherApps = true
			}
		} else {
			linkOtherApps = true
		}

		if linkOtherApps {
			for _, app := range apps {
				for _, img := range images {
					if img.Id == app.Image_id && img.Type == "database" {
						name := fmt.Sprintf("app-%s-%s:%s", account.Name, app.Name, app.Name)
						links = append(links, name)
					}
				}
			}
		}
	}

	return &docker.HostConfig{
		Binds:      []string{"/home/" + account.Name + ":/home/" + account.Name},
		ExtraHosts: []string{"mysql:172.17.42.1", "postgres:172.17.42.1"},
		Links:      links,
	}, nil
}

func StartContainer(dockerClient *docker.Client, containerId string) error {
	err := dockerClient.StartContainer(containerId, nil)

	return err
}

func GetAllContainers(context *shared.SharedContext) ([]docker.APIContainers, error) {
	containers, err := context.DockerClient.ListContainers(docker.ListContainersOptions{
		All: true,
	})

	return containers, err
}

func ReadOutputFromBuildImage(context *shared.SharedContext, task models.Task, buf *bytes.Buffer) error {
	//read output from building the image

	tx := context.PersistentDB.Begin()

	var line = ""
	scanner := bufio.NewScanner(buf)
	for scanner.Scan() {
		line = scanner.Text()

		l := models.TaskLog{
			TaskId:   task.Id,
			Added_at: time.Now(),
			Value:    line,
			Type:     "log",
		}

		tx.Save(&l)
	}

	tx.Commit()

	if err := scanner.Err(); err != nil {
		l := models.TaskLog{
			TaskId:   task.Id,
			Added_at: time.Now(),
			Value:    fmt.Sprintf("error encountered while reading output: %s", err),
			Type:     "error",
		}

		context.PersistentDB.Save(&l)
		return errors.New(fmt.Sprintf("error encountered while reading output: %s", err))
	}

	if !strings.Contains(line, "Successfully built") {
		l := models.TaskLog{
			TaskId:   task.Id,
			Added_at: time.Now(),
			Value:    fmt.Sprintf("last line doesn't contain 'Successfully built'"),
			Type:     "error",
		}

		context.PersistentDB.Save(&l)
		return errors.New(fmt.Sprintf("last line doesn't contain 'Successfully built'"))
	}

	return nil
}
