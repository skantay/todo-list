package usecase

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/skantay/todo-list/internal/entity"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/mongo"
)

func Test_Create(t *testing.T) {
	type args struct {
		ctx      context.Context
		title    string
		activeAt entity.TaskDate
	}

	type fields struct {
		taskRepo *MocktaskRepo
	}

	set := func(field *fields, id string, err error) {
		field.taskRepo.EXPECT().Create(gomock.Any(), gomock.Any()).Return(id, err)
	}

	tests := []struct {
		name    string
		setup   func(f *fields)
		args    args
		wantID  string
		wantErr error
	}{
		{
			name: "#1 valid",
			setup: func(f *fields) {
				set(f, "1", nil)
			},
			args: args{
				ctx:      context.Background(),
				title:    "valid",
				activeAt: entity.TaskDate(time.Now()),
			},
			wantID:  "1",
			wantErr: nil,
		},
		{
			name: "#2 invlaid title",
			setup: nil,
			args: args{
				ctx:      context.Background(),
				title:    "loooooooooooooooooooo123ooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooo000000000000000000000000000000000000000000oong",
				activeAt: entity.TaskDate(time.Now()),
			},
			wantID:  "",
			wantErr: entity.ErrInvalidTitle,
		},
		{
			name: "#3 repository error",
			setup: func(f *fields) {
				set(f, "", mongo.ErrNoDocuments)
			},
			args: args{
				ctx:      context.Background(),
				title:    "title",
				activeAt: entity.TaskDate(time.Now()),
			},
			wantID:  "",
			wantErr: mongo.ErrNoDocuments,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			taskRepo := NewMocktaskRepo(ctrl)

			taskUsecase := newTaskUsecase(taskRepo, nil)

			fields := &fields{taskRepo}

			if tt.setup != nil {
				tt.setup(fields)
			}

			id, err := taskUsecase.Create(tt.args.ctx, tt.args.title, tt.args.activeAt)

			assert.Equal(t, tt.wantID, id)
			if tt.wantErr != nil {
				if err == nil {
					t.Errorf("\nexpected error: %v \nbut got nil error", tt.wantErr)
				} else {
					if !errors.Is(err, tt.wantErr) {
						t.Errorf("\nexpected error:%v \ninvalid error: %v", tt.wantErr.Error(), err.Error())
					}
				}
			} else if err != nil {
				t.Errorf("\nunexpeceted error: %v", err)
			}
			ctrl.Finish()
		})
	}
}
