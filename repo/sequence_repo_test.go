package repo

import (
	"regexp"
	"testing"

	"github.com/miledxz/salesforge/db"
	"github.com/miledxz/salesforge/model"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
)

func setupMockDB(t *testing.T) (sqlmock.Sqlmock, func()) {
	sqlDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Error initializing sqlmock: %s", err)
	}
	db.DB = sqlx.NewDb(sqlDB, "postgres")

	return mock, func() {
		_ = db.DB.Close()
	}
}

func TestCreateSequence(t *testing.T) {
	mock, teardown := setupMockDB(t)
	defer teardown()

	repo := New()

	seq := &model.Sequence{
		Name:                 "Test Sequence",
		OpenTrackingEnabled:  true,
		ClickTrackingEnabled: false,
		Steps: []model.SequenceStep{
			{Subject: "Subject 1", Content: "Content 1", WaitingDays: 1},
			{Subject: "Subject 2", Content: "Content 2", WaitingDays: 2},
		},
	}

	mock.ExpectBegin()
	mock.ExpectQuery(`INSERT INTO sequences`).
		WithArgs(seq.Name, seq.OpenTrackingEnabled, seq.ClickTrackingEnabled).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))

	mock.ExpectQuery(`INSERT INTO sequence_steps`).
		WithArgs(1, "Subject 1", "Content 1", 1).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(10))

	mock.ExpectQuery(`INSERT INTO sequence_steps`).
		WithArgs(1, "Subject 2", "Content 2", 2).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(11))

	mock.ExpectCommit()

	result, err := repo.CreateSequence(seq)
	assert.NoError(t, err)
	assert.Equal(t, 1, result.ID)
	assert.Equal(t, 2, len(result.Steps))
	assert.Equal(t, 10, result.Steps[0].ID)
	assert.Equal(t, 11, result.Steps[1].ID)
}

func TestUpdateStep(t *testing.T) {
	mock, teardown := setupMockDB(t)
	defer teardown()

	repo := New()
	step := &model.SequenceStep{
		Subject: "Updated Subject",
		Content: "Updated Content",
	}

	mock.ExpectExec(regexp.QuoteMeta(`UPDATE sequence_steps SET subject=$1, content=$2 WHERE id=$3`)).
		WithArgs(step.Subject, step.Content, 5).
		WillReturnResult(sqlmock.NewResult(0, 1))

	err := repo.UpdateStep(5, step)
	assert.NoError(t, err)
}

func TestDeleteStep(t *testing.T) {
	mock, teardown := setupMockDB(t)
	defer teardown()

	repo := New()

	mock.ExpectExec(regexp.QuoteMeta(`DELETE FROM sequence_steps WHERE id=$1`)).
		WithArgs(7).
		WillReturnResult(sqlmock.NewResult(0, 1))

	err := repo.DeleteStep(7)
	assert.NoError(t, err)
}

func TestUpdateTracking(t *testing.T) {
	mock, teardown := setupMockDB(t)
	defer teardown()

	repo := New()
	open := true
	click := false

	mock.ExpectExec(regexp.QuoteMeta(
		`UPDATE sequences SET open_tracking_enabled=$1, click_tracking_enabled=$2 WHERE id=$3`)).
		WithArgs(open, click, 3).
		WillReturnResult(sqlmock.NewResult(0, 1))

	err := repo.UpdateTracking(3, &open, &click)
	assert.NoError(t, err)
}

func TestUpdateTracking_OnlyOpen(t *testing.T) {
	mock, teardown := setupMockDB(t)
	defer teardown()

	repo := New()
	open := true

	mock.ExpectExec(regexp.QuoteMeta(
		`UPDATE sequences SET open_tracking_enabled=$1 WHERE id=$2`)).
		WithArgs(open, 42).
		WillReturnResult(sqlmock.NewResult(0, 1))

	err := repo.UpdateTracking(42, &open, nil)
	assert.NoError(t, err)
}

func TestUpdateTracking_OnlyClick(t *testing.T) {
	mock, teardown := setupMockDB(t)
	defer teardown()

	repo := New()
	click := false

	mock.ExpectExec(regexp.QuoteMeta(
		`UPDATE sequences SET click_tracking_enabled=$1 WHERE id=$2`)).
		WithArgs(click, 55).
		WillReturnResult(sqlmock.NewResult(0, 1))

	err := repo.UpdateTracking(55, nil, &click)
	assert.NoError(t, err)
}

func TestUpdateTracking_None(t *testing.T) {
	repo := New()

	// Should skip execution if both open and click are nil
	err := repo.UpdateTracking(10, nil, nil)
	assert.NoError(t, err)
}
