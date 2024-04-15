package entity

import (
	"encoding/json"
	"fmt"
	"time"
)

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
	status   string   `json:"-"`
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
	return t.status
}

func (t *Task) SetStatusDone() {
	t.status = Done
}

func (t *Task) SetStatusActive() {
	t.status = Active
}

func (td TaskDate) Time() time.Time {
	return time.Time(td)
}

func (td *TaskDate) UnmarshalJSON(data []byte) error {
	var rawDate string
	if err := json.Unmarshal(data, &rawDate); err != nil {
		return fmt.Errorf("failed to unmarshal: %w",err)
	}
	parsedDate, err := time.Parse(dateFormat, rawDate)
	if err != nil {
		return fmt.Errorf("failed to parse: %w",err)
	}
	*td = TaskDate(parsedDate)
	return nil
}

func (td TaskDate) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprintf(`"%s"`, time.Time(td).Format(dateFormat))), nil
}
