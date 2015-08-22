package container

import (
	"../models"
	"../shared"
	"fmt"
	"github.com/fsouza/go-dockerclient"
)

type Container struct {
}

func StartContainer(account *models.Account, context *shared.SharedContext, dockerClient *docker.Client, containerId string) error {
	var links []string

	if context != nil {
		apps := account.Apps()

		var images []models.Image
		context.PersistentDB.Find(&images)

		for _, app := range apps {
			for _, img := range images {
				if img.Id == app.Image_id && img.Type == "database" {
					name := fmt.Sprintf("app-%s-%s:%s", account.Name, app.Name, app.Name)
					links = append(links, name)
				}
			}
		}
	}

	err := dockerClient.StartContainer(containerId,
		&docker.HostConfig{
			Binds:      []string{"/home/" + account.Name + ":/home/" + account.Name},
			ExtraHosts: []string{"mysql:172.17.42.1"},
			Links:      links,
		})

	return err
}

func GetAllContainers(context *shared.SharedContext) ([]docker.APIContainers, error) {
	containers, err := context.DockerClient.ListContainers(docker.ListContainersOptions{
		All: true,
	})

	return containers, err
}
