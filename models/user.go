package models

import (
	"context"
	"fmt"

	"github.com/kartikjoshi267/simpletodolist/database"
	"github.com/kartikjoshi267/simpletodolist/lib"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID       primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Email    string             `json:"email,omitempty" bson:"email,omitempty"`
	Password string             `json:"password,omitempty" bson:"password,omitempty"`
	TaskIds  []string           `json:"tasks,omitempty" bson:"tasks,omitempty"`
}

var ctx = context.TODO()

func isValidUser(user *User) bool {
	return user.Email != "" && user.Password != ""
}

func CreateUser(user User) (*User, error) {
	if !isValidUser(&user) {
		fmt.Println(user)
		return nil, fmt.Errorf("invalid user")
	}

	if database.UsersCollection.FindOne(ctx, bson.M{"email": user.Email}).Err() == nil {
		return nil, fmt.Errorf("User with email %s already exists", user.Email)
	}

	password, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	user.ID = primitive.NewObjectID()
	user.TaskIds = []string{}

	user.Password = string(password)
	_, err = database.UsersCollection.InsertOne(ctx, user)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func LoginUser(email string, password string) (string, error) {
	var user User
	err := database.UsersCollection.FindOne(ctx, bson.M{"email": email}).Decode(&user)
	if err != nil {
		return "", fmt.Errorf("User with email %s not found", email)
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))

	if err != nil {
		return "", fmt.Errorf("invalid password")
	}

	token, err := lib.CreateToken(user.ID.Hex())
	if err != nil {
		return "", err
	}
	return token, nil
}

func GetUser(token string) (*User, error) {
	userId, err := lib.VerifyToken(token)
	if err != nil {
		return nil, err
	}

	var user User
	objectId, err := primitive.ObjectIDFromHex(userId)
	if err != nil {
		return nil, err
	}

	err = database.UsersCollection.FindOne(ctx, bson.M{"_id": objectId}).Decode(&user)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func AddTaskToUser(taskId string, userId string) error {
	objectId, _ := primitive.ObjectIDFromHex(userId)
	_, err := database.UsersCollection.UpdateOne(ctx, bson.M{"_id": objectId}, bson.M{"$push": bson.M{"tasks": taskId}})
	if err != nil {
		return err
	}

	return nil
}

func DeleteTaskFromUser(taskId string, userId string) error {
	objectId, _ := primitive.ObjectIDFromHex(userId)
	_, err := database.UsersCollection.UpdateOne(ctx, bson.M{"_id": objectId}, bson.M{"$pull": bson.M{"tasks": taskId}})
	if err != nil {
		return err
	}

	return nil
}
