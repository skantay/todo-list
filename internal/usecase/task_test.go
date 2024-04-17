package usecase

import "testing"

func TestTaskUsecase_Create(t *testing.T) {
	mockRepo := &mocks.MockTaskRepository{}
	usecase := NewTaskUsecase(mockRepo)

	// Mock repository behavior
	mockRepo.EXPECT().Create(gomock.Any(), gomock.Any()).Return("id123", nil)

	// Call the method being tested
	id, err := usecase.Create(context.Background(), "title", time.Now())

	// Assertions
	assert.NoError(t, err)
	assert.Equal(t, "id123", id)
}
