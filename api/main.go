package main

import (
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func ReminderListHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "RemindersList: %s\n", "")
}

func CreateReminderHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
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
