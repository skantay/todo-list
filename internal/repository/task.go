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

func (t taskRepository) Create(ctx context.Context, task entity.Task) (string, error) {
	existingTask, err := t.findTask(ctx, task.Title, task.ActiveAt)
	if err != nil {
		return "", fmt.Errorf("failed to check task uniqueness: %w", err)
	}
	if existingTask {
		return "", entity.ErrAlreadyExists
	}

	taskBSON := bson.M{
		"title":    task.Title,
		"activeAt": task.ActiveAt.Time(),
		"status":   task.Status,
	}

	result, err := t.collection.InsertOne(ctx, taskBSON)
	if err != nil {
		return "", fmt.Errorf("failed to insert a task into db: %w", err)
	}

	oidResult, ok := result.InsertedID.(primitive.ObjectID)
	if !ok {
		return "", errors.New("failed to retrieve inserted ID")
	}

	return oidResult.Hex(), nil
}

func (t taskRepository) findTask(ctx context.Context, title string, activeAt entity.TaskDate) (bool, error) {
	filter := bson.M{
		"title":    title,
		"activeAt": activeAt.Time(),
	}

	result := t.collection.FindOne(ctx, filter)
	if err := result.Err(); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return false, nil
		}
		return false, fmt.Errorf("failed to find task: %w", err)
	}

	if result.Err() == mongo.ErrNoDocuments {
		return false, nil
	}

	return true, nil
}

type sTask struct {
	ID       primitive.ObjectID `bson:"_id"`
	ActiveAt time.Time
	Title    string
	Status   string
}

func (t taskRepository) List(ctx context.Context, status string, now time.Time) ([]entity.Task, error) {
	filter := bson.M{
		"status":   status,
		"activeAt": bson.M{"$lte": now},
	}

	sort := bson.D{{"activeAt", 1}}

	cursor, err := t.collection.Find(ctx, filter, options.Find().SetSort(sort))
	if err != nil {
		return nil, fmt.Errorf("failed to cursor a collection: %w", err)
	}
	defer cursor.Close(ctx)

	var tasks []entity.Task

	for cursor.Next(ctx) {
		var task sTask
		if err := cursor.Decode(&task); err != nil {
			return nil, fmt.Errorf("failed to decode task: %w", err)
		}

		taskEntity := entity.Task{
			ID:       task.ID.Hex(),
			Status:   task.Status,
			ActiveAt: entity.TaskDate(task.ActiveAt),
			Title:    task.Title,
		}

		tasks = append(tasks, taskEntity)
	}

	if err := cursor.Err(); err != nil {
		return nil, fmt.Errorf("cursor error occured: %w", err)
	}

	return tasks, nil
}

func (t taskRepository) Update(ctx context.Context, task entity.Task) error {
	existingTask, err := t.findTask(ctx, task.Title, task.ActiveAt)
	if err != nil {
		return fmt.Errorf("failed to check task uniqueness: %w", err)
	}
	if existingTask {
		return entity.ErrAlreadyExists
	}

	id, err := primitive.ObjectIDFromHex(task.ID)
	if err != nil {
		return fmt.Errorf("failed to convert ObjectId: %w", err)
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

func (t taskRepository) MarkDone(ctx context.Context, id string) error {
	idObj, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return fmt.Errorf("failed to convert ObjectId: %w", err)
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

func (t taskRepository) Delete(ctx context.Context, id string) error {
	idObj, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return fmt.Errorf("failed to convert ObjectId: %w", err)
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
