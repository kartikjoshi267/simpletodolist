package routes

import (
	"github.com/gorilla/mux"
	"github.com/kartikjoshi267/simpletodolist/backend/controllers"
)

func UserRouter() *mux.Router {
	router := mux.NewRouter().PathPrefix("/api/user").Subrouter()
	router.HandleFunc("/signup", controllers.CreateUser).Methods("POST")
	router.HandleFunc("/login", controllers.LoginUser).Methods("POST")
	router.HandleFunc("/", controllers.GetUser).Methods("GET")

	return router
}
