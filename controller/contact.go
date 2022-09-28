package controller

import (
	"html/template"
	"net/http"
)

func GetContactMe(w http.ResponseWriter, r *http.Request) {
	var view, err = template.ParseFiles("views/contact.html")	
	if err != nil {
		panic(err.Error())
	}
	view.Execute(w, nil)
}