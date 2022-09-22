package utils

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"reminders/app/codec"

	"go.temporal.io/sdk/workflow"
)

type ReminderDetails struct {
	FromTime     time.Time
	NMinutes     time.Duration // editable
	ReminderTime time.Time     // automatically updates
	ReminderText string        // editable
	ReminderName string        // editable
	Phone        string
	WorkflowId   string
	RunId        string
	ReferenceId  string
}

type ReminderInput struct {
	FromTime     time.Time
	NMinutes     int
	ReminderText string
	ReminderName string
	Phone        string
	ReferenceId  string
}

type ReminderResponse struct {
	ReminderTime string
	ReminderText string
	ReminderName string
	ReferenceId  string
}

type UpdateReminderSignal struct {
	NMinutes     int
	ReminderText string
	ReminderName string
	Phone        string
}

func GetReminderTime(startTime time.Time, duration time.Duration) time.Time {
	return startTime.Add(duration)
}

func (r *ReminderDetails) GetMinutesToReminder(ctx workflow.Context) time.Duration {
	return r.ReminderTime.Sub(workflow.Now(ctx))
}

func (r *ReminderDetails) GetReminderTime() time.Time {
	return r.ReminderTime
}

func MakeReferenceId(workflowId string, runId string) (string, error) {
	if workflowId == "" || runId == "" {
		return "", errors.New(fmt.Sprintf("Unable to create referenceId from workflowId %s runId%s", workflowId, runId))
	}
	rawReferenceId := fmt.Sprintf("%s_%s", workflowId, runId)
	return codec.Encode(rawReferenceId), nil
}

func GetInternalIdsFromReferenceId(referenceId string) (string, string, error) {
	if referenceId == "" {
		return "", "", errors.New("Missing ReferenceId.")
	}
	decoded, err := codec.Decode(referenceId)
	if err != nil {
		return "", "", err
	}
	idComponents := strings.Split(decoded, "_")
	if len(idComponents) != 2 {
		return "", "", errors.New(fmt.Sprintf("Unable to decode referenceId %s", referenceId))
	}
	return idComponents[0], idComponents[1], nil
}
