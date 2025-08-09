package models

import "time"

type Event struct {
	ID     string    `json:"id" form:"id"`
	UserID string    `json:"user_id" form:"user_id"`
	Date   time.Time `json:"date" form:"date"`
	Text   string    `json:"text" form:"text"`
}
