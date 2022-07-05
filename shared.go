package app

import (
	"time"

	"go.temporal.io/sdk/workflow"
)

const ReminderTaskQueueName = "REMINDER_TASK_QUEUE"
const UpdateReminderSignalChannelName = "update-reminder-signal"

const TIME_FORMAT = "Mon Jan 2 2006 15:04:05 MST"

type ReminderDetails struct {
	CreatedAt    time.Time
	NMinutes     time.Duration
	ReminderTime time.Time
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

func (r *ReminderDetails) GetMinutesToReminder(ctx workflow.Context) time.Duration {
	return r.ReminderTime.Sub(workflow.Now(ctx))
}

func (r *ReminderDetails) GetReminderTime() time.Time {
	return r.ReminderTime
}
