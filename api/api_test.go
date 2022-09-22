package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"reminders/app"
	"reminders/app/utils"
	"testing"
	"time"

	"github.com/gorilla/mux"
	"github.com/stretchr/testify/suite"
	"github.com/tidwall/gjson"
	"go.temporal.io/sdk/testsuite"
)

const FAKE_FROM_PHONE = "16505551111"

type UnitTestSuite struct {
	suite.Suite
	testsuite.WorkflowTestSuite

	env *testsuite.TestWorkflowEnvironment
}

func (s *UnitTestSuite) SetupTest() {
	s.env = s.NewTestWorkflowEnvironment()
}

func (s *UnitTestSuite) AfterTest(suiteName, testName string) {
	s.env.AssertExpectations(s.T())
}

func TestUnitTestSuite(t *testing.T) {
	suite.Run(t, new(UnitTestSuite))
}

func (t *UnitTestSuite) TestCreateReminderHandlerEmpty() {
	// Sending an empty body to /reminders results in an error response.
	req, err := http.NewRequest("POST", "/reminders", nil)
	if err != nil {
		t.Fail(err.Error())
	}

	cr := httptest.NewRecorder()
	m := mux.NewRouter()
	requestHandler := RequestHandler{utils.MockWorkflowClient{}}
	m.HandleFunc("/reminders", requestHandler.HandleCreate)
	m.ServeHTTP(cr, req)

	status := cr.Code
	t.True(status == http.StatusBadRequest, fmt.Sprintf("status = %v, expected %v", status, http.StatusBadRequest))
}

func (t *UnitTestSuite) TestCreateReminderHandler() {
	r := httptest.NewRecorder()
	m := mux.NewRouter()
	resp := createReminder(t, r, m)
	t.True(resp.ReminderTime != "", "Empty reminder time.")
}

func (t *UnitTestSuite) TestUpdateReminderHandler() {
	r := httptest.NewRecorder()
	m := mux.NewRouter()
	resp := createReminder(t, r, m)

	createReminderTime := resp.ReminderTime
	referenceId := resp.ReferenceId

	// Send a PUT to update the reminder to remind even earlier
	r = httptest.NewRecorder()
	updateReq := fmt.Sprintf(`{"NMinutes": 0}`)
	var query = []byte(updateReq)
	url := fmt.Sprintf("/reminders/%s", referenceId)
	requestHandler := RequestHandler{utils.MockWorkflowClient{}}
	m.HandleFunc("/reminders/{referenceId}", requestHandler.HandleUpdate)
	req, err := http.NewRequest("PUT", url, bytes.NewBuffer(query))
	if err != nil {
		t.Fail(err.Error())
	}

	m.ServeHTTP(r, req)
	if err != nil {
		t.Fail(err.Error())
	}

	if status := r.Code; status != http.StatusAccepted {
		t.Fail("UpdateReminderHandler returned wrong status code: got %v want %v",
			status, http.StatusAccepted)
	}

	updateRespBody, _ := ioutil.ReadAll(r.Body)
	updateReminderTime := gjson.GetBytes(updateRespBody, "ReminderTime").String()
	if updateReminderTime == "" {
		t.Fail("UpdateReminderHandler returned an empty ReminderTime.")
	}

	// The reminder should now be set to remind sooner than the original reminder
	createTs, _ := time.Parse(app.TIME_FORMAT, createReminderTime)
	updateTs, _ := time.Parse(app.TIME_FORMAT, updateReminderTime)
	t.True(updateTs.Before(createTs), fmt.Sprintf("Expected %v to be after %v", createTs, updateTs))
}

func (t *UnitTestSuite) TestDeleteReminderHandler() {
	var phone string
	t.env.RegisterDelayedCallback(func() {
		res, err := t.env.QueryWorkflow("getPhone")
		t.NoError(err)
		err = res.Get(&phone)
		t.NoError(err)
		t.Equal(phone, FAKE_FROM_PHONE)
	}, time.Millisecond*1)

	r := httptest.NewRecorder()
	m := mux.NewRouter()
	resp := createReminder(t, r, m)

	// Delete the reminder
	r = httptest.NewRecorder()
	url := fmt.Sprintf("/reminders/%s", resp.ReferenceId)
	workflowRequestHandler := RequestHandler{utils.MockWorkflowClient{}}
	m.HandleFunc("/reminders/{referenceId}", workflowRequestHandler.HandleDelete)
	req, err := http.NewRequest("DELETE", url, nil)
	if err != nil {
		t.Fail(err.Error())
	}

	m.ServeHTTP(r, req)
	if err != nil {
		t.Fail(err.Error())
	}

	status := r.Code
	t.True(status == http.StatusAccepted, fmt.Sprintf("status = %v, expected %v", status, http.StatusAccepted))
}

