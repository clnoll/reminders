package main

import (
	"context"
	"log"
	"time"

	"go.temporal.io/sdk/client"

	"reminders/app"
	"reminders/app/utils"
	"reminders/app/workflows"
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
		TaskQueue: app.ReminderTaskQueueName,
	}
	reminderDetails := utils.ReminderDetails{
		FromTime:     time.Now(),
		NMinutes:     time.Second * 60,
		ReminderText: "Book return flights from Jakarta",
		ReminderName: "Flights",
	}
	we, err := c.ExecuteWorkflow(context.Background(), options, workflows.MakeReminderWorkflow, reminderDetails)
	if err != nil {
		log.Fatalln("error starting Reminder workflow", err)
	}
	printResults(reminderDetails, we.GetID(), we.GetRunID())
}

// @@@SNIPEND

func printResults(reminderDetails utils.ReminderDetails, workflowId, runId string) {
	log.Printf(
		"\nCreating reminder for %s (%s) at %s. workflowId=%s runId=%s\n",
		reminderDetails.ReminderName,
		reminderDetails.ReminderText,
		utils.GetReminderTime(reminderDetails.FromTime, reminderDetails.NMinutes),
		workflowId,
		runId,
	)
	log.Printf(
		"\nWorkflowID: %s RunID: %s\n",
		workflowId,
		runId,
	)
}
