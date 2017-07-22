package models

type EnvironmentList map[string][]Application

type Environment struct {
	Environment string `json:"environment"`
	Owner string `json:"owner"`
	Purpose string `json:"purpose"`
	Description string `json:"description"`
}