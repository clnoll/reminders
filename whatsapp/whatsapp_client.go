package whatsapp

import (
	"errors"
	"fmt"
	"reminders/app"

	"golang.org/x/exp/slices"
)

func GetWhatsappClient() IWhatsappClient {
	if slices.Contains([]string{"PROD", "DEV"}, app.ENV) {
		return _LiveWhatsappClient{
			app.WhatsappToken,
			app.WhatsappAccountId,
		}
	} else {
		return _MockWhatsappClient{"", ""}
	}
}

func WhatsappClientError() error {
	return errors.New(fmt.Sprintf("Error getting WhatsappClient; unrecognized environment. ENV=%s", app.ENV))
}
