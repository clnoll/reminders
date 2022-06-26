package app

// @@@SNIPSTART reminders/app
const ReminderTaskQueue = "REMINDER_TASK_QUEUE"

// @@@SNIPEND

type ReminderDetails struct {
	CreatedAt    float32
	ReminderTime float32
	ReminderText string
	ReminderName string
	ReminderId   string
}
