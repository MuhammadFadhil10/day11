package controller

import (
	"context"
	"fmt"
	"html/template"
	"mvcweb/connection"
	"mvcweb/helper"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
)

type ProjectData struct {
	Id int
	Name,Description,Image,Duration, StringStartDate, StringEndDate string
	StartDate,EndDate time.Time
	Technologies[]string	
}

type SessionStruct struct {
	Email,Name string
	IsLogin bool
}


func GetHome(w http.ResponseWriter, r *http.Request) {
	data, err := connection.Conn.Query(context.Background(), "SELECT id,name,start_date,end_date,description,technologies,image FROM public.tb_projects ORDER BY posted_at DESC")

	if err != nil {
		fmt.Println(err)
		panic(err.Error())
	}

	var dataResult []ProjectData

	for data.Next() {
		var project = ProjectData{}

		var err = data.Scan(&project.Id,&project.Name, &project.StartDate, &project.EndDate, &project.Description, &project.Technologies, &project.Image);
		if err != nil {
			panic(err.Error())
		}
		
		project.Duration = helper.GetDuration(project.StartDate.Format("2006-01-02"), project.EndDate.Format("2006-01-02"))
		project.Description = helper.CutString(project.Description, 30)
		project.Name = helper.CutString(project.Name, 20)

		dataResult = append(dataResult, project)
	}

	var view, templErr = template.ParseFiles("views/index.html", "views/layout/navigation.html");
	if templErr != nil {
		fmt.Println(templErr.Error())
		return
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
		"blogData": dataResult,
		"sessionData": sessionData,
	}
	
	view.Execute(w, responseData)
}

func GetAddProject(w http.ResponseWriter, r *http.Request) {
	var view, err = template.ParseFiles("views/project.html", "views/layout/navigation.html")	
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

func PostAddProject(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	
	name := r.PostForm.Get("name")
	description := r.PostForm.Get("description")
	startDate := r.PostForm.Get("start-date")
	endDate := r.PostForm.Get("end-date")
	techlist := r.PostForm["checkbox"]

	query := `INSERT INTO public.tb_projects(
		name, start_date, end_date, description, technologies, image)
		VALUES ($1,$2,$3,$4,$5,$6)`

	_ , err := connection.Conn.Exec(context.Background(),query, name,startDate,endDate,description,techlist,"https://images.unsplash.com/photo-1498050108023-c5249f4df085?ixlib=rb-1.2.1&ixid=MnwxMjA3fDB8MHxzZWFyY2h8Mnx8Y29kaW5nfGVufDB8fDB8fA%3D%3D&auto=format&fit=crop&w=500&q=60")

	if err != nil {
		panic(err.Error())
	}
	http.Redirect(w,r,"/form-add-project", http.StatusFound)
}

func GetEditProject(w http.ResponseWriter, r *http.Request) {
	projectId,idErr := strconv.Atoi(mux.Vars(r)["index"])
	if idErr != nil {
		panic(idErr.Error())
	}

	queryString := `
	SELECT id,name,start_date,end_date,description,technologies,image FROM public.tb_projects WHERE id = ($1)`

	data, err := connection.Conn.Query(context.Background(), queryString, projectId);

	if err != nil {
		panic(err.Error())
	}

	var project = ProjectData{};
	
	for data.Next() {
		err := data.Scan(&project.Id,&project.Name,&project.StartDate,&project.EndDate,&project.Description, &project.Technologies,&project.Image);
		if err != nil {
			fmt.Println(err.Error())
		}
		project.StringStartDate = project.StartDate.Format("2006-01-02")
		project.StringEndDate = project.EndDate.Format("2006-01-02")
	} 
	
	var view, viewErr = template.ParseFiles("views/edit-project.html", "views/layout/navigation.html")

	if viewErr != nil {
		panic(viewErr.Error())
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
		"project": project,
		"sessionData": sessionData,
	}
	view.Execute(w, responseData)
}

func UpdateProject(w http.ResponseWriter, r *http.Request) {
	parseErr := r.ParseForm();
	if parseErr != nil {
		panic(parseErr.Error())
	}
	projectId,idErr := strconv.Atoi(mux.Vars(r)["index"]);
	if idErr != nil {
		panic(idErr.Error())
	}

	name := r.PostForm.Get("name")
	startDate := r.PostForm.Get("start-date")
	endDate := r.PostForm.Get("end-date")
	description := r.PostForm.Get("description")
	checkbox := r.PostForm["checkbox"]
	
	queryString2 := `
		UPDATE public.tb_projects
		SET name=$1, start_date=$2, end_date=$3, description=$4, technologies=array_cat(technologies, $5)
		WHERE id = ($6)
	`

	_, queryErr := connection.Conn.Exec(context.Background(),queryString2,name,startDate,endDate,description,checkbox ,projectId)
	
	if queryErr != nil {
		panic(queryErr.Error())
	}
	
	http.Redirect(w,r,"/",http.StatusFound)
}

func GetProjectDetail(w http.ResponseWriter, r *http.Request) {
	projectId, indexError := strconv.Atoi(mux.Vars(r)["projectId"]);
	if indexError != nil {
		panic(indexError.Error())
	}

	queryString := `
		SELECT id,name,start_date,end_date,description,technologies,image 
		FROM public.tb_projects WHERE id = $1
	`
	data, err := connection.Conn.Query(context.Background(),queryString, projectId )

	if err != nil {
		fmt.Println(err)
		panic(err.Error())
	}

	var project = ProjectData{}
	for data.Next() {
		// var project = ProjectData{}

		var scanErr = data.Scan(&project.Id,&project.Name, &project.StartDate, &project.EndDate, &project.Description, &project.Technologies, &project.Image)
		if scanErr != nil {
			fmt.Println(scanErr.Error())
			return
		}
		project.Duration = helper.GetDuration(project.StartDate.Format("2006-01-02"), project.EndDate.Format("2006-01-02"))
		project.StringStartDate = project.StartDate.Format("January 02, 2006")
		project.StringEndDate = project.EndDate.Format("January 02, 2006")
	}

	var view,viewErr = template.ParseFiles("views/projectDetail.html", "views/layout/navigation.html")
	if viewErr != nil {
		panic(viewErr.Error())
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
		"project": project,
		"sessionData": sessionData,
	}
	view.Execute(w, responseData)

}

func DeleteProject(w http.ResponseWriter, r *http.Request) {
	projectId, idErr := strconv.Atoi(mux.Vars(r)["projectId"]);

	if idErr != nil {
		panic(idErr.Error())
	}

	queryString := `
		DELETE FROM public.tb_projects WHERE id = $1
	`

	_, queryErr := connection.Conn.Exec(context.Background(), queryString, projectId)

	if queryErr != nil {
		panic(queryErr.Error())
	}

	http.Redirect(w,r,"/",http.StatusFound)
}