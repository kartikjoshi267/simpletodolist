package models

import (
	"fmt"

	"github.com/kartikjoshi267/simpletodolist/database"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Task struct {
	ID        primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Text      string             `json:"text,omitempty" bson:"text,omitempty"`
	User      primitive.ObjectID `json:"user,omitempty" bson:"user,omitempty"`
	Completed bool               `json:"completed" bson:"completed"`
}

func isValidTask(task Task) bool {
	return task.Text != ""
}

func CreateTask(task Task) (primitive.ObjectID, error) {
	if !isValidTask(task) {
		return primitive.ObjectID{}, fmt.Errorf("invalid task")
	}

	inserted, err := database.TasksCollection.InsertOne(ctx, task)
	if err != nil {
		return primitive.ObjectID{}, err
	}

	insertedID, ok := inserted.InsertedID.(primitive.ObjectID)
	if !ok {
		return primitive.ObjectID{}, fmt.Errorf("failed to convert insertedID to primitive.ObjectID")
	}

	return insertedID, nil
}

func GetTasks(userId string) ([]Task, error) {
	var tasks []Task
	objectId, _ := primitive.ObjectIDFromHex(userId)
	cursor, err := database.TasksCollection.Find(ctx, bson.M{"user": objectId})
	if err != nil {
		return nil, err
	}
	err = cursor.All(ctx, &tasks)
	if err != nil {
		return nil, err
	}

	return tasks, nil
}

func GetTask(taskId string) (*Task, error) {
	var task Task
	objectId, _ := primitive.ObjectIDFromHex(taskId)
	err := database.TasksCollection.FindOne(ctx, bson.M{"_id": objectId}).Decode(&task)
	if err != nil {
		return nil, err
	}

	return &task, nil
}

func DeleteTask(taskId string) error {
	objectId, _ := primitive.ObjectIDFromHex(taskId)
	_, err := database.TasksCollection.DeleteOne(ctx, bson.M{"_id": objectId})
	if err != nil {
		return err
	}

	return nil
}

func UpdateTask(taskId string) error {
	task, err := GetTask(taskId)
	if err != nil {
		return err
	}

	task.Completed = !task.Completed

	objectId, _ := primitive.ObjectIDFromHex(taskId)
	_, err = database.TasksCollection.UpdateOne(ctx, bson.M{"_id": objectId}, bson.M{"$set": bson.M{"completed": task.Completed}})
	if err != nil {
		return err
	}

	return nil
}
