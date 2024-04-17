package entity

import (
	"testing"
	"time"
)

func TestTask(t *testing.T) {
	task := NewTask("title", time.Now())
	if task.GetStatus() != Active {
		t.Errorf("status must be \"active\", but got: %v", task.GetStatus())
	}

	task.SetStatusDone()

	if task.GetStatus() != Done {
		t.Errorf("status must be \"done\", but got: %v", task.GetStatus())
	}
}
