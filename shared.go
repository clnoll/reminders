package app

import "time"

const ReminderTaskQueueName = "REMINDER_TASK_QUEUE"
const UpdateReminderSignalChannelName = "update-reminder-signal"

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

type CancelReminderSignal struct {
	WorkflowId string
	RunId      string
}

type UpdateReminderSignal struct {
	NMinutes     int
	ReminderText string
	ReminderName string
	Phone        string
}

func GetReminderTime(startTime time.Time, duration time.Duration) time.Time {
	return startTime.Add(duration)
}
}
