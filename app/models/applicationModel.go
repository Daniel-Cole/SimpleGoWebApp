package models

type ApplicationList map[string][]Application

type Application struct {
	Environment string `json:"environment"`
	Name string `json:"name"`
	Artifact string `json:"artifact"`
	Tomcat string `json:"tomcat"`
	LastUpdated string `json:"lastUpdated"`
}