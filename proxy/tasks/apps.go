package tasks

import (
	"../container"
	"../helpers"
	"../models"
	"../shared"
	"archive/tar"
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/fsouza/go-dockerclient"
	"io/ioutil"
	"path/filepath"
	"strings"
	"time"
)

type File struct {
	Name string `json:"name"`
	Body string `json:"body"`
}

type RedeployAppVariable struct {
	Name     string
	Value    string
	Filename string
}

func copyImageFile(contents string, a models.Account, app *models.App, task models.Task, variables []RedeployAppVariable, context *shared.SharedContext) (string, error) {
	contents = strings.Replace(contents, "#user#", a.Name, -1)

	uid := a.Uid()
	if uid == nil {
		task.Log(fmt.Sprintf("cannot get uid for user: %s", a.Name), "error", context)
		return "", errors.New("cannot get uid for user")
	}
	contents = strings.Replace(contents, "#uid#", *uid, -1)

	contents = strings.Replace(contents, "#appname#", app.Name, -1)

	for _, v := range variables {
		name := fmt.Sprintf("#variable_%s#", v.Name)
		contents = strings.Replace(contents, name, v.Value, -1)
	}

	return contents, nil
}

func StartApp(app *models.App, a *models.Account, user int, context *shared.SharedContext) models.Task {
	//create task
	vars, _ := json.Marshal(struct {
		App models.App `json:"app"`
	}{
		App: *app,
	})

	task := models.NewTask("start-app", string(vars), user)
	context.PersistentDB.Save(&task)
	task.NotifyUser(*context, user)

	var success = false
	defer func() {
		task.Duration = time.Now().Sub(task.Added_at).Seconds()
		task.Finished = true
		task.Success = success
		context.PersistentDB.Save(&task)
		task.NotifyUser(*context, user)
	}()

	id := ""
	name := app.ContainerName(a.Name)

	containers, err := container.GetAllContainers(context)
	for _, c := range containers {
		for _, n := range c.Names {
			if strings.Replace(n, "/", "", -1) == name {
				id = c.ID
			}
		}
	}

	if id == "" {
		task.Log(fmt.Sprintf("container doesn't exist: %s", name), "error", context)
		return task
	}

	err = container.StartContainer(a, context, context.DockerClient, app, id)
	if err != nil {
		task.Log(fmt.Sprintf("error starting the container: %s", err), "error", context)
		return task
	}

	task.Log(fmt.Sprintf("container with the name %s started", name), "info", context)
	success = true

	return task
}

func StopApp(app *models.App, a *models.Account, user int, context *shared.SharedContext) models.Task {
	//create task
	vars, _ := json.Marshal(struct {
		App models.App `json:"app"`
	}{
		App: *app,
	})

	task := models.NewTask("stop-app", string(vars), user)
	context.PersistentDB.Save(&task)
	task.NotifyUser(*context, user)

	var success = false
	defer func() {
		task.Duration = time.Now().Sub(task.Added_at).Seconds()
		task.Finished = true
		task.Success = success
		context.PersistentDB.Save(&task)
		task.NotifyUser(*context, user)
	}()

	id := ""
	name := app.ContainerName(a.Name)

	containers, _ := container.GetAllContainers(context)
	for _, c := range containers {
		for _, n := range c.Names {
			if strings.Replace(n, "/", "", -1) == name {
				id = c.ID
			}
		}
	}

	if id == "" {
		task.Log(fmt.Sprintf("container doesn't exist: %s", name), "error", context)
		return task
	}

	if err := context.DockerClient.StopContainer(id, 10); err != nil {
		task.Log(fmt.Sprintf("cannot stop container: %s", err), "error", context)
		return task
	}

	task.Log(fmt.Sprintf("container with the name %s stopped", name), "info", context)
	success = true

	return task
}

