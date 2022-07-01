package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"reminders/app"
	"reminders/app/workflows"

	"github.com/gorilla/mux"
)

func ReminderListHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "RemindersList: %s\n", "")
}

func CreateReminderHandler(w http.ResponseWriter, r *http.Request) {
	var input app.ReminderInput

	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	reminderInfo, workflowId, runId, err := workflows.StartWorkflow(input.Phone, input.NMinutes)
	if err != nil {
		log.Printf("failed to start workflow: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(
		map[string]string{
			"workflowId":   workflowId,
			"runId":        runId,
			"reminderId":   reminderInfo.ReminderId,
			"reminderTime": reminderInfo.ReminderTime.Format(app.TIME_FORMAT),
		})
}

func GetReminderHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	io.WriteString(w, `{"alive": true}`)
}

func UpdateReminderHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	io.WriteString(w, `{"alive": true}`)
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
	w.WriteHeader(http.StatusCreated)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(
		map[string]string{
			"workflowId": workflowId,
			"runId":      runId,
			"status":     "cancelled",
		})
}

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/reminders", ReminderListHandler).Methods("GET")
	r.HandleFunc("/reminders", CreateReminderHandler).Methods("POST")
	r.HandleFunc("/reminders/{workflowId}/{runId}", GetReminderHandler).Methods("GET")
	r.HandleFunc("/reminders/{workflowId}/{runId}", UpdateReminderHandler).Methods("PATCH")
	r.HandleFunc("/reminders/{workflowId}/{runId}", DeleteReminderHandler).Methods("DELETE")
	http.Handle("/", r)

	log.Fatal(http.ListenAndServe(":8000", r))
}
