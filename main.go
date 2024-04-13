package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/kartikjoshi267/simpletodolist/routes"
)

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Welcome to Simple Todo List API"))
	})

	r.PathPrefix("/api/user").Handler(routes.UserRouter())
	r.PathPrefix("/api/task").Handler(routes.TaskRouter())
	log.Fatal(http.ListenAndServe(":8000", r))
	fmt.Println("Server started at port 8000")
}
