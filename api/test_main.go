package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestCreateReminderHandler(t *testing.T) {
	req, err := http.NewRequest("POST", "/reminders", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(CreateReminderHandler)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusCreated)
	}
}
