package main

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"reminders/app"
	"reminders/app/workflows"
	"strconv"
	"time"

	"github.com/gorilla/mux"
	"github.com/tidwall/gjson"
)

func ReminderListHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "RemindersList: %s\n", "")
}

func CreateReminderHandler(w http.ResponseWriter, r *http.Request) {
	if r.Body == nil {
		http.Error(w, "Bad request.", http.StatusBadRequest)
		return
	}

	var input app.ReminderInput

	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	input.FromTime = time.Now()
	reminderInfo, err := workflows.StartWorkflow(input)
	log.Printf("Creating reminder for Phone %s", input.Phone)
	if err != nil {
		log.Printf("failed to start workflow: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	log.Printf("Created reminder for workflowId %s runId %s", reminderInfo.WorkflowId, reminderInfo.RunId)

	w.WriteHeader(http.StatusCreated)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(
		map[string]string{
			"workflowId":   reminderInfo.WorkflowId,
			"runId":        reminderInfo.RunId,
			"reminderTime": app.GetReminderTime(reminderInfo.FromTime, reminderInfo.NMinutes).Format(app.TIME_FORMAT),
		})
}

func GetReminderHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	io.WriteString(w, `{"alive": true}`)
}

func UpdateReminderHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	runId := vars["runId"]
	workflowId := vars["workflowId"]
	var input app.ReminderInput

	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	reminderInfo, err := workflows.UpdateWorkflow(workflowId, runId, input)
	if err != nil {
		log.Printf("Failed to update workflow: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	log.Printf("Updated reminder for workflowId %s runId %s", workflowId, runId)
	w.WriteHeader(http.StatusAccepted)
	w.Header().Set("Content-Type", "application/json")

	json.NewEncoder(w).Encode(
		map[string]string{
			"workflowId":   reminderInfo.WorkflowId,
			"runId":        reminderInfo.RunId,
			"reminderTime": app.GetReminderTime(reminderInfo.FromTime, reminderInfo.NMinutes).Format(app.TIME_FORMAT),
		})
}

func DeleteReminderHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	runId := vars["runId"]
	workflowId := vars["workflowId"]

	err := workflows.DeleteWorkflow(workflowId, runId)
	if err != nil {
		log.Printf("failed to delete workflow %s (runID %s): %v", workflowId, runId, err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	log.Printf("Deleted reminder for workflowId %s runId %s", workflowId, runId)
	w.WriteHeader(http.StatusAccepted)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(
		map[string]string{
			"workflowId": workflowId,
			"runId":      runId,
			"status":     "cancelled",
		})
}

func WhatsappResponseHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("WhatsApp message received.")
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Unable to read request body.", http.StatusBadRequest)
		return
	}

	results := gjson.GetManyBytes(
		body,
		"entry.0.changes.0.value.messages.0.from",
		"entry.0.changes.0.value.messages.0.timestamp",
		"entry.0.changes.0.value.messages.0.text.body",
	)

	fromPhone := results[0].Str
	timestampStr := results[1].Str
	message := results[2].Str

	if fromPhone == "" {
		http.Error(w, "From phone number not found in request.", http.StatusBadRequest)
		return
	}

	timestampInt, err := strconv.ParseInt(timestampStr, 10, 64)
	if err != nil {
		http.Error(w, "Invalid WhatsApp request.", http.StatusBadRequest)
	}
	fromTime := time.Unix(timestampInt, 0)

	err = doMessageAction(fromPhone, message, fromTime)

	if err != nil {
		app.SendWhatsappMessage(fromPhone, "Unable to create reminder; unrecognized request format.")
		http.Error(w, "Unrecognized reminder request format.", http.StatusBadRequest)
	}

	w.WriteHeader(http.StatusOK)
}

func doMessageAction(phone string, message string, fromTime time.Time) error {
	if name, text, nMinutes, err := app.ParseCreateReminderMessage(message); err == nil {
		return createReminderFromMessage(phone, name, text, nMinutes, fromTime)
	}
	return app.ReminderParseError(fmt.Sprintf("Unable to create reminder from request %s", message))
}

func createReminderFromMessage(phone string, reminderName string, reminderText string, nMinutes int, fromTime time.Time) error {
	input := app.ReminderInput{
		FromTime:     fromTime,
		NMinutes:     nMinutes,
		ReminderText: reminderText,
		ReminderName: reminderName,
		Phone:        phone,
	}
	reminderInfo, err := workflows.StartWorkflow(input)
	log.Printf("Creating reminder for Phone %s", input.Phone)
	if err != nil {
		log.Printf("failed to start workflow: %v", err)
		return err
	}
	log.Printf("Created reminder for workflowId %s runId %s", reminderInfo.WorkflowId, reminderInfo.RunId)
	err = app.SendWhatsappMessage(
		phone,
		fmt.Sprintf(
			"Created reminder %s: %s at %s. workflowId=%s runId=%s",
			reminderInfo.ReminderName, reminderInfo.ReminderText,
			app.GetReminderTime(reminderInfo.FromTime, reminderInfo.NMinutes).Format(app.TIME_FORMAT),
			reminderInfo.WorkflowId, reminderInfo.RunId,
		),
	)
	return err
}

func sendErrorMessage(phone string, message string) {
	app.SendWhatsappMessage(phone, fmt.Sprintf(
		`Error creating reminder: "%s". Please use the format "New Reminder <Reminder Name>: <Reminder Text>: <1H 30M>"`,
		message,
	))
}

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/reminders", ReminderListHandler).Methods("GET")
	r.HandleFunc("/reminders", CreateReminderHandler).Methods("POST")
	r.HandleFunc("/reminders/{workflowId}/{runId}", GetReminderHandler).Methods("GET")
	r.HandleFunc("/reminders/{workflowId}/{runId}", UpdateReminderHandler).Methods("PUT")
	r.HandleFunc("/reminders/{workflowId}/{runId}", DeleteReminderHandler).Methods("DELETE")
	r.HandleFunc("/external/reminders/whatsapp", WhatsappResponseHandler).Methods("POST")
	http.Handle("/", r)

	log.Fatal(http.ListenAndServe(":8000", r))
}
