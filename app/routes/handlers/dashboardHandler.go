package handlers

import (
	"net/http"
	"text/template"
)

func DashboardHandler(w http.ResponseWriter, r *http.Request){

	t, _ := template.ParseFiles("templates/base.tmpl", "templates/menu.tmpl", "templates/application.tmpl", "templates/environment.tmpl")
	t.ExecuteTemplate(w, "base", "")

}