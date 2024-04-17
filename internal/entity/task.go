package entity

import (
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Определение общих ошибок для сущности "Задача"
var (
	ErrAlreadyExists = errors.New("task already exists")
	ErrTaskNotFound  = errors.New("task does not exist")
	ErrInvalidTitle  = errors.New("invalid title")
	ErrInvalidStatus = errors.New("invalid status")
	ErrInvalidID     = errors.New("invalid id")
)

// Константы для статусов задачи и формата даты
const (
	Active     = "active"
	Done       = "done"
	dateFormat = "2006-01-02"
)

// TaskDate определяет пользовательский тип для даты задачи
type TaskDate time.Time

type Task struct {
	ID       string   `json:"id"`
	Title    string   `json:"title"`
	ActiveAt TaskDate `json:"activeAt"`
	Status   string   `json:"-"`
}

// NewTask создает новую задачу
// Со значением status="active"
func NewTask(title string, activeAt TaskDate) Task {
	t := Task{
		Title:    title,
		ActiveAt: activeAt,
	}

	t.SetStatusActive()

	return t
}

func (t *Task) SetStatusDone() {
	t.Status = Done
}

func (t *Task) SetStatusActive() {
	t.Status = Active
}

func (td TaskDate) Time() time.Time {
	return time.Time(td)
}

// UnmarshalJSON разбирает JSON TaskDate
func (td *TaskDate) UnmarshalJSON(data []byte) error {
	var rawDate string
	if err := json.Unmarshal(data, &rawDate); err != nil {
		return fmt.Errorf("failed to unmarshal: %w", err)
	}

	parsedDate, err := time.Parse(dateFormat, rawDate)
	if err != nil {
		return fmt.Errorf("failed to parse: %w", err)
	}

	*td = TaskDate(parsedDate)
	
	return nil
}

// MarshalJSON преобразует TaskDate в JSON
func (td TaskDate) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprintf(`"%s"`, time.Time(td).Format(dateFormat))), nil
}

// UnmarshalBSON разбирает BSON Task
func (t *Task) UnmarshalBSON(data []byte) error {
	var rawTask struct {
		ID       primitive.ObjectID `bson:"_id"`
		Title    string             `bson:"title"`
		ActiveAt time.Time          `bson:"activeAt"`
		Status   string             `bson:"status"`
	}

	if err := bson.Unmarshal(data, &rawTask); err != nil {
		return fmt.Errorf("failed to unmarshal Task: %w", err)
	}

	t.ID = rawTask.ID.Hex()
	
	t.Title = rawTask.Title
	
	t.ActiveAt = TaskDate(rawTask.ActiveAt)
	
	t.Status = rawTask.Status
	
	return nil
}

// MarshalBSON преобразует Task в BSON
func (t Task) MarshalBSON() ([]byte, error) {
	id, err := primitive.ObjectIDFromHex(t.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to convert ObjectId: %w", err)
	}

	return bson.Marshal(struct {
		ID       primitive.ObjectID `bson:"_id"`
		Title    string             `bson:"title"`
		ActiveAt time.Time          `bson:"activeAt"`
		Status   string             `bson:"status"`
	}{
		ID:       id,
		Title:    t.Title,
		ActiveAt: time.Time(t.ActiveAt),
		Status:   t.Status,
	})
}
