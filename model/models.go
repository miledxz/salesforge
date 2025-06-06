package model

type SequenceStep struct {
	ID          int    `db:"id" json:"id"`
	Subject     string `db:"subject" json:"subject"`
	Content     string `db:"content" json:"content"`
	WaitingDays int    `db:"waiting_days" json:"waitingDays"`
	SequenceID  int    `db:"sequence_id" json:"-"`
}

type Sequence struct {
	ID                   int            `db:"id" json:"id"`
	Name                 string         `db:"name" json:"name"`
	OpenTrackingEnabled  bool           `db:"open_tracking_enabled" json:"openTrackingEnabled"`
	ClickTrackingEnabled bool           `db:"click_tracking_enabled" json:"clickTrackingEnabled"`
	Steps                []SequenceStep `json:"steps"`
}
