package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
	"github.com/kartikjoshi267/simpletodolist/models"
	"go.mongodb.org/mongo-driver/bson"
)

func CreateTask(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Methods", "POST")

	token := r.Header.Get("Authorization")
	token = strings.Split(strings.TrimSpace(token), " ")[1]

	user, err := models.GetUser(token)
	if err != nil {
		fmt.Println(err, token)
		w.WriteHeader(http.StatusUnauthorized)
		_ = json.NewEncoder(w).Encode(bson.M{"error": "Unauthorized"})
		return
	}

	var task models.Task
	err = json.NewDecoder(r.Body).Decode(&task)

	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(bson.M{"error": "Invalid request body"})
		return
	}

	task.User = user.ID
	task.Completed = false

	newTaskId, err := models.CreateTask(task)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(bson.M{"error": "Internal server error"})
		return
	}

	err = models.AddTaskToUser(newTaskId.Hex(), user.ID.Hex())
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(bson.M{"error": "Internal server error"})
		return
	}

	newTask, err := models.GetTask(newTaskId.Hex())
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(bson.M{"error": "Internal server error"})
		return
	}

	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(bson.M{
		"_id":       newTaskId.Hex(),
		"text":      newTask.Text,
		"completed": newTask.Completed,
	})
}

func GetTasks(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Methods", "GET")

	token := r.Header.Get("Authorization")
	token = strings.Split(strings.TrimSpace(token), " ")[1]

	user, err := models.GetUser(token)
	if err != nil {
		fmt.Println(err, token)
		w.WriteHeader(http.StatusUnauthorized)
		_ = json.NewEncoder(w).Encode(bson.M{"error": "Unauthorized"})
		return
	}

	tasks, err := models.GetTasks(user.ID.Hex())
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(bson.M{"error": "Something went wrong"})
		return
	}

	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(tasks)
}

func DeleteTask(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Methods", "DELETE")

	token := r.Header.Get("Authorization")
	token = strings.Split(strings.TrimSpace(token), " ")[1]

	user, err := models.GetUser(token)
	if err != nil {
		fmt.Println(err, token)
		w.WriteHeader(http.StatusUnauthorized)
		_ = json.NewEncoder(w).Encode(bson.M{"error": "Unauthorized"})
		return
	}

	params := mux.Vars(r)
	id := params["id"]

	if id == "" {
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(bson.M{"error": "Task Id is required"})
		return
	}

	err = models.DeleteTask(id)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(bson.M{"error": "Something went wrong"})
		return
	}

	err = models.DeleteTaskFromUser(id, user.ID.Hex())
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(bson.M{"error": "Internal server error"})
		return
	}

	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(bson.M{
		"message": "Task was successfully deleted",
	})
}

func UpdateTask(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Methods", "PUT")

	token := r.Header.Get("Authorization")
	token = strings.Split(strings.TrimSpace(token), " ")[1]
	user, err := models.GetUser(token)
	if err != nil {
		fmt.Println(err, token)
		w.WriteHeader(http.StatusUnauthorized)
		_ = json.NewEncoder(w).Encode(bson.M{"error": "Unauthorized"})
		return
	}

	params := mux.Vars(r)
	id := params["id"]

	if id == "" {
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(bson.M{"error": "Task Id is required"})
		return
	}

	task, err := models.GetTask(id)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(bson.M{"error": "Task not found"})
		return
	}

	if task.User.Hex() != user.ID.Hex() {
		w.WriteHeader(http.StatusUnauthorized)
		_ = json.NewEncoder(w).Encode(bson.M{"error": "Unauthorized"})
		return
	}

	err = models.UpdateTask(id)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(bson.M{"error": "Something went wrong"})
		return
	}

	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(bson.M{
		"_id":       task.ID.Hex(),
		"text":      task.Text,
		"completed": !task.Completed,
	})
}
