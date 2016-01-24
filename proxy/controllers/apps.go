package controllers

import (
	"../container"
	"../models"
	"../shared"
	"archive/tar"
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/fsouza/go-dockerclient"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"io/ioutil"
	"path/filepath"
	"regexp"
	"strings"
	"time"
)

type AppsAPI struct {
	Context *shared.SharedContext
}

type File struct {
	Name string `json:"name"`
	Body string `json:"body"`
}

type RedeployAppVariable struct {
	Name     string
	Value    string
	Filename string
}

func (api *AppsAPI) ListApps(c *gin.Context) {
	a := models.AccountFromContext(c)
	var apps []models.App
	api.Context.PersistentDB.Where("account_id = ?", a.Id).Find(&apps)
	c.JSON(200, apps)
}

func (api *AppsAPI) GetApp(c *gin.Context) {
	a := models.AccountFromContext(c)
	var app models.App

	if err := api.Context.PersistentDB.Where("account_id = ? AND id = ?", a.Id, c.Params.ByName("id")).First(&app).Error; err != nil {
		c.String(404, "")
		return
	}

	api.Context.PersistentDB.Where("app_id = ?", app.Id).Find(&app.Variables)

	c.JSON(200, app)
}

func (api *AppsAPI) validate(c *gin.Context) (*models.App, *shared.FormErrors) {
	var form models.App
	c.BindWith(&form, binding.JSON)

	fe := shared.NewFormErrors()

	if form.Name == "" {
		fe.Add("name", "This field is required.")
	}

	match, _ := regexp.MatchString("^([a-zA-Z][0-9a-zA-Z_-]*)$", form.Name)
	if !match {
		fe.Add("name", "Only alphanumeric characters, underscore and '-' are allowed.")
	}

	if form.Memory < 0 || form.Memory > 16000 {
		fe.Add("memory", "Valid memory is between 0 and 1600 MB.")
	}

	var images []models.Image
	api.Context.PersistentDB.Find(&images)
	found := false
	for _, image := range images {
		if image.Id == form.Image_id {
			found = true
			break
		}
	}
	if !found {
		fe.Add("image_id", "Image doesn't exist.")
	}

	if fe.HasErrors() {
		return nil, &fe
	}

	return &form, nil
}

func (api *AppsAPI) EditApp(c *gin.Context) {
	a := models.AccountFromContext(c)
	app := models.App{}

	id := c.Params.ByName("id")

	if id != "" {
		if err := api.Context.PersistentDB.Where("account_id = ? AND id = ?", a.Id, id).First(&app).Error; err != nil {
			c.String(404, "")
			return
		}
	} else {
		app.Added_at = time.Now()
		app.Added_by_id = c.MustGet("user").(models.User).Id
		app.Account_id = a.Id
	}

	form, fe := api.validate(c)

	if fe != nil {
		c.JSON(400, fe)
		return
	}

	app.Name = form.Name
	app.Memory = form.Memory
	app.Image_id = form.Image_id

	api.Context.PersistentDB.Save(&app)

	for _, variable := range form.Variables {
		if variable.Value == "" {
			api.Context.PersistentDB.Where("app_id=? AND name=?", app.Id, variable.Name).Delete(models.AppVariable{})
		} else {
			v := models.AppVariable{}
			if err := api.Context.PersistentDB.Where("app_id=? AND name=?", app.Id, variable.Name).First(&v).Error; err != nil {
				v.Name = variable.Name
				v.App_id = app.Id
			}

			v.Value = variable.Value
			api.Context.PersistentDB.Save(&v)
		}
	}
}

