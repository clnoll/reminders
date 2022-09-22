package whatsapp

import (
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

func WhatsappRequestError(resp *http.Response) error {
	return errors.New(fmt.Sprintf("Error sending WhatsApp request. status=%s", resp.Status))
}

type IWhatsappClient interface {
	SendMessage(toPhone string, message string) error
}

type _LiveWhatsappClient struct {
	AuthToken string
	AccountId string
}

func (w _LiveWhatsappClient) SendMessage(toPhone string, message string) error {
	url := fmt.Sprintf("https://graph.facebook.com/v13.0/%s/messages", w.AccountId)
	auth := fmt.Sprintf("Bearer %s", w.AuthToken)

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
		return err
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	log.Println("WhatsApp response received. status:", resp.Status, "headers:", resp.Header, "body:", string(body))
	if resp.StatusCode >= 400 {
		log.Println("Error sending WhatsApp request")
		return WhatsappRequestError(resp)
	}
	return nil
}

type _MockWhatsappClient struct {
	AuthToken string
	AccountId string
}

func (f _MockWhatsappClient) SendMessage(toPhone string, message string) error {
	return nil
}
