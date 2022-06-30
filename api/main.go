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

	workflows.StartWorkflow(input.Phone, input.NMinutes)
	if err != nil {
		log.Printf("failed to start workflow: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	io.WriteString(w, `{"alive": true}`)
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
	w.WriteHeader(http.StatusOK)
	io.WriteString(w, `{"alive": true}`)
}

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/reminders", ReminderListHandler).Methods("GET")
	r.HandleFunc("/reminders", CreateReminderHandler).Methods("POST")
	r.HandleFunc("/reminders/{reminder_id}", GetReminderHandler).Methods("GET")
	r.HandleFunc("/reminders/{reminder_id}", UpdateReminderHandler).Methods("PATCH")
	r.HandleFunc("/reminders/{reminder_id}", DeleteReminderHandler).Methods("DELETE")
	http.Handle("/", r)

	log.Fatal(http.ListenAndServe(":8000", r))
}
