package models

import "time"

type Event struct {
	ID     string
	UserID string
	Date   time.Time
	Text   string
}
