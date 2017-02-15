package models

type ImageVariable struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Default     string `json:"default"`
	Filename    string `json:"filename"`
}
