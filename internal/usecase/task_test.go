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
			name:  "#2 invlaid title",
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

func Test_UpdateTask(t *testing.T) {
	type args struct {
		ctx  context.Context
		task entity.Task
	}

	type fields struct {
		taskRepo *MocktaskRepo
	}

	set := func(field *fields, err error) {
		field.taskRepo.EXPECT().Update(gomock.Any(), gomock.Any()).Return(err)
	}

	tests := []struct {
		name    string
		setup   func(f *fields)
		args    args
		wantErr error
	}{
		{
			name: "#1 valid",
			setup: func(f *fields) {
				set(f, nil)
			},
			args: args{
				ctx:  context.Background(),
				task: entity.NewTask("title", entity.TaskDate(time.Now())),
			},
			wantErr: nil,
		},
		{
			name:  "#2 invlaid title",
			setup: nil,
			args: args{
				ctx:  context.Background(),
				task: entity.NewTask("loooooooooooooooooooo123ooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooo000000000000000000000000000000000000000000oong", entity.TaskDate(time.Now())),
			},
			wantErr: entity.ErrInvalidTitle,
		},
		{
			name: "#3 repository error",
			setup: func(f *fields) {
				set(f, mongo.ErrEmptySlice)
			},
			args: args{
				ctx:  context.Background(),
				task: entity.NewTask("title", entity.TaskDate(time.Now())),
			},
			wantErr: mongo.ErrEmptySlice,
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

			err := taskUsecase.UpdateTask(tt.args.ctx, tt.args.task)
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

func Test_MarkTaskDone(t *testing.T) {
	type args struct {
		ctx context.Context
		id  string
	}

	type fields struct {
		taskRepo *MocktaskRepo
	}

	set := func(field *fields, err error) {
		field.taskRepo.EXPECT().MarkDone(gomock.Any(), gomock.Any()).Return(err)
	}

	tests := []struct {
		name    string
		setup   func(f *fields)
		args    args
		wantErr error
	}{
		{
			name: "#1 valid",
			setup: func(f *fields) {
				set(f, nil)
			},
			args: args{
				ctx: context.Background(),
				id:  "1",
			},
			wantErr: nil,
		},
		{
			name: "#2 repository error",
			setup: func(f *fields) {
				set(f, mongo.ErrEmptySlice)
			},
			args: args{
				ctx: context.Background(),
				id:  "1",
			},
			wantErr: mongo.ErrEmptySlice,
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

			err := taskUsecase.MarkTaskDone(tt.args.ctx, tt.args.id)
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

func Test_Delete(t *testing.T) {
	type args struct {
		ctx context.Context
		id  string
	}

	type fields struct {
		taskRepo *MocktaskRepo
	}

	set := func(field *fields, err error) {
		field.taskRepo.EXPECT().Delete(gomock.Any(), gomock.Any()).Return(err)
	}

	tests := []struct {
		name    string
		setup   func(f *fields)
		args    args
		wantErr error
	}{
		{
			name: "#1 valid",
			setup: func(f *fields) {
				set(f, nil)
			},
			args: args{
				ctx: context.Background(),
				id:  "1",
			},
			wantErr: nil,
		},
		{
			name: "#2 repository error",
			setup: func(f *fields) {
				set(f, mongo.ErrEmptySlice)
			},
			args: args{
				ctx: context.Background(),
				id:  "1",
			},
			wantErr: mongo.ErrEmptySlice,
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

			err := taskUsecase.Delete(tt.args.ctx, tt.args.id)
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
func Test_List(t *testing.T) {
	type args struct {
		ctx    context.Context
		status string
	}

	type fields struct {
		taskRepo *MocktaskRepo
	}

	set := func(field *fields, tasks []entity.Task, err error) {
		field.taskRepo.EXPECT().List(gomock.Any(), gomock.Any(), gomock.Any()).Return(tasks, err)
	}

	tests := []struct {
		name      string
		setup     func(f *fields)
		args      args
		wantTasks []entity.Task
		wantErr   error
	}{
		{
			name: "#1 valid",
			setup: func(f *fields) {
				set(
					f,
					[]entity.Task{
						{
							Title: "BTC",
						},
					},
					nil,
				)
			},
			args: args{
				ctx: context.Background(),
				status: "active",
			},
			wantTasks: []entity.Task{
				{
					Title: "BTC",
				},
			},
			wantErr: nil,
		},
		{
			name: "#2 valid empty status",
			setup: func(f *fields) {
				set(
					f,
					[]entity.Task{
						{
							Title: "BTC",
						},
					},
					nil,
				)
			},
			args: args{
				ctx: context.Background(),
				status: "",
			},
			wantTasks: []entity.Task{
				{
					Title: "BTC",
				},
			},
			wantErr: nil,
		},
		{
			name: "#3 invalid status",
			setup: nil,
			args: args{
				ctx: context.Background(),
				status: "invalid",
			},
			wantTasks: nil,
			wantErr: entity.ErrInvalidStatus,
		},
		{
			name: "#4 repository error",
			setup: func(f *fields) {
				set(
					f,
					nil,
					mongo.ErrClientDisconnected,
				)
			},
			args: args{
				ctx: context.Background(),
				status: "active",
			},
			wantTasks: nil,
			wantErr: mongo.ErrClientDisconnected,
		},
		{
			name: "#5 weekend usecase",
			setup: func(f *fields) {
				set(
					f,
					[]entity.Task{
						{
							Title: "BTC",
							ActiveAt: entity.TaskDate(time.Date(2024,04,14,1,1,1,1,time.Local)),
						},
					},
					nil,
				)
			},
			args: args{
				ctx: context.Background(),
				status: "active",
			},
			wantTasks: []entity.Task{
				{
					Title: "ВЫХОДНОЙ - BTC",
					ActiveAt: entity.TaskDate(time.Date(2024,04,14,1,1,1,1,time.Local)),
				},
			},
			wantErr: nil,
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

			tasks, err := taskUsecase.List(tt.args.ctx, tt.args.status)
			assert.Equal(t, tt.wantTasks, tasks)
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
