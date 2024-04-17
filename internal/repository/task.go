package repository

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"github.com/skantay/todo-list/internal/entity"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type taskRepository struct {
	collection *mongo.Collection
	log        *slog.Logger
}

func newTaskRepository(collection *mongo.Collection, log *slog.Logger) taskRepository {
	return taskRepository{
		collection: collection,
		log:        log,
	}
}

// Create создаёт задачу в коллекции.
func (t taskRepository) Create(ctx context.Context, task entity.Task) (string, error) {
	existingTask, err := t.findTask(ctx, task)
	if err != nil {
		return "", fmt.Errorf("failed to check task uniqueness: %w", err)
	}
	if existingTask {
		return "", entity.ErrAlreadyExists
	}

	task.ID = primitive.NewObjectID().Hex()

	result, err := t.collection.InsertOne(ctx, task)
	if err != nil {
		return "", fmt.Errorf("failed to insert a task into db: %w", err)
	}

	oidResult, ok := result.InsertedID.(primitive.ObjectID)
	if !ok {
		return "", errors.New("failed to retrieve inserted ID")
	}

	return oidResult.Hex(), nil
}

// List возвращает список задач с колекции на основе указанных параметров(status, now time.Time).
func (t taskRepository) List(ctx context.Context, status string, now time.Time) ([]entity.Task, error) {
	var filter bson.M

	if status == entity.Active {
		filter = bson.M{
			"status":   status,
			"activeAt": bson.M{"$lte": now},
		}
	} else {
		filter = bson.M{
			"status": status,
		}
	}

	sort := bson.D{{"activeAt", 1}}

	cursor, err := t.collection.Find(ctx, filter, options.Find().SetSort(sort))
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

// Update обновляет title и activeAt задачи в колекции на основе указанных параметров(task entity.Task).
func (t taskRepository) Update(ctx context.Context, task entity.Task) error {
	// Конвертируем строку ID в тип ObjectID
	id, err := primitive.ObjectIDFromHex(task.ID)
	if err != nil {
		return entity.ErrInvalidID
	}

	// Проверка на наличие задачи с такими же полями
	// Если присутсвуют тогда возврщает entity.ErrAlreadyExists
	existingTask, err := t.findTask(ctx, task)
	if err != nil {
		return fmt.Errorf("failed to check task uniqueness: %w", err)
	}
	if existingTask {
		return entity.ErrAlreadyExists
	}

	filter := bson.M{"_id": id}

	update := bson.M{
		"$set": bson.M{
			"title":    task.Title,
			"activeAt": task.ActiveAt.Time(),
		},
	}

	result, err := t.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return fmt.Errorf("update failed: %w", err)
	}

	if result.ModifiedCount == 0 {
		return entity.ErrTaskNotFound
	}

	return nil
}

// MarkDone маркирует задачу завершенной в колекции на основе указанных параметров(id).
func (t taskRepository) MarkDone(ctx context.Context, id string) error {
	// Конвертируем строку ID в тип ObjectID
	idObj, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return entity.ErrInvalidID
	}

	filter := bson.M{"_id": idObj}

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
		return entity.ErrTaskNotFound
	}

	return nil
}

// Delete удаляет задачу в колекции на основе указанных параметров(id).
func (t taskRepository) Delete(ctx context.Context, id string) error {
	// Конвертируем строку ID в тип ObjectID
	idObj, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return entity.ErrInvalidID
	}

	filter := bson.M{"_id": idObj}

	result, err := t.collection.DeleteOne(ctx, filter)
	if err != nil {
		return fmt.Errorf("failed to delete: %w", err)
	}

	if result.DeletedCount == 0 {
		return entity.ErrTaskNotFound
	}

	return nil
}

// findTask ищет задачу в коллекции на основе указанных параметров(title, activeAt).
func (t taskRepository) findTask(ctx context.Context, task entity.Task) (bool, error) {
	filter := bson.M{
		"title":    task.Title,
		"activeAt": task.ActiveAt.Time(),
		"status":   task.Status,
	}

	t.log.Debug("", "task", filter)

	result := t.collection.FindOne(ctx, filter)
	if err := result.Err(); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return false, nil
		}
		return false, fmt.Errorf("failed to find task: %w", err)
	}

	return true, nil
}
