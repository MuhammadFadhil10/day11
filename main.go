package main

import (
	"fmt"
	"mvcweb/connection"
	"mvcweb/controller"
	"net/http"

	"github.com/gorilla/mux"
)




func main() {
	router := mux.NewRouter()

	directory := http.Dir("./public")
	fileServer := http.FileServer(directory)

    router.PathPrefix("/static/").Handler(http.StripPrefix("/static", fileServer))

	// router
	// get
	router.HandleFunc("/", controller.GetHome).Methods("GET")
	router.HandleFunc("/form-add-project", controller.GetAddProject).Methods("GET")
	router.HandleFunc("/form-edit-project/{index}", controller.GetEditProject).Methods("GET")
	router.HandleFunc("/contact-me", controller.GetContactMe).Methods("GET")
	router.HandleFunc("/project/{projectId}", controller.GetProjectDetail).Methods("GET")
	// post
	router.HandleFunc("/add-project", controller.PostAddProject).Methods("POST")
	router.HandleFunc("/update-project/{index}", controller.UpdateProject).Methods("POST")
	router.HandleFunc("/delete-project/{projectId}", controller.DeleteProject).Methods("POST")
	
	connection.DatabaseConnect(func() {
		fmt.Println("running on port 5000");
		http.ListenAndServe("localhost:5000", router);
	})

}











