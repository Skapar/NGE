package models

import (
	"fmt"
	"time"
)

// Event struct represents an event with a date and description
type Event struct {
	Date        time.Time `json:"date"`
	Description string    `json:"description"`
}

// Calendar struct represents a calendar with events
type Calendar struct {
	Events []Event `json:"events"`
}

// makeEvent marks a day on the calendar with a description
func (c *Calendar) makeEvent(date time.Time, description string) {
	event := Event{
		Date:        date,
		Description: description,
	}
	c.Events = append(c.Events, event)
}

// updateEvent updates the description of an existing event
func (c *Calendar) updateEvent(date time.Time, newDescription string) {
	for i, event := range c.Events {
		if event.Date.Equal(date) {
			c.Events[i].Description = newDescription
			return
		}
	}
	fmt.Println("Event not found on", date.Format("2006-01-02"))
}

// deleteEvent deletes an event from the calendar
func (c *Calendar) deleteEvent(date time.Time) {
	for i, event := range c.Events {
		if event.Date.Equal(date) {
			c.Events = append(c.Events[:i], c.Events[i+1:]...)
			return
		}
	}
	fmt.Println("Event not found on", date.Format("2006-01-02"))
}

// findEvent finds an event on a given date
func (c *Calendar) findEvent(date time.Time) *Event {
	for _, event := range c.Events {
		if event.Date.Equal(date) {
			return &event
		}
	}
	return nil
}

// printEvents prints all events in the calendar
func (c *Calendar) GetEvents() []string {
	var events []string
	for _, event := range c.Events {
		eventDescription := fmt.Sprintf("%s: %s", event.Date.Format("2006-01-02"), event.Description)
		events = append(events, eventDescription)
	}
	return events
}
