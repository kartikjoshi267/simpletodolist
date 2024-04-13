package routes

import (
	"github.com/gorilla/mux"
	"github.com/kartikjoshi267/simpletodolist/contollers"
)

func UserRouter() *mux.Router {
	router := mux.NewRouter().PathPrefix("/api/user").Subrouter()
	router.HandleFunc("/signup", contollers.CreateUser).Methods("POST")
	router.HandleFunc("/login", contollers.LoginUser).Methods("POST")
	router.HandleFunc("/", contollers.GetUser).Methods("GET")

	return router
}
