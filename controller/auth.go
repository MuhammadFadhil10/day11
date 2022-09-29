package controller

import (
	"context"
	"fmt"
	"html/template"
	"mvcweb/connection"
	"net/http"

	"github.com/gorilla/sessions"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	Name, Email, password string
}

func GetRegisterForm(w http.ResponseWriter, r *http.Request) {
	view, viewErr := template.ParseFiles("views/register.html", "views/layout/navigation.html")

	if viewErr != nil {
		fmt.Println(viewErr.Error())
	}

	view.Execute(w,nil);

}
func GetLoginForm(w http.ResponseWriter, r *http.Request) {
	view, viewErr := template.ParseFiles("views/login.html", "views/layout/navigation.html")

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
		SELECT email,name,password FROM public.tb_user WHERE email = $1
	`

	data, dataErr := connection.Conn.Query(context.Background(),queryString,email)

	if dataErr != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("message:" + dataErr.Error()))
		fmt.Println(dataErr.Error())
	}

	var user = User{}
	for data.Next() {
		scanErr := data.Scan(&user.Email, &user.Name,&user.password);
		if scanErr != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("message:" + scanErr.Error()))
			fmt.Println(scanErr)
		}
	}

	passwordMatch := bcrypt.CompareHashAndPassword([]byte(user.password), []byte(password))

	if passwordMatch != nil {
		w.WriteHeader(http.StatusBadRequest)
		http.Redirect(w,r,"/form-login", 301)
		w.Write([]byte("Wrong password!, " + passwordMatch.Error()));
		return;
	}

	store := sessions.NewCookieStore([]byte("MY_SESSION_KEY"))
	session, _ := store.Get(r, "MY_SESSION_KEY")

	session.Values["Email"] = user.Email;
	session.Values["Name"] = user.Name;
	session.Values["IsLogin"] = true;
	session.Options.MaxAge = 7200;

	session.AddFlash("Login succesfully!", "login message")

	session.Save(r,w)
	

	http.Redirect(w,r,"/", http.StatusMovedPermanently)
}

func Logout(w http.ResponseWriter, r *http.Request) {
	store := sessions.NewCookieStore([]byte("MY_SESSION_KEY"))
	session, _ := store.Get(r, "MY_SESSION_KEY")
	session.Options.MaxAge = -1;
	err := session.Save(r,w)
	if err != nil {
		fmt.Println(err.Error())
	}
	http.Redirect(w,r,"/", http.StatusTemporaryRedirect)
}