package utils

import "github.com/daniel-cole/SimpleGoWebApp/app/models"

func GenEnvKey(environment models.Environment) string {
	return environment.Environment
}