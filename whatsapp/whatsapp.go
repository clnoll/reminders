package whatsapp

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"reminders/app"
)

var WhatsappAuth = fmt.Sprintf("Bearer %s", app.WhatsappToken)
var SendMessageUrl = fmt.Sprintf("https://graph.facebook.com/v13.0/%s/messages", app.WhatsappAccountId)

func SendMessage(toPhone string, message string) error {
	url := SendMessageUrl
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
