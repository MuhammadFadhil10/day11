package controller

import (
	"context"
	"fmt"
	"html/template"
	"mvcweb/connection"
	"net/http"

	"golang.org/x/crypto/bcrypt"
)

type User struct {
	Name, password string
}

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

// post 

func Register(w http.ResponseWriter, r *http.Request) {
	r.ParseForm();

	name := r.PostForm.Get("name") 
	email := r.PostForm.Get("email") 
	password := r.PostForm.Get("password");

	hashedPassword, hashedErr := bcrypt.GenerateFromPassword([]byte(password), 12);
	if hashedErr != nil {
		fmt.Println(hashedErr)
	} 

	queryString := `
		INSERT INTO public.tb_user(name,email,password) VALUES ($1,$2,$3)
	`

	_, dbErr := connection.Conn.Exec(context.Background(), queryString, name, email, hashedPassword);

	if dbErr != nil {
		fmt.Println(dbErr)
	}

	http.Redirect(w,r,"/form-login", 301)

}

func Login(w http.ResponseWriter, r *http.Request) {
	r.ParseForm();

	email := r.PostForm.Get("email") 
	password := r.PostForm.Get("password");

	queryString := `
		SELECT email,password FROM public.tb_user WHERE email = $1
	`

	data, dataErr := connection.Conn.Query(context.Background(),queryString,email)

	if dataErr != nil {
		fmt.Println(dataErr.Error())
	}

	var user = User{}
	for data.Next() {
		scanErr := data.Scan(&user.Name, &user.password);
		if scanErr != nil {
			fmt.Println(scanErr)
		}
	}

	passwordMatch := bcrypt.CompareHashAndPassword([]byte(user.password), []byte(password))

	if passwordMatch != nil {
		http.Redirect(w,r,"/form-login", 301)
	} else {
		http.Redirect(w,r,"/",301)
	}

}