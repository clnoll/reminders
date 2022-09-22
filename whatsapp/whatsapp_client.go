package whatsapp

import (
	"errors"
	"fmt"
	"reminders/app"
)

var WhatsappClient = GetWhatsappClient()

func GetWhatsappClient() IWhatsappClient {
	if app.ENV == "PROD" {
		return _LiveWhatsappClient{}
	} else {
		return _MockWhatsappClient{}
	}
}

func WhatsappClientError() error {
	return errors.New(fmt.Sprintf("Error getting WhatsappClient; unrecognized environment. ENV=%s", app.ENV))
}
