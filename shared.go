package app

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
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

func SendViaWhatsapp(reminderDetails ReminderDetails) error {
	url := WhatsappSendMessageUrl
	auth := WhatsappAuth

	data := fmt.Sprintf(`{
		"messaging_product": "whatsapp",
  		"recipient_type": "individual",
  		"to": "%s",
  		"type": "text",
  		"text": {
			"body": "%s: %s"
		}
	}`, reminderDetails.Phone, reminderDetails.ReminderName, reminderDetails.ReminderText)
	log.Println("Sending WhatsApp reminder. url:", url, "data:", data)

	var query = []byte(data)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(query))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", auth)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Println("Panicked sending WhatsApp request")
		panic(err)
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	log.Println("WhatsApp response received. status:", resp.Status, "headers:", resp.Header, "body:", string(body))

	return nil
}
