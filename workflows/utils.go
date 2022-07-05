package workflows

import (
	"context"
	"log"
	"reminders/app"
	"time"

	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/workflow"
)

func StartWorkflow(phone string, nMins int) (app.ReminderDetails, string, string, error) {
	c, err := client.NewClient(client.Options{})
	if err != nil {
		log.Fatalln("unable to create Temporal client", err)
	}
	defer c.Close()
	options := client.StartWorkflowOptions{
		ID:        "reminder-workflow",
		TaskQueue: app.ReminderTaskQueueName,
	}
	createdAt := time.Now()
	remindInMinutes := time.Minute * time.Duration(nMins)
	reminderDetails := app.ReminderDetails{
		CreatedAt:    createdAt,
		NMinutes:     remindInMinutes,
		ReminderTime: createdAt.Add(remindInMinutes),
		ReminderText: "Book return flights from Jakarta",
		ReminderName: "Flights",
		ReminderId:   "Test",
	}
	log.Println("Starting workflow to remind in", remindInMinutes, "minutes, at", reminderDetails.GetReminderTime().Format(app.TIME_FORMAT))
	we, err := c.ExecuteWorkflow(context.Background(), options, MakeReminderWorkflow, reminderDetails)
	if err != nil {
		log.Fatalln("error starting Reminder workflow", err)
	}
	return reminderDetails, we.GetID(), we.GetRunID(), err
}

func UpdateWorkflow(workflowId string, runId string, phone string, nMinutes int) (app.ReminderDetails, error) {
	c, err := client.NewClient(client.Options{})
	if err != nil {
		log.Fatalln("unable to create Temporal client", err)
	}
	defer c.Close()

	signal := app.UpdateReminderSignal{
		Phone:    phone,
		NMinutes: nMinutes,
	}
	ctx := context.Background()
	var reminderDetails app.ReminderDetails
	err = c.SignalWorkflow(ctx, workflowId, runId, app.UpdateReminderSignalChannelName, signal)
	if err != nil {
		log.Fatalln("Error sending the UpdateReminder Signal", err)
		return reminderDetails, err
	}
	return reminderDetails, err
}

func updateReminderDetails(ctx workflow.Context, reminderUpdate *app.UpdateReminderSignal, reminderDetails *app.ReminderDetails) *app.ReminderDetails {
	newReminderTime := time.Duration(reminderUpdate.NMinutes) * time.Minute
	reminderDetails.NMinutes = newReminderTime
	reminderDetails.ReminderTime = app.GetReminderTime(workflow.Now(ctx), newReminderTime)
	if reminderUpdate.Phone != "" {
		reminderDetails.Phone = reminderUpdate.Phone
	}
	if reminderUpdate.ReminderText != "" {
		reminderDetails.ReminderText = reminderUpdate.ReminderText
	}
	if reminderUpdate.ReminderName != "" {
		reminderDetails.ReminderName = reminderUpdate.ReminderName
	}
	return reminderDetails
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
