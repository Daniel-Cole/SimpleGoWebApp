package utils

import (
	"github.com/go-yaml/yaml"
	"strings"
	"io/ioutil"
	"os"
	"github.com/boltdb/bolt"
	"encoding/json"
	"github.com/daniel-cole/SimpleGoWebApp/app/database"
	"github.com/daniel-cole/SimpleGoWebApp/app/log"
	"github.com/daniel-cole/SimpleGoWebApp/app/models"
	"time"
	"bufio"
	"fmt"
)

type AppYAML struct {
	Tomcats map[string]interface{}
}

type AppConfig struct {
	DBName string
	Port float64
	DBTimeout float64
}

func ParseConfig() (string, int, int){
	log.LogInfo.Printf("Attempting to parse application configuration defined in config.json")

	pwd ,_ := os.Getwd()
	config, err := ioutil.ReadFile(pwd + string(os.PathSeparator) + "config.json")
	if err != nil {
		fmt.Println("Unable to find config file - please specify configuration manually")

		reader := bufio.NewReader(os.Stdin)

		fmt.Print("Specify the database name for the application: ")
		dbName, _ := reader.ReadString('\n')
		fmt.Print("Specify the port for the application: ")

		var port int
		_, err := fmt.Scanf("%d", &port)
		if(err != nil){
			log.LogFatal("Unable to parse port - ", err)
		}

		defaultDBTimeout := 10
		return strings.TrimSpace(dbName), port, defaultDBTimeout
	}

	var appConfig AppConfig
	if err := json.Unmarshal(config, &appConfig); err != nil {
		log.LogFatal("Unable to parse config JSON. Please check that the configuration is well formed JSON", err)
	}

	dbName := appConfig.DBName
	dbTimeout := int(appConfig.DBTimeout)
	port := int(appConfig.Port)

	log.LogInfo.Println("Successfully parsed application config")
	return dbName, port, dbTimeout
}

func ParseTomcats(dbConn *bolt.DB, dbBucket []byte){
	pwd ,_ := os.Getwd()

	applicationsYaml, err := ioutil.ReadFile(pwd + string(os.PathSeparator) + "applications-tomcat.yaml")
	if err != nil {
		log.LogFatal("Unable to read applications yaml file.", err)
	}

	var appYAML AppYAML
	yaml.Unmarshal(applicationsYaml, &appYAML)

	tcDeployedPackages := []models.Application{}

	//refactor this out into search method for YAML/JSON
	for tomcat, attributes := range appYAML.Tomcats {
		for attribute, attributeValue := range attributes.(map[interface{}]interface{}) {
			if attribute == "applications" {
				for app, pkg := range attributeValue.(map[interface{}]interface{}) {
					for identifier, pkgName := range pkg.(map[interface{}]interface{}){
						if identifier == "package" {
							//trim the package name so that only the artifact remains
							pkgSplit := strings.Split(pkgName.(string),"/")
							artifact := pkgSplit[len(pkgSplit)-1]
							tcDeployedPackages = append(tcDeployedPackages, models.Application{
								Environment: "baseline",
								Name:        app.(string),
								Artifact:     artifact,
								Tomcat:      tomcat,
							})
						}
					}
				}
			}

		}
	}

	//update database with tomcats after parsing yaml file
	lastUpdated := time.Time.Format(time.Now(),"2006-01-02 15:04:05")

	for _, application := range tcDeployedPackages {
		application.LastUpdated = lastUpdated
		applicationJSON, err := json.Marshal(application)
		if err != nil {
			log.LogError.Println("Failed to parse Tomcat applications YAML file")
			return
		}
		database.InsertDBValue(dbConn, dbBucket, []byte(GenAppKey(application)), applicationJSON)
	}

	log.LogInfo.Println("Successfully parsed Tomcat applications YAML file")

}

func ParseSpringBoots(dbConn *bolt.DB, dbBucket []byte){

}