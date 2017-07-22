package routes

import (
	"net/http"
	"github.com/gorilla/mux"
	"github.com/daniel-cole/SimpleGoWebApp/app/routes/handlers"
	"github.com/daniel-cole/SimpleGoWebApp/app/utils"
	"github.com/daniel-cole/SimpleGoWebApp/app/log"
)

type ContextHandler struct {
	Context utils.Context
	ContextFunc ContextFunc
}

type ContextFunc func(context utils.Context, w http.ResponseWriter, r *http.Request)

func (handler ContextHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	handler.ContextFunc(handler.Context, w, r)
}

func Init(context utils.Context){
	log.LogInfo.Println("Setting up routes for SCAVM")
	r := mux.NewRouter()

	//************************
	//ENVIRONMENT ROUTES
	//************************

	//GET
	r.Handle("/environments", ContextHandler{context, handlers.EnvReadHandler}).Methods("GET")
	r.Handle("/environments/{owner}", ContextHandler{context, handlers.EnvReadHandler}).Methods("GET")
	r.Handle("/environments/{owner}/{environment}", ContextHandler{context, handlers.EnvReadHandler}).Methods("GET")

	//POST
	r.Handle("/environments/{owner}/{environment}", ContextHandler{context, handlers.EnvReadHandler}).Methods("POST")

	//DELETE
	r.Handle("/environments/{owner}/{environment}", ContextHandler{context, handlers.EnvReadHandler}).Methods("DELETE")

	//************************
	//APPLICATION ROUTES
	//************************
	//GET
	r.Handle("/applications", ContextHandler{context, handlers.AppReadHandler}).Methods("GET")
	r.Handle("/applications/{environment}", ContextHandler{context, handlers.AppReadHandler}).Methods("GET")
	r.Handle("/applications/{environment}/{application}", ContextHandler{context, handlers.AppReadHandler}).Methods("GET")

	//POST
	r.Handle("/applications", ContextHandler{context, handlers.AppUpdateHandler}).Methods("POST")

	//DELETE
	r.Handle("/applications/{environment}/{application}", ContextHandler{context, handlers.AppDeleteHandler}).Methods("DELETE")



	//handler to server dashboard
	r.HandleFunc("/", handlers.DashboardHandler)

	//ensure static directory is served on '/'
	r.PathPrefix("/").Handler(http.FileServer(http.Dir("./static/")))

	//required for gorilla mux routes to work
	http.Handle("/", r)
}
