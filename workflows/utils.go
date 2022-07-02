package workflows

import (
	"context"
	"log"
	"time"

	"go.temporal.io/sdk/client"

	"reminders/app"
)

type CancelReminderSignal struct {
	WorkflowId string
	RunId      string
}

func StartWorkflow(phone string, nMins int) (app.ReminderDetails, string, string, error) {
	c, err := client.NewClient(client.Options{})
	if err != nil {
		log.Fatalln("unable to create Temporal client", err)
	}
	defer c.Close()
	options := client.StartWorkflowOptions{
		ID:        "reminder-workflow",
		TaskQueue: app.ReminderTaskQueue,
	}
	createdAt := time.Now()
	remindAt := time.Minute * time.Duration(nMins)
	reminderDetails := app.ReminderDetails{
		CreatedAt:    createdAt,
		NMinutes:     remindAt,
		ReminderText: "Book return flights from Jakarta",
		ReminderName: "Flights",
		ReminderId:   "Test",
	}
	we, err := c.ExecuteWorkflow(context.Background(), options, MakeReminderWorkflow, reminderDetails)
	if err != nil {
		log.Fatalln("error starting Reminder workflow", err)
	}
	return reminderDetails, we.GetID(), we.GetRunID(), err
}

func DeleteWorkflow(workflowId string, runId string) error {
	// Delete the reminder
	c, err := client.NewClient(client.Options{})
	if err != nil {
		log.Fatalln("unable to create Temporal client", err)
	}
	defer c.Close()
	err = c.CancelWorkflow(context.Background(), workflowId, runId)
	if err != nil {
		log.Fatalln("error deleting Reminder workflow", workflowId, runId, err)
	}
	return err
}

func printResults(reminderDetails app.ReminderDetails, workflowID, runID string) {
	log.Printf(
		"\nCreating reminder for %s (%s) at %s. ReminderId: %s\n",
		reminderDetails.ReminderName,
		reminderDetails.ReminderText,
		app.GetReminderTime(reminderDetails.CreatedAt, reminderDetails.NMinutes),
		reminderDetails.ReminderId,
	)
	log.Printf(
		"\nWorkflowID: %s RunID: %s\n",
		workflowID,
		runID,
	)
}
