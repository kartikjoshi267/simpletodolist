package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/kartikjoshi267/simpletodolist/models"
	"go.mongodb.org/mongo-driver/bson"
)

func CreateUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Methods", "POST")
	var user models.User
	err := json.NewDecoder(r.Body).Decode(&user)

	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(bson.M{"error": "Invalid request body"})
		return
	}

	user.TaskIds = []string{}

	newUser, err := models.CreateUser(user)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(bson.M{"error": "Internal server error"})
		return
	}

	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(bson.M{
		"email": newUser.Email,
		"tasks": newUser.TaskIds,
	})
}

func LoginUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Methods", "POST")
	var user models.User
	err := json.NewDecoder(r.Body).Decode(&user)

	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(bson.M{"error": "Invalid request body"})
		return
	}

	token, err := models.LoginUser(user.Email, user.Password)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusUnauthorized)
		_ = json.NewEncoder(w).Encode(bson.M{"error": "Invalid email or password"})
		return
	}

	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(bson.M{
		"token": token,
	})
}

func GetUser(w http.ResponseWriter, r *http.Request) {
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

	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(bson.M{
		"email": user.Email,
		"tasks": user.TaskIds,
	})
}