func (t *UnitTestSuite) TestWhatsappResponseHandlerCreate() {
	r := httptest.NewRecorder()
	m := mux.NewRouter()
	sendWhatsappMessageReminderRequest(t, r, m, whatsappCreateBody)
}

func (t *UnitTestSuite) TestWhatsappResponseHandlerUpdate() {
	r := httptest.NewRecorder()
	m := mux.NewRouter()
	sendWhatsappMessageReminderRequest(t, r, m, whatsappCreateBody)
	sendWhatsappMessageReminderRequest(t, r, m, whatsappUpdateBody)
}

func (t *UnitTestSuite) TestWhatsappResponseHandlerDelete() {
	r := httptest.NewRecorder()
	m := mux.NewRouter()
	sendWhatsappMessageReminderRequest(t, r, m, whatsappCreateBody)
	sendWhatsappMessageReminderRequest(t, r, m, whatsappDeleteBody)
}

func createReminder(t *UnitTestSuite, r *httptest.ResponseRecorder, m *mux.Router) utils.ReminderResponse {
	body := fmt.Sprintf(`{
		"NMinutes": 1,
  		"ReminderText": "Book return flight",
  		"ReminderName": "Flights",
  		"Phone": "%s"
	}`, FAKE_FROM_PHONE)
	requestHandler := RequestHandler{utils.MockWorkflowClient{}}
	status, resp := post(t, r, m, "/reminders", requestHandler.HandleCreate, body)
	t.True(status == http.StatusCreated, fmt.Sprintf("status %v, expected %v", status, http.StatusCreated))
	return resp
}

func sendWhatsappMessageReminderRequest(t *UnitTestSuite, r *httptest.ResponseRecorder, m *mux.Router, body string) {
	requestHandler := RequestHandler{utils.MockWorkflowClient{}}
	status, _ := post(t, r, m, "/external/reminders/whatsapp", requestHandler.HandleWhatsappCreate, body)
	t.True(status == http.StatusOK, fmt.Sprintf("status %v, expected %v", status, http.StatusOK))
}

func post(
	t *UnitTestSuite, r *httptest.ResponseRecorder, m *mux.Router,
	url string, handler func(http.ResponseWriter, *http.Request,
	), body string) (int, utils.ReminderResponse) {
	var query = []byte(body)
	m.HandleFunc(url, handler)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(query))
	if err != nil {
		t.Fail(err.Error())
	}
	m.ServeHTTP(r, req)

	var resp utils.ReminderResponse
	err = json.NewDecoder(r.Body).Decode(&resp)
	if err != nil {
		return r.Code, utils.ReminderResponse{}
	}
	return r.Code, resp
}

var whatsappCreateBody = fmt.Sprintf(`{
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
				  "wa_id": "%s"
				}
			  ],
			  "messages": [
				{
				  "from": "%s",
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
}`, FAKE_FROM_PHONE, FAKE_FROM_PHONE)

var whatsappUpdateBody = fmt.Sprintf(`{
	"object": "whatsapp_business_account",
	"entry": [
		{
		"id": "0",
		"changes": [
			{
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
					"wa_id": "%s"
				}
				],
				"messages": [
				{
					"context": {
					"from": "%s",
					"id": "ABGGFlA5Fpa"
					},
					"from": "%s",
					"id": "wamid.HBgLMTUwMjc0MTI0ODAVAgASGBQzRUIwMkJDMDQ5MkNCMzc1NUY0NgA=",
					"timestamp": "1657724009",
					"text": {
					"body": "Update XXXXXXX: 0h 5m"
					},
					"type": "text"
				}
				]
			},
			"field": "messages"
			}
		]
		}
	]
	}
`, FAKE_FROM_PHONE, FAKE_FROM_PHONE, FAKE_FROM_PHONE)

var whatsappDeleteBody = fmt.Sprintf(`{
	"object": "whatsapp_business_account",
	"entry": [
		{
		"id": "0",
		"changes": [
			{
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
					"wa_id": "%s"
				}
				],
				"messages": [
				{
					"context": {
					"from": "%s",
					"id": "ABGGFlA5Fpa"
					},
					"from": "%s",
					"id": "wamid.HBgLMTUwMjc0MTI0ODAVAgASGBQzRUIwMkJDMDQ5MkNCMzc1NUY0NgA=",
					"timestamp": "1657724009",
					"text": {
					"body": "Update XXXXXXX"
					},
					"type": "text"
				}
				]
			},
			"field": "messages"
			}
		]
		}
	]
	}
`, FAKE_FROM_PHONE, FAKE_FROM_PHONE, FAKE_FROM_PHONE)
