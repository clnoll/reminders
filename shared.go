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

type ReminderDetails struct {
	CreatedAt    time.Time
	NMinutes     time.Duration // editable
	ReminderTime time.Time     // automatically updates
	ReminderText string        // editable
	ReminderName string        // editable
	Phone        string        // editable
	WorkflowId   string
	RunId        string
}

type ReminderInput struct {
	NMinutes     int
	ReminderText string
	ReminderName string
	Phone        string
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

func SendWhatsappMessage(toPhone string, message string) error {
	url := WhatsappSendMessageUrl
	auth := WhatsappAuth

	data := fmt.Sprintf(`{
		"messaging_product": "whatsapp",
  		"recipient_type": "individual",
  		"to": "%s",
  		"type": "text",
  		"text": {
			"body": "%s"
		}
	}`, toPhone, message)
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