func (api *AppsAPI) copyImageFile(contents string, a models.Account, app models.App, task models.Task, variables []RedeployAppVariable) (string, error) {
	contents = strings.Replace(contents, "#user#", a.Name, -1)

	uid := a.Uid()
	if uid == nil {
		task.Log(fmt.Sprintf("cannot get uid for user: %s", a.Name), "error", api.Context)
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

func (api *AppsAPI) StartApp(c *gin.Context) {
	a := models.AccountFromContext(c)
	var app models.App

	if err := api.Context.PersistentDB.Where("account_id = ? AND id = ?", a.Id, c.Params.ByName("id")).First(&app).Error; err != nil {
		c.String(404, "")
		return
	}

	user := c.MustGet("user").(models.User).Id

	//create task
	vars, _ := json.Marshal(struct {
		App models.App `json:"app"`
	}{
		App: app,
	})

	task := models.NewTask("start-app", string(vars), user)
	api.Context.PersistentDB.Save(&task)
	task.NotifyUser(*api.Context, user)

	var success = false
	defer func() {
		task.Duration = time.Now().Sub(task.Added_at).Seconds()
		task.Finished = true
		task.Success = success
		api.Context.PersistentDB.Save(&task)
		task.NotifyUser(*api.Context, user)
	}()

	id := ""
	name := app.ContainerName(a.Name)

	containers, err := container.GetAllContainers(api.Context)
	for _, c := range containers {
		for _, n := range c.Names {
			if strings.Replace(n, "/", "", -1) == name {
				id = c.ID
			}
		}
	}

	if id == "" {
		task.Log(fmt.Sprintf("container doesn't exist: %s", name), "error", api.Context)
		return
	}

	err = container.StartContainer(a, api.Context, api.Context.DockerClient, &app, id)
	if err != nil {
		task.Log(fmt.Sprintf("error starting the container: %s", err), "error", api.Context)
		return
	}

	task.Log(fmt.Sprintf("container with the name %s started", name), "info", api.Context)
	success = true
}

func (api *AppsAPI) StopApp(c *gin.Context) {
	a := models.AccountFromContext(c)
	var app models.App

	if err := api.Context.PersistentDB.Where("account_id = ? AND id = ?", a.Id, c.Params.ByName("id")).First(&app).Error; err != nil {
		c.String(404, "")
		return
	}

	user := c.MustGet("user").(models.User).Id

	//create task
	vars, _ := json.Marshal(struct {
		App models.App `json:"app"`
	}{
		App: app,
	})

	task := models.NewTask("stop-app", string(vars), user)
	api.Context.PersistentDB.Save(&task)
	task.NotifyUser(*api.Context, user)

	var success = false
	defer func() {
		task.Duration = time.Now().Sub(task.Added_at).Seconds()
		task.Finished = true
		task.Success = success
		api.Context.PersistentDB.Save(&task)
		task.NotifyUser(*api.Context, user)
	}()

	id := ""
	name := app.ContainerName(a.Name)

	containers, _ := container.GetAllContainers(api.Context)
	for _, c := range containers {
		for _, n := range c.Names {
			if strings.Replace(n, "/", "", -1) == name {
				id = c.ID
			}
		}
	}

	if id == "" {
		task.Log(fmt.Sprintf("container doesn't exist: %s", name), "error", api.Context)
		return
	}

	if err := api.Context.DockerClient.StopContainer(id, 10); err != nil {
		task.Log(fmt.Sprintf("cannot stop container: %s", err), "error", api.Context)
		return
	}

	task.Log(fmt.Sprintf("container with the name %s stopped", name), "info", api.Context)
	success = true
}

func (api *AppsAPI) RedeployApp(c *gin.Context) {
	a := models.AccountFromContext(c)
	var app models.App

	if err := api.Context.PersistentDB.Where("account_id = ? AND id = ?", a.Id, c.Params.ByName("id")).First(&app).Error; err != nil {
		c.String(404, "")
		return
	}

	user := c.MustGet("user").(models.User).Id

	//create task
	vars, _ := json.Marshal(struct {
		App models.App `json:"app"`
	}{
		App: app,
	})

	task := models.NewTask("redeploy-app", string(vars), user)
	api.Context.PersistentDB.Save(&task)
	task.NotifyUser(*api.Context, user)

	var success = false
	defer func() {
		task.Duration = time.Now().Sub(task.Added_at).Seconds()
		task.Finished = true
		task.Success = success
		api.Context.PersistentDB.Save(&task)
		task.NotifyUser(*api.Context, user)
		c.JSON(200, task)
	}()

	//get image
	var image models.Image

	if err := api.Context.PersistentDB.Where("id = ?", app.Image_id).First(&image).Error; err != nil {
		task.Log(fmt.Sprintf("image doesn't exist: %s", app.Image_id), "error", api.Context)
		return
	}

	name := app.ContainerName(a.Name)
	id := ""

	//get existing container name
	containers, err := container.GetAllContainers(api.Context)
	for _, c := range containers {
		for _, n := range c.Names {
			if strings.Replace(n, "/", "", -1) == name {
				id = c.ID
			}
		}
	}

	if id != "" {
		task.Log(fmt.Sprintf("stopping container with the name %s", name), "info", api.Context)

		//stop container
		if err := api.Context.DockerClient.StopContainer(id, 10); err != nil {
			task.Log(fmt.Sprintf("cannot stop container: %s", err), "info", api.Context)
		}

		task.Log(fmt.Sprintf("removing container with the name %s", name), "info", api.Context)

		//remove container
		if err := api.Context.DockerClient.RemoveContainer(docker.RemoveContainerOptions{
			ID: id,
		}); err != nil {
			task.Log(fmt.Sprintf("cannot remove container: %s", err), "error", api.Context)
			return
		}
	} else {
		task.Log(fmt.Sprintf("couldn't find container with name %s", name), "info", api.Context)
	}

	//remove image
	if err := api.Context.DockerClient.RemoveImage(fmt.Sprintf("manager/%s", name)); err != nil {
		task.Log(fmt.Sprintf("cannot remove image: %s", err), "info", api.Context)
	}

	//create dockerfile files
	var files []File

	//get image files
	folder := fmt.Sprintf("../images/%s/", image.Name)
	imageFiles, err := ioutil.ReadDir(folder)
	if err != nil {
		task.Log(fmt.Sprintf("image folder doesnt exist: %s", image.Name), "error", api.Context)
		return
	}

	//image variables
	imageVariables := image.GetVariables(api.Context)
	appVariables := app.GetVariables(api.Context)
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
		task.Log(fmt.Sprintf("copying image file %s", f.Name()), "info", api.Context)

		data, err := ioutil.ReadFile(filepath.Join(folder, f.Name()))
		if err != nil {
			task.Log(fmt.Sprintf("cannot read image file: %s", f.Name), "error", api.Context)
			return
		}

		var contents = string(data)

		if f.Name() == "Dockerfile" || f.Name() == "start.sh" {
			contents, err = api.copyImageFile(contents, *a, app, task, variables)
			if err != nil {
				return
			}
		} else {
			for _, v := range variables {
				if v.Filename != "" && f.Name() == v.Filename {
					contents = fmt.Sprintf("%s\n%s\n", contents, v.Value)
					contents, err = api.copyImageFile(contents, *a, app, task, variables)
					if err != nil {
						return
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
		task.Log(fmt.Sprintf("converting to json failed: %s", err), "error", api.Context)
		return
	}
	task.Log(string(b), "info", api.Context)

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
			task.Log(fmt.Sprintf("cannot write tar header: %s", err), "error", api.Context)
			return
		}
		if _, err := tw.Write([]byte(file.Body)); err != nil {
			task.Log(fmt.Sprintf("cannot write to tar: %s", err), "error", api.Context)
			return
		}
	}
	if err := tw.Close(); err != nil {
		task.Log(fmt.Sprintf("error closing tar archive: %s", err), "error", api.Context)
		return
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
	err = api.Context.DockerClient.BuildImage(opts)
	if err != nil {
		l := models.TaskLog{
			TaskId:   task.Id,
			Added_at: time.Now(),
			Value:    fmt.Sprintf("error encountered while building the image: %s", err),
			Type:     "error",
		}

		api.Context.PersistentDB.Save(&l)
		return
	}

	err = container.ReadOutputFromBuildImage(api.Context, task, buf)
	if err != nil {
		return
	}

	//create container
	cont, err := api.Context.DockerClient.CreateContainer(docker.CreateContainerOptions{
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
		task.Log(fmt.Sprintf("error creating the container: %s", err), "error", api.Context)
		return
	}

	//start container
	err = container.StartContainer(a, api.Context, api.Context.DockerClient, &app, cont.ID)
	if err != nil {
		task.Log(fmt.Sprintf("error starting the container: %s", err), "error", api.Context)
		return
	}

	//TODO: create task that will reload nginx/apache config

	success = true
}
