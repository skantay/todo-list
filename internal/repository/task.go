package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/skantay/todo-list/internal/entity"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type taskRepository struct {
	collection *mongo.Collection
}

func newTaskRepository(collection *mongo.Collection) taskRepository {
	return taskRepository{
		collection: collection,
	}
}

func (t taskRepository) Create(ctx context.Context, task entity.Task) (string, error) {
	// Check if the task already exists
	existingTask, err := t.getByTitle(ctx, task.Title)
	if err != nil {
		return "", fmt.Errorf("failed to check task uniqueness: %w", err)
	}
	if existingTask != nil {
		return "", fmt.Errorf("task with title '%s' already exists", task.Title)
	}

	// If the task is unique, insert it into the database
	result, err := t.collection.InsertOne(ctx, task)
	if err != nil {
		return "", fmt.Errorf("failed to insert a task into db: %w", err)
	}

	if oidResult, ok := result.InsertedID.(primitive.ObjectID); ok {
		return oidResult.Hex(), nil
	}

	return primitive.NilObjectID.Hex(), errors.New("failed to retrieve inserted ID")
}

func (t taskRepository) getByTitle(ctx context.Context, title string, activeAt TaskDg) (bool, error) {
	filter := bson.M{
		"title":title,
		"activeAt":activeAt
	}
} 

func (t taskRepository) List(ctx context.Context, status string, now time.Time) ([]entity.Task, error) {
	filter := bson.M{
		"status":   status,
		"activeAt": bson.M{"$lte": now},
	}

	cursor, err := t.collection.Find(ctx, filter)
	if err != nil {
		return nil, fmt.Errorf("failed to cursor a collection: %w", err)
	}
	defer cursor.Close(ctx)

	var tasks []entity.Task
	for cursor.Next(ctx) {
		var task entity.Task
		if err := cursor.Decode(&task); err != nil {
			return nil, fmt.Errorf("failed to decode task: %w", err)
		}
		tasks = append(tasks, task)
	}

	if err := cursor.Err(); err != nil {
		return nil, fmt.Errorf("cursor error occured: %w", err)
	}

	return tasks, nil
}

func (t taskRepository) Update(ctx context.Context, task entity.Task) error {
	filter := bson.M{"_id": task.ID}

	update := bson.M{
		"$set": bson.M{
			"title":    task.Title,
			"activeAt": task.ActiveAt,
		},
	}

	result, err := t.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return fmt.Errorf("update failed: %w", err)
	}

	if result.ModifiedCount == 0 {
		return mongo.ErrNoDocuments
	}

	return nil
}

func (t taskRepository) MarkDone(ctx context.Context, id string) error {
	filter := bson.M{"_id": id}
	update := bson.M{
		"$set": bson.M{
			"status": entity.Done,
		},
	}

	result, err := t.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return fmt.Errorf("update failed: %w", err)
	}

	if result.ModifiedCount == 0 {
		return mongo.ErrNoDocuments
	}

	return nil
}

func (t taskRepository) Delete(ctx context.Context, id string) error {
	filter := bson.M{"_id": id}

	result, err := t.collection.DeleteOne(ctx, filter)
	if err != nil {
		return fmt.Errorf("failed to delete: %w", err)
	}

	if result.DeletedCount == 0 {
		return mongo.ErrNoDocuments
	}

	return nil
}
