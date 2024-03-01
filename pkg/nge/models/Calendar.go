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

func DeleteEvent(db *gorm.DB, eventID uint) error {
	result := db.Delete(&Event{}, eventID)
	return result.Error // Returns nil if the deletion is successful
}

func UpdateEvent(db *gorm.DB, eventID uint, newDate time.Time, newDescription string) error {
	result := db.Model(&Event{}).Where("id = ?", eventID).Updates(Event{Date: newDate, Description: newDescription})
	return result.Error // Returns nil if the update is successful
}

func GetEventByID(db *gorm.DB, eventID uint) (Event, error) {
	var event Event
	result := db.First(&event, eventID)
	return event, result.Error // Returns the event and nil if found, else an error
}

func GetAllEvents(db *gorm.DB, userID uint) ([]Event, error) {
	var events []Event
	query := db.Model(&Event{})
	if userID != 0 {
		query = query.Where("user_id = ?", userID)
	}
	result := query.Find(&events)
	return events, result.Error // Returns the list of events and nil if successful, else an error
}
