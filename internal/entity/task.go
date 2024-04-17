package entity

import (
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
)

var ErrAlreadyExists = errors.New("task already exists")

const (
	Active     = "active"
	Done       = "done"
	dateFormat = "2006-01-02"
)

type TaskDate time.Time

type Task struct {
	ID       string   `json:"id"`
	Title    string   `json:"title"`
	ActiveAt TaskDate `json:"active_at"`
	Status   string   `json:"-"`
}

func NewTask(title string, activeAt TaskDate) Task {
	t := Task{
		Title:    title,
		ActiveAt: activeAt,
	}
	t.SetStatusActive()

	return t
}

func (t *Task) GetStatus() string {
	return t.Status
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

func (td TaskDate) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprintf(`"%s"`, time.Time(td).Format(dateFormat))), nil
}

func (td *TaskDate) UnmarshalBSON(data []byte) error {
	var t time.Time
	if err := bson.Unmarshal(data, &t); err != nil {
		return err
	}
	*td = TaskDate(t)
	return nil
}
