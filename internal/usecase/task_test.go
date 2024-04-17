package usecase

import (
	"context"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/skantay/todo-list/internal/entity"
	"github.com/stretchr/testify/assert"
)

func TestTaskUsecase_Create(t *testing.T) {
	usecase := NewMocktaskRepo(gomock.NewController(t))

	// Mock repository behavior
	usecase.EXPECT().Create(gomock.Any(), gomock.Any()).Return("id123", nil)

	task := entity.NewTask("title", entity.TaskDate(time.Now()))
	id, err := usecase.Create(context.Background(), task)

	if id != "id123" {
		t.Error(id)
	}

	// Assertions
	assert.NoError(t, err)
	assert.Equal(t, "id123", id)
}
