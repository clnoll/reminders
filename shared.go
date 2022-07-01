package app

import "time"

// @@@SNIPSTART reminders/app
const ReminderTaskQueue = "REMINDER_TASK_QUEUE"

// @@@SNIPEND

const TIME_FORMAT = "Mon Jan 2 2006 15:04:05 MST"

type ReminderDetails struct {
	CreatedAt    time.Time
	ReminderTime time.Time
	ReminderText string
	ReminderName string
	ReminderId   string
	Phone        string
}

type ReminderInput struct {
	Phone    string
	NMinutes int
}
