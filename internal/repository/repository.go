package repository

import "go.mongodb.org/mongo-driver/mongo"

type Repository struct {
	TaskRepository taskRepository
}

type Collections struct {
	Task string
}

func New(client *mongo.Client, database string, collection Collections) Repository {
	taskCollection := client.Database(database).Collection(collection.Task)

	return Repository{
		TaskRepository: newTaskRepository(taskCollection),
	}
}