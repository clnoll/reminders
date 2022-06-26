package main

import (
	"context"
	"log"

	"go.temporal.io/sdk/client"

	"reminders/app"
)

// @@@SNIPSTART reminders-start-workflow
func main() {
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
		CreatedAt:    1656191661,
		ReminderTime: 1656278061,
		ReminderText: "Book return flights from Jakarta",
		ReminderName: "Flights",
		ReminderId:   "Test",
	}
	we, err := c.ExecuteWorkflow(context.Background(), options, app.MakeReminderWorkflow, reminderDetails)
	if err != nil {
		log.Fatalln("error starting Reminder workflow", err)
	}
	printResults(reminderDetails, we.GetID(), we.GetRunID())
}

// @@@SNIPEND

func printResults(reminderDetails app.ReminderDetails, workflowID, runID string) {
	log.Printf(
		"\nCreating reminder for %s (%s) at %f. ReminderId: %s\n",
		reminderDetails.ReminderName,
		reminderDetails.ReminderText,
		reminderDetails.ReminderTime,
		reminderDetails.ReminderId,
	)
	log.Printf(
		"\nWorkflowID: %s RunID: %s\n",
		workflowID,
		runID,
	)
}
