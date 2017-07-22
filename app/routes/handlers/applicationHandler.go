package handlers

import (
	"github.com/daniel-cole/SimpleGoWebApp/app/database"
	"net/http"
	"io/ioutil"
	"github.com/daniel-cole/SimpleGoWebApp/app/utils"
	"github.com/daniel-cole/SimpleGoWebApp/app/models"
	"github.com/daniel-cole/SimpleGoWebApp/app/log"
	"encoding/json"
	"github.com/gorilla/mux"
	"strings"
	"github.com/boltdb/bolt"
)

//An AppUpdateHandler to update entries in the database
func AppUpdateHandler(context utils.Context, w http.ResponseWriter, r *http.Request) {

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.LogInfo.Printf("Error reading body: %v", err)
		http.Error(w, "unable to read request body", http.StatusBadRequest)
		return
	}

	log.LogInfo.Printf("Attempting to update : %s\n", string(body))

	//use decode here instead of unmarshal
	application, err := utils.ParseJson(body)

	if err != nil {
		log.LogError.Println(err)
		http.Error(w, "unable to parse JSON", http.StatusBadRequest)
		return
	}

	dbConn := context.DBConn
	dbBucket := context.DBBucketApp

	key := []byte(application.Environment + "_" + application.Name)
	value := []byte(body)

	if err := database.InsertDBValue(dbConn, dbBucket, key, value); err != nil {
		log.LogInfo.Printf("Failed to update DB: %v", err)
	}

	log.LogInfo.Printf("application environment: %s\n", application.Environment)

}

//A AppReadHandler to read entries from the database
func AppReadHandler(context utils.Context, w http.ResponseWriter, r *http.Request) {
	log.LogInfo.Println("accessing route:", r.URL.Path)

	vars := mux.Vars(r)

	env := vars["environment"]
	app := vars["application"]

	dbConn := context.DBConn
	dbBucket := context.DBBucketApp

	//return everything
	if env == "" && app == "" {
		getAllApplications(dbConn, dbBucket, w)
		return
	}

	if env != "" && app == "" {
		getAllAppsInEnv(env, dbConn, dbBucket, w)
		return
	}

	getAppInEnv(env, app, dbConn, dbBucket, w)

}

func getAllApplications(dbConn *bolt.DB, dbBucket []byte, w http.ResponseWriter) {
	values, err := database.ReadAllDBValues(dbConn, dbBucket);
	if err != nil {
		log.LogInfo.Printf("Failed to read db values: %v", err)
	}

	//initialise application list
	applicationList := models.ApplicationList{}

	for _, value := range values {
		var application models.Application
		err := json.Unmarshal([]byte(value), &application)
		if err != nil {
			log.LogInfo.Println("failed to retrieve values from database: ", err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		//add each application found
		applicationList[application.Environment] = append(applicationList[application.Environment], application)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(&applicationList)
}

func getAllAppsInEnv(env string, dbConn *bolt.DB, dbBucket []byte, w http.ResponseWriter) {
	//return all apps in specified environment

	values, err := database.ReadAllDBValues(dbConn, dbBucket);
	if err != nil {
		log.LogInfo.Printf("Failed to read db value: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var applications []string
	for key, value := range values {
		application := strings.Split(key, "_")[0]
		log.LogInfo.Println(application)
		if application == env {
			applications = append(applications, value)
		}
	}

	if applications == nil {
		log.LogInfo.Printf("no applications found for route")
		http.Error(w, "no applications found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(applications)

}

func getAppInEnv(env string, app string, dbConn *bolt.DB, dbBucket []byte, w http.ResponseWriter) {
	key := []byte(env + "_" + app)

	val, err := database.ReadDBValue(dbConn, dbBucket, key);
	if err != nil {
		log.LogInfo.Printf("Failed to read db value: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if val == nil {
		http.Error(w, "app not found", http.StatusNotFound)
		return
	}

	log.LogInfo.Println(val)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(string(val))
}

//A AppDeleteHandler to delete an entry from the database
func AppDeleteHandler(context utils.Context, w http.ResponseWriter, r *http.Request) {

	dbConn := context.DBConn
	dbBucket := context.DBBucketApp

	vars := mux.Vars(r)

	env := vars["environment"]
	app := vars["application"]

	key := []byte(env + "_" + app)

	if err := database.DeleteDBValue(dbConn, dbBucket, key); err != nil {
		log.LogInfo.Printf("Failed to read db value: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("200 - OK: valued deleted or was not found"))

}
