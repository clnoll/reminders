package app

import (
	"fmt"
	"os"
)

var WhatsappAccountId = os.Getenv("WHATSAPP_ACCOUNT_ID")
var WhatsappSendMessageUrl = fmt.Sprintf("https://graph.facebook.com/v13.0/%s/messages", WhatsappAccountId)
var WhatsappToken = os.Getenv("WHATSAPP_TOKEN")
var WhatsappAuth = fmt.Sprintf("Bearer %s", WhatsappToken)
