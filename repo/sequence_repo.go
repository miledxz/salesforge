package repo

import (
	"fmt"
	"strings"

	"github.com/miledxz/salesforge/db"
	"github.com/miledxz/salesforge/model"
)

type SequenceRepository interface {
	CreateSequence(seq *model.Sequence) (*model.Sequence, error)
	UpdateStep(id int, step *model.SequenceStep) error
	DeleteStep(id int) error
	UpdateTracking(id int, open, click *bool) error
}

type repoImpl struct{}

func New() SequenceRepository {
	return &repoImpl{}
}

func (r *repoImpl) CreateSequence(seq *model.Sequence) (*model.Sequence, error) {
	tx := db.DB.MustBegin()
	err := tx.QueryRowx(
		`INSERT INTO sequences (name, open_tracking_enabled, click_tracking_enabled) 
		 VALUES ($1, $2, $3) RETURNING id`,
		seq.Name, seq.OpenTrackingEnabled, seq.ClickTrackingEnabled,
	).Scan(&seq.ID)
	if err != nil {
		tx.Rollback()
		return nil, err
	}
	for i := range seq.Steps {
		err = tx.QueryRowx(
			`INSERT INTO sequence_steps (sequence_id, subject, content, waiting_days) 
			 VALUES ($1, $2, $3, $4) RETURNING id`,
			seq.ID, seq.Steps[i].Subject, seq.Steps[i].Content, seq.Steps[i].WaitingDays,
		).Scan(&seq.Steps[i].ID)
		if err != nil {
			tx.Rollback()
			return nil, err
		}
	}
	tx.Commit()
	return seq, nil
}

func (r *repoImpl) UpdateStep(id int, step *model.SequenceStep) error {
	_, err := db.DB.Exec(`UPDATE sequence_steps SET subject=$1, content=$2 WHERE id=$3`, step.Subject, step.Content, id)
	return err
}

func (r *repoImpl) DeleteStep(id int) error {
	_, err := db.DB.Exec(`DELETE FROM sequence_steps WHERE id=$1`, id)
	return err
}

func (r *repoImpl) UpdateTracking(id int, open, click *bool) error {
	fields := []string{}
	params := []interface{}{}
	idx := 1

	if open != nil {
		fields = append(fields, fmt.Sprintf("open_tracking_enabled=$%d", idx))
		params = append(params, *open)
		idx++
	}

	if click != nil {
		fields = append(fields, fmt.Sprintf("click_tracking_enabled=$%d", idx))
		params = append(params, *click)
		idx++
	}

	if len(fields) == 0 {
		return nil
	}

	query := fmt.Sprintf("UPDATE sequences SET %s WHERE id=$%d", strings.Join(fields, ", "), idx)
	params = append(params, id)

	_, err := db.DB.Exec(query, params...)
	return err
}
