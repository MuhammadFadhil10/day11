package controller

import (
	"html/template"
	"net/http"

	"github.com/gorilla/sessions"
)

func GetContactMe(w http.ResponseWriter, r *http.Request) {
	var view, err = template.ParseFiles("views/contact.html", "views/layout/navigation.html")	
	if err != nil {
		panic(err.Error())
	}
	store := sessions.NewCookieStore([]byte("MY_SESSION_KEY"))
	session, _ := store.Get(r,"MY_SESSION_KEY");
	if session.Values["IsLogin"] == nil {
		session.Values["IsLogin"] = false
	} 
	if session.Values["Name"] == nil {
		session.Values["Name"] = ""
	} 
	if session.Values["Email"] == nil{
		session.Values["Email"] = ""
	}
	sessionData := SessionStruct{
		Name: session.Values["Name"].(string),
		Email: session.Values["Email"].(string),
		IsLogin: session.Values["IsLogin"].(bool),
	}
	responseData := map[string]interface{} {
		"sessionData": sessionData,
	}
	view.Execute(w, responseData)
}