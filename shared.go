package app

import "time"

const ReminderTaskQueue = "REMINDER_TASK_QUEUE"

const TIME_FORMAT = "Mon Jan 2 2006 15:04:05 MST"

type ReminderDetails struct {
	CreatedAt    time.Time
	NMinutes     time.Duration
	ReminderText string
	ReminderName string
	ReminderId   string
	Phone        string
}

type ReminderInput struct {
	ReminderId string
	RunId      string
	Phone      string
	NMinutes   int
}

func GetReminderTime(createdAt time.Time, duration time.Duration) time.Time {
	return createdAt.Add(duration)
}
