package routes

import (
	"github.com/gorilla/mux"
	"github.com/kartikjoshi267/simpletodolist/contollers"
)

func TaskRouter() *mux.Router {
	router := mux.NewRouter().PathPrefix("/api/task").Subrouter()

	router.HandleFunc("/", contollers.CreateTask).Methods("POST")
	router.HandleFunc("/", contollers.GetTasks).Methods("GET")
	router.HandleFunc("/{id}", contollers.UpdateTask).Methods("PUT")
	router.HandleFunc("/{id}", contollers.DeleteTask).Methods("DELETE")

	return router
}
