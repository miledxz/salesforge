package service

import (
	"errors"
	"testing"

	"github.com/miledxz/salesforge/model"
	"github.com/stretchr/testify/assert"
)

type mockRepo struct {
	CreateSequenceFn func(*model.Sequence) (*model.Sequence, error)
	UpdateStepFn     func(int, *model.SequenceStep) error
	DeleteStepFn     func(int) error
	UpdateTrackingFn func(int, *bool, *bool) error
}

func (m *mockRepo) CreateSequence(seq *model.Sequence) (*model.Sequence, error) {
	return m.CreateSequenceFn(seq)
}

func (m *mockRepo) UpdateStep(id int, step *model.SequenceStep) error {
	return m.UpdateStepFn(id, step)
}

func (m *mockRepo) DeleteStep(id int) error {
	return m.DeleteStepFn(id)
}

func (m *mockRepo) UpdateTracking(id int, open, click *bool) error {
	return m.UpdateTrackingFn(id, open, click)
}

// --- Tests ---

func TestCreateSequence_Success(t *testing.T) {
	mock := &mockRepo{
		CreateSequenceFn: func(seq *model.Sequence) (*model.Sequence, error) {
			seq.ID = 42
			return seq, nil
		},
	}
	SetRepository(mock)

	input := &model.Sequence{Name: "Test"}
	result, err := CreateSequence(input)

	assert.NoError(t, err)
	assert.Equal(t, 42, result.ID)
}

func TestUpdateStep_Success(t *testing.T) {
	mock := &mockRepo{
		UpdateStepFn: func(id int, step *model.SequenceStep) error {
			if id != 1 {
				return errors.New("invalid id")
			}
			return nil
		},
	}
	SetRepository(mock)

	err := UpdateStep(1, &model.SequenceStep{Subject: "Updated"})
	assert.NoError(t, err)
}

func TestDeleteStep_Success(t *testing.T) {
	mock := &mockRepo{
		DeleteStepFn: func(id int) error {
			if id == 99 {
				return errors.New("not found")
			}
			return nil
		},
	}
	SetRepository(mock)

	err := DeleteStep(5)
	assert.NoError(t, err)
}

func TestUpdateTracking_Success(t *testing.T) {
	mock := &mockRepo{
		UpdateTrackingFn: func(id int, open, click *bool) error {
			if id == 0 {
				return errors.New("invalid id")
			}
			return nil
		},
	}
	SetRepository(mock)

	open := true
	click := false
	err := UpdateTracking(10, &open, &click)

	assert.NoError(t, err)
}
