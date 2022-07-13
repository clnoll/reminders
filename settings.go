package app

import (
	"os"
)

var WhatsappAccountId = os.Getenv("WHATSAPP_ACCOUNT_ID")
var WhatsappToken = os.Getenv("WHATSAPP_TOKEN")

const ReminderTaskQueueName = "REMINDER_TASK_QUEUE"
const UpdateReminderSignalChannelName = "update-reminder-signal"
const TIME_FORMAT = "Mon Jan 2 2006 15:04:05 MST"
