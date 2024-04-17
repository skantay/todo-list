package repository

import (
	"log/slog"

	"go.mongodb.org/mongo-driver/mongo"
)

type Repository struct {
	TaskRepository taskRepository
}

type Collections struct {
	Task string
}

func New(client *mongo.Client, database string, collection Collections, log *slog.Logger) Repository {
	taskCollection := client.Database(database).Collection(collection.Task)

	return Repository{
		TaskRepository: newTaskRepository(taskCollection, log),
	}
}
