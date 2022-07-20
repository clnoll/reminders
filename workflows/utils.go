package workflows

import (
	"context"
	"log"
	"reminders/app"
	"reminders/app/whatsapp"
	"time"

	enums "go.temporal.io/api/enums/v1"
	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/workflow"
	"golang.org/x/exp/slices"
)

func StartWorkflow(c client.Client, wc whatsapp.WhatsappClientDefinition, input *app.ReminderInput) (app.ReminderDetails, error) {
	options := client.StartWorkflowOptions{
		ID:        "reminder-workflow",
		TaskQueue: app.ReminderTaskQueueName,
	}
	remindInMinutes := time.Minute * time.Duration(input.NMinutes)
	reminderDetails := app.ReminderDetails{
		FromTime:     input.FromTime,
		NMinutes:     remindInMinutes,
		Phone:        input.Phone,
		ReminderTime: input.FromTime.Add(remindInMinutes),
		ReminderText: input.ReminderText,
		ReminderName: input.ReminderName,
	}
	log.Println("Starting workflow to remind", input.Phone, "in", remindInMinutes, "minutes, at", reminderDetails.GetReminderTime().Format(app.TIME_FORMAT))
	we, err := c.ExecuteWorkflow(context.Background(), options, MakeReminderWorkflow, wc, reminderDetails)
	if err != nil {
		log.Fatalln("error starting Reminder workflow", err)
	}
	workflowId, runId := we.GetID(), we.GetRunID()
	reminderDetails.RunId = runId
	reminderDetails.WorkflowId = workflowId
	referenceId, err := app.MakeReferenceId(workflowId, runId)
	if err != nil {
		return reminderDetails, err
	}
	reminderDetails.ReferenceId = referenceId
	return reminderDetails, err
}

func UpdateWorkflow(c client.Client, workflowId string, runId string, input *app.ReminderInput) (app.ReminderDetails, error) {
	signal := app.UpdateReminderSignal{
		Phone:        input.Phone,
		NMinutes:     input.NMinutes,
		ReminderName: input.ReminderName,
		ReminderText: input.ReminderText,
	}
	ctx := context.Background()
	var reminderDetails app.ReminderDetails
	err := c.SignalWorkflow(ctx, workflowId, runId, app.UpdateReminderSignalChannelName, signal)
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

var WorkflowStatusDone = []enums.WorkflowExecutionStatus{
	enums.WORKFLOW_EXECUTION_STATUS_COMPLETED,
	enums.WORKFLOW_EXECUTION_STATUS_FAILED,
	enums.WORKFLOW_EXECUTION_STATUS_CANCELED,
	enums.WORKFLOW_EXECUTION_STATUS_TIMED_OUT,
}

func workflowStatusIsDone(c client.Client, ctx context.Context, workflowId string, runId string) (enums.WorkflowExecutionStatus, bool, error) {
	status, err := c.DescribeWorkflowExecution(ctx, workflowId, runId)
	if err != nil {
		return 0, false, err
	}
	if slices.Contains(WorkflowStatusDone, status.WorkflowExecutionInfo.Status) {
		return status.WorkflowExecutionInfo.Status, true, nil
	} else {
		return status.WorkflowExecutionInfo.Status, false, nil
	}
}

func getPhone(c client.Client, ctx context.Context, workflowId string, runId string) (string, error) {
	toPhone, err := c.QueryWorkflow(ctx, workflowId, runId, "getPhone")
	log.Printf("ERR!!!!!!!!!!!!!!!!! querying - %s", err)
	if err != nil {
		return "", err
	}
	var result string
	err = toPhone.Get(&result)
	return result, err
}

func DeleteWorkflow(c client.Client, wc whatsapp.WhatsappClientDefinition, workflowId string, runId string) error {
	ctx := context.Background()

	status, done, err := workflowStatusIsDone(c, ctx, workflowId, runId)
	log.Println("Workflow status", status, "done=", done, "err=", err)
	if err != nil {
		return err
	}

	phone, err := getPhone(c, ctx, workflowId, runId)
	log.Println("Workflow phone", phone, "err=", err)
	if err != nil {
		return err
	}

	if done == true {
		log.Println("Workflow", workflowId, runId, "already complete with status", status)
		wc.SendMessage(phone, "")
		return nil
	}

	// Delete the reminder
	return c.CancelWorkflow(ctx, workflowId, runId)
}

func printResults(reminderDetails app.ReminderDetails, workflowId, runId string) {
	log.Printf(
		"\nCreating reminder for %s (%s) at %s. workflowId=%s runID=%s\n",
		reminderDetails.ReminderName,
		reminderDetails.ReminderText,
		app.GetReminderTime(reminderDetails.FromTime, reminderDetails.NMinutes),
		workflowId,
		runId,
	)
}
