package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"reminders/app"
	"testing"
	"time"

	"github.com/gorilla/mux"
	"github.com/tidwall/gjson"
)

func TestCreateReminderHandlerEmpty(t *testing.T) {
	// Sending an empty body to /reminders results in an error response.
	req, err := http.NewRequest("POST", "/reminders", nil)
	if err != nil {
		t.Fatal(err)
	}

	cr := httptest.NewRecorder()
	m := mux.NewRouter()
	m.HandleFunc("/reminders", CreateReminderHandler)
	m.ServeHTTP(cr, req)

	if status := cr.Code; status != http.StatusBadRequest {
		t.Errorf("CreateReminderHandler returned wrong status code: got %v want %v",
			status, http.StatusCreated)
	}
}

func TestCreateReminderHandler(t *testing.T) {
	r := httptest.NewRecorder()
	m := mux.NewRouter()
	createReminder(t, r, m)
	respBody, _ := ioutil.ReadAll(r.Body)
	if results := gjson.GetBytes(respBody, "reminderTime").String(); results == "" {
		t.Errorf("CreateReminderHandler returned an empty reminderTime")
	}
}

func TestUpdateReminderHandler(t *testing.T) {
	r := httptest.NewRecorder()
	m := mux.NewRouter()
	createReminder(t, r, m)
	createRespBody, _ := ioutil.ReadAll(r.Body)

	results := gjson.GetManyBytes(createRespBody, "reminderTime", "workflowId", "runId")
	createReminderTime := results[0].String()
	workflowId := results[1].String()
	runId := results[2].String()

	// Send a PUT to update the reminder to remind even earlier
	r = httptest.NewRecorder()
	updateReq := fmt.Sprintf(`{"NMinutes": 0}`)
	var query = []byte(updateReq)
	url := fmt.Sprintf("/reminders/%s/%s", workflowId, runId)
	m.HandleFunc("/reminders/{workflowId}/{runId}", UpdateReminderHandler)
	req, err := http.NewRequest("PUT", url, bytes.NewBuffer(query))
	if err != nil {
		t.Fatal(err)
	}

	m.ServeHTTP(r, req)
	if err != nil {
		t.Fatal(err)
	}

	if status := r.Code; status != http.StatusAccepted {
		t.Errorf("UpdateReminderHandler returned wrong status code: got %v want %v",
			status, http.StatusAccepted)
	}

	updateRespBody, _ := ioutil.ReadAll(r.Body)
	updateReminderTime := gjson.GetBytes(updateRespBody, "reminderTime").String()
	if updateReminderTime == "" {
		t.Errorf("UpdateReminderHandler returned an empty reminderTime.")
	}

	// The reminder should now be set to remind sooner than the original reminder
	createTs, _ := time.Parse(app.TIME_FORMAT, createReminderTime)
	updateTs, _ := time.Parse(app.TIME_FORMAT, updateReminderTime)
	if !(updateTs.Before(createTs)) {
		t.Errorf("UpdateReminderHandler did not update the reminder: got %v want less than %v",
			updateReminderTime, createReminderTime)
	}
}

func TestDeleteReminderHandler(t *testing.T) {
	r := httptest.NewRecorder()
	m := mux.NewRouter()
	createReminder(t, r, m)
	createRespBody, _ := ioutil.ReadAll(r.Body)

	results := gjson.GetManyBytes(createRespBody, "workflowId", "runId")
	workflowId := results[0].String()
	runId := results[1].String()

	// Delete the reminder
	r = httptest.NewRecorder()
	url := fmt.Sprintf("/reminders/%s/%s", workflowId, runId)
	m.HandleFunc("/reminders/{workflowId}/{runId}", DeleteReminderHandler)
	req, err := http.NewRequest("DELETE", url, nil)
	if err != nil {
		t.Fatal(err)
	}

	m.ServeHTTP(r, req)
	if err != nil {
		t.Fatal(err)
	}

	if status := r.Code; status != http.StatusAccepted {
		t.Errorf("DeleteReminderHandler returned wrong status code: got %v want %v",
			status, http.StatusAccepted)
	}
}

func TestWhatsappResponseHandlerCreate(t *testing.T) {
	r := httptest.NewRecorder()
	m := mux.NewRouter()
	createReminderFromWhatsappMessage(t, r, m)
}

func createReminder(t *testing.T, r *httptest.ResponseRecorder, m *mux.Router) {
	body := fmt.Sprintf(`{
		"NMinutes": 1,
  		"ReminderText": "Book return flight",
  		"ReminderName": "Flights",
  		"Phone": "15555555555"
	}`)
	post(t, r, m, "/reminders", CreateReminderHandler, body, http.StatusCreated)
}

func createReminderFromWhatsappMessage(t *testing.T, r *httptest.ResponseRecorder, m *mux.Router) {
	body := fmt.Sprintf(`{
		"object": "whatsapp_business_account",
		"entry": [
		  {
			"id": "0",
			"changes": [
			  {
				"field": "messages",
				"value": {
				  "messaging_product": "whatsapp",
				  "metadata": {
					"display_phone_number": "16505551111",
					"phone_number_id": "123456123"
				  },
				  "contacts": [
					{
					  "profile": {
						"name": "test user name"
					  },
					  "wa_id": "16315551181"
					}
				  ],
				  "messages": [
					{
					  "from": "16315551181",
					  "id": "ABGGFlA5Fpa",
					  "timestamp": "1504902988",
					  "type": "text",
					  "text": {
						"body": "New Reminder Family: call mom about test results: 3h 30m"
					  }
					}
				  ]
				}
			  }
			]
		  }
		]
	}`)
	post(t, r, m, "/external/reminders/whatsapp", WhatsappResponseHandler, body, http.StatusOK)
}

func post(
	t *testing.T, r *httptest.ResponseRecorder, m *mux.Router,
	url string, handler func(http.ResponseWriter, *http.Request,
	), body string, statusExpected int) {
	var query = []byte(body)
	m.HandleFunc(url, handler)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(query))
	if err != nil {
		t.Fatal(err)
	}
	m.ServeHTTP(r, req)

	if status := r.Code; status != statusExpected {
		t.Errorf("WhatsappResponseHandler returned wrong status code: got %v want %v",
			status, statusExpected)
	}
}
