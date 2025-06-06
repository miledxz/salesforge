package service

import (
	"github.com/miledxz/salesforge/model"
	"github.com/miledxz/salesforge/repo"
)

type Service interface {
	CreateSequence(seq *model.Sequence) (*model.Sequence, error)
	UpdateStep(id int, step *model.SequenceStep) error
	DeleteStep(id int) error
	UpdateTracking(id int, open, click *bool) error
}

var repository repo.SequenceRepository = repo.New()

func SetRepository(r repo.SequenceRepository) {
	repository = r
}

func CreateSequence(seq *model.Sequence) (*model.Sequence, error) {
	return repository.CreateSequence(seq)
}

func UpdateStep(id int, step *model.SequenceStep) error {
	return repository.UpdateStep(id, step)
}

func DeleteStep(id int) error {
	return repository.DeleteStep(id)
}

func UpdateTracking(id int, open, click *bool) error {
	return repository.UpdateTracking(id, open, click)
}
