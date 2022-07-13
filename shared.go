package app

import (
	"time"

	"go.temporal.io/sdk/workflow"
)

type ReminderDetails struct {
	FromTime     time.Time
	NMinutes     time.Duration // editable
	ReminderTime time.Time     // automatically updates
	ReminderText string        // editable
	ReminderName string        // editable
	Phone        string
	WorkflowId   string
	RunId        string
	ReferenceId  string
}

type ReminderInput struct {
	FromTime     time.Time
	NMinutes     int
	ReminderText string
	ReminderName string
	Phone        string
	ReferenceId  string
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
