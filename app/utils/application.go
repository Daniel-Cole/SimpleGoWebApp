package utils

import (
	"encoding/json"
	"github.com/daniel-cole/SimpleGoWebApp/app/models"
)

func GenAppKey(application models.Application) string {
	return application.Environment + "_" + application.Name
}

func ParseJson(body []byte) (models.Application, error) {
	var a models.Application
	if err := json.Unmarshal(body, &a); err != nil {
		return a, err
	}
	return a, nil
}
