package controller

import (
	"fmt"
	"html/template"
	"net/http"
)

func GetRegisterForm(w http.ResponseWriter, r *http.Request) {
	view, viewErr := template.ParseFiles("views/register.html", "views/layout/layout.html")

	if viewErr != nil {
		fmt.Println(viewErr.Error())
	}

	view.Execute(w,nil);

}
func GetLoginForm(w http.ResponseWriter, r *http.Request) {
	view, viewErr := template.ParseFiles("views/login.html", "views/layout/layout.html")

	if viewErr != nil {
		fmt.Println(viewErr.Error())
	}

	view.Execute(w,nil);

}