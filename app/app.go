package main

import (
	"io/ioutil"
	"os"
	"github.com/daniel-cole/SimpleGoWebApp/app/database"
	"github.com/daniel-cole/SimpleGoWebApp/app/log"
	"github.com/daniel-cole/SimpleGoWebApp/app/routes"
	"github.com/daniel-cole/SimpleGoWebApp/app/utils"
	"net/http"
	"strconv"
)


func main() {
	log.Init(ioutil.Discard, os.Stdout, os.Stdout, os.Stderr)

	environment := os.Getenv("ENV")
	log.LogInfo.Printf("SimpleGoWebApp starting up in environment: %s", environment)

	//read in application specific configuration from config.json
	dbName, port, dbTimeout := utils.ParseConfig()

	//set up context for application handlers
	context := utils.Context{
		DBConn:      database.InitDB(dbName),
		DBTimeout:   dbTimeout,
		DBBucketApp: []byte("apps"),
		DBBucketEnv: []byte("envs"),
	}

	//parse tomcat/springboot applications to add baseline versions
	utils.ParseTomcats(context.DBConn, context.DBBucketApp)
	utils.ParseSpringBoots(context.DBConn, context.DBBucketApp)

	//initalise routes
	routes.Init(context)

	//start up app server
	log.LogInfo.Printf("SCAVM starting up on port %d\n", port)
	log.LogFatal("", http.ListenAndServe(":" + strconv.Itoa(port), nil))

}
