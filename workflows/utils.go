package workflows

import (
	"context"
	"log"
	"reminders/app"
	"time"

	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/workflow"
)

func StartWorkflow(input app.ReminderInput) (app.ReminderDetails, error) {
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
	remindInMinutes := time.Minute * time.Duration(input.NMinutes)
	reminderDetails := app.ReminderDetails{
		CreatedAt:    createdAt,
		NMinutes:     remindInMinutes,
		Phone:        input.Phone,
		ReminderTime: createdAt.Add(remindInMinutes),
		ReminderText: input.ReminderText,
		ReminderName: input.ReminderName,
	}
	log.Println("Starting workflow to remind", input.Phone, "in", remindInMinutes, "minutes, at", reminderDetails.GetReminderTime().Format(app.TIME_FORMAT))
	we, err := c.ExecuteWorkflow(context.Background(), options, MakeReminderWorkflow, reminderDetails)
	if err != nil {
		log.Fatalln("error starting Reminder workflow", err)
	}
	reminderDetails.RunId = we.GetRunID()
	reminderDetails.WorkflowId = we.GetID()
	return reminderDetails, err
}

func UpdateWorkflow(workflowId string, runId string, input app.ReminderInput) (app.ReminderDetails, error) {
	c, err := client.NewClient(client.Options{})
	if err != nil {
		log.Fatalln("unable to create Temporal client", err)
	}
	defer c.Close()

	signal := app.UpdateReminderSignal{
		Phone:        input.Phone,
		NMinutes:     input.NMinutes,
		ReminderName: input.ReminderName,
		ReminderText: input.ReminderText,
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

func printResults(reminderDetails app.ReminderDetails, workflowId, runId string) {
	log.Printf(
		"\nCreating reminder for %s (%s) at %s. workflowId=%s runID=%s\n",
		reminderDetails.ReminderName,
		reminderDetails.ReminderText,
		app.GetReminderTime(reminderDetails.CreatedAt, reminderDetails.NMinutes),
		workflowId,
		runId,
	)
}
