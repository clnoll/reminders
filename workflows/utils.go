package workflows

import (
	"context"
	"log"
	"time"

	"go.temporal.io/sdk/client"

	"reminders/app"
)

func StartWorkflow(phone string, nMins int) (app.ReminderDetails, string, string) {
	// Create the client object just once per process
	c, err := client.NewClient(client.Options{})
	if err != nil {
		log.Fatalln("unable to create Temporal client", err)
	}
	defer c.Close()
	options := client.StartWorkflowOptions{
		ID:        "reminder-workflow",
		TaskQueue: app.ReminderTaskQueue,
	}
	reminderDetails := app.ReminderDetails{
		CreatedAt:    time.Now(),
		ReminderTime: time.Now().Add(time.Minute * time.Duration(nMins)),
		ReminderText: "Book return flights from Jakarta",
		ReminderName: "Flights",
		ReminderId:   "Test",
	}
	we, err := c.ExecuteWorkflow(context.Background(), options, MakeReminderWorkflow, reminderDetails)
	if err != nil {
		log.Fatalln("error starting Reminder workflow", err)
	}
	return reminderDetails, we.GetID(), we.GetRunID()
}

func printResults(reminderDetails app.ReminderDetails, workflowID, runID string) {
	log.Printf(
		"\nCreating reminder for %s (%s) at %s. ReminderId: %s\n",
		reminderDetails.ReminderName,
		reminderDetails.ReminderText,
		reminderDetails.ReminderTime.Format(app.TIME_FORMAT),
		reminderDetails.ReminderId,
	)
	log.Printf(
		"\nWorkflowID: %s RunID: %s\n",
		workflowID,
		runID,
	)
}