func RedeployApp(app *models.App, a *models.Account, user int, context *shared.SharedContext) models.Task {
	//create task
	vars, _ := json.Marshal(struct {
		App models.App `json:"app"`
	}{
		App: *app,
	})

	task := models.NewTask("redeploy-app", string(vars), user)
	context.PersistentDB.Save(&task)
	task.NotifyUser(*context, user)

	var success = false
	defer func() {
		task.Duration = time.Now().Sub(task.Added_at).Seconds()
		task.Finished = true
		task.Success = success
		context.PersistentDB.Save(&task)
		task.NotifyUser(*context, user)
	}()

	//get image
	var image models.Image

	if err := context.PersistentDB.Where("id = ?", app.Image_id).First(&image).Error; err != nil {
		task.Log(fmt.Sprintf("image doesn't exist: %s", app.Image_id), "error", context)
		return task
	}

	name := app.ContainerName(a.Name)
	id := ""

	//get existing container name
	containers, err := container.GetAllContainers(context)
	for _, c := range containers {
		for _, n := range c.Names {
			if strings.Replace(n, "/", "", -1) == name {
				id = c.ID
			}
		}
	}

	if id != "" {
		task.Log(fmt.Sprintf("stopping container with the name %s", name), "info", context)

		//stop container
		if err := context.DockerClient.StopContainer(id, 10); err != nil {
			task.Log(fmt.Sprintf("cannot stop container: %s", err), "info", context)
		}

		task.Log(fmt.Sprintf("removing container with the name %s", name), "info", context)

		//remove container
		if err := context.DockerClient.RemoveContainer(docker.RemoveContainerOptions{
			ID: id,
		}); err != nil {
			task.Log(fmt.Sprintf("cannot remove container: %s", err), "error", context)
			return task
		}
	} else {
		task.Log(fmt.Sprintf("couldn't find container with name %s", name), "info", context)
	}

	//remove image
	if err := context.DockerClient.RemoveImage(fmt.Sprintf("manager/%s", name)); err != nil {
		task.Log(fmt.Sprintf("cannot remove image: %s", err), "info", context)
	}

	//create dockerfile files
	var files []File

	//get image files
	folder := fmt.Sprintf("../images/%s/", image.Name)
	imageFiles, err := ioutil.ReadDir(folder)
	if err != nil {
		task.Log(fmt.Sprintf("image folder doesnt exist: %s", image.Name), "error", context)
		return task
	}

	//image variables
	imageVariables := image.GetVariables(context)
	appVariables := app.GetVariables(context)
	var variables []RedeployAppVariable

	for _, i := range imageVariables {
		found := false
		for _, a := range appVariables {
			if a.Name == i.Name {
				variables = append(variables, RedeployAppVariable{
					Name:     i.Name,
					Value:    a.Value,
					Filename: i.Filename,
				})
				found = true
				break
			}
		}

		if !found {
			//use default value
			variables = append(variables, RedeployAppVariable{
				Name:     i.Name,
				Value:    i.Default,
				Filename: i.Filename,
			})
		}
	}

	for _, f := range imageFiles {
		task.Log(fmt.Sprintf("copying image file %s", f.Name()), "info", context)

		data, err := ioutil.ReadFile(filepath.Join(folder, f.Name()))
		if err != nil {
			task.Log(fmt.Sprintf("cannot read image file: %s", f.Name), "error", context)
			return task
		}

		var contents = string(data)

		if f.Name() == "Dockerfile" || f.Name() == "start.sh" {
			contents, err = copyImageFile(contents, *a, app, task, variables, context)
			if err != nil {
				return task
			}
		} else {
			for _, v := range variables {
				if v.Filename != "" && f.Name() == v.Filename {
					contents = fmt.Sprintf("%s\n%s\n", contents, v.Value)
					contents, err = copyImageFile(contents, *a, app, task, variables, context)
					if err != nil {
						return task
					}
				}
			}
		}

		files = append(files, File{
			Name: f.Name(),
			Body: contents,
		})
	}

	b, err := json.Marshal(files)
	if err != nil {
		task.Log(fmt.Sprintf("converting to json failed: %s", err), "error", context)
		return task
	}
	task.Log(string(b), "info", context)

	//build image
	inputbuf := bytes.NewBuffer(nil)

	tw := tar.NewWriter(inputbuf)
	for _, file := range files {
		hdr := &tar.Header{
			Name: file.Name,
			Mode: 0644,
			Size: int64(len(file.Body)),
		}
		if err := tw.WriteHeader(hdr); err != nil {
			task.Log(fmt.Sprintf("cannot write tar header: %s", err), "error", context)
			return task
		}
		if _, err := tw.Write([]byte(file.Body)); err != nil {
			task.Log(fmt.Sprintf("cannot write to tar: %s", err), "error", context)
			return task
		}
	}
	if err := tw.Close(); err != nil {
		task.Log(fmt.Sprintf("error closing tar archive: %s", err), "error", context)
		return task
	}

	buf := bytes.NewBuffer(nil)
	opts := docker.BuildImageOptions{
		Name:                fmt.Sprintf("manager/%s", name),
		NoCache:             true,
		RmTmpContainer:      true,
		ForceRmTmpContainer: true,
		OutputStream:        buf,
		RawJSONStream:       true,
		SuppressOutput:      false,
		InputStream:         inputbuf,
	}

	//call build image API
	err = context.DockerClient.BuildImage(opts)
	if err != nil {
		l := models.TaskLog{
			TaskId:   task.Id,
			Added_at: time.Now(),
			Value:    fmt.Sprintf("error encountered while building the image: %s", err),
			Type:     "error",
		}

		context.PersistentDB.Save(&l)
		return task
	}

	err = container.ReadOutputFromBuildImage(context, task, buf)
	if err != nil {
		return task
	}

	//create container
	cont, err := context.DockerClient.CreateContainer(docker.CreateContainerOptions{
		Config: &docker.Config{
			User:   a.Name,
			Image:  fmt.Sprintf("manager/%s", name),
			Memory: int64(app.Memory * 1024 * 1024),
		},
		HostConfig: &docker.HostConfig{
			RestartPolicy: docker.RestartPolicy{
				Name:              "always",
				MaximumRetryCount: 100,
			},
		},
		Name: name,
	})

	if err != nil {
		task.Log(fmt.Sprintf("error creating the container: %s", err), "error", context)
		return task
	}

	//start container
	err = container.StartContainer(a, context, context.DockerClient, app, cont.ID)
	if err != nil {
		task.Log(fmt.Sprintf("error starting the container: %s", err), "error", context)
		return task
	}

	// reload web servers' config
	go helpers.SyncWebServersForAccount(a, user, context)

	success = true

	return task
}
