package models

import (
	"time"

	"gorm.io/gorm"
)

// Event struct represents an event with a date and description
type Event struct {
	gorm.Model
	Id          uint      `json:"id" gorm:"unique;primaryKey;autoIncrement"`
	Date        time.Time `json:"date"`
	Description string    `json:"description" `
	UserID      uint      `json:"user_id" `
}

// Calendar struct represents a calendar with events

func AddEvent(db *gorm.DB, date time.Time, description string) error {
	event := Event{
		Date:        date,
		Description: description,
	}
	return db.Create(&event).Error
}
