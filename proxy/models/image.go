package models

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"path/filepath"
	"strings"

	"../shared"
)

type Image struct {
	Id          int    `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Type        string `json:"type"`

	Ports     []ImagePort     `json:"ports"`
	Variables []ImageVariable `json:"variables"`
}

func (c Image) GetVariables(context *shared.SharedContext) []ImageVariable {
	var variables []ImageVariable

	images := GetImages(context)
	var image Image
	for _, i := range images {
		if i.Id == c.Id {
			image = i
		}
	}
	if image.Id != 0 {
		variables = image.Variables
	}

	return variables
}

func ParseImages(context *shared.SharedContext) error {
	var images []Image

	folder := "./images/"

	files, _ := ioutil.ReadDir(folder)
	for _, f := range files {
		// ignore shell images
		if strings.HasSuffix(f.Name(), "-shell") {
			continue
		}

		config, err := ioutil.ReadFile(filepath.Join(folder, f.Name(), "config.json"))
		if err != nil {
			return fmt.Errorf("config.json for image %s doesn't exist", f.Name())
		}

		dec := json.NewDecoder(bytes.NewReader(config))
		var image Image
		dec.Decode(&image)

		images = append(images, image)
	}

	context.Images = images

	return nil
}

func GetImages(context *shared.SharedContext) []Image {
	images := context.Images.([]Image)

	return images
}
