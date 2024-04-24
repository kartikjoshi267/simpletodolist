package routes

import (
	"github.com/gorilla/mux"
	"github.com/kartikjoshi267/simpletodolist/backend/controllers"
)

func TaskRouter() *mux.Router {
	router := mux.NewRouter().PathPrefix("/api/task").Subrouter()

	router.HandleFunc("/", controllers.CreateTask).Methods("POST")
	router.HandleFunc("/", controllers.GetTasks).Methods("GET")
	router.HandleFunc("/{id}", controllers.UpdateTask).Methods("PUT")
	router.HandleFunc("/{id}", controllers.DeleteTask).Methods("DELETE")

	return router
}
