package workflows

import (
	"reminders/app"
	"time"

	"go.temporal.io/sdk/temporal"
	"go.temporal.io/sdk/workflow"
)

// @@@SNIPSTART reminders-workflow
func MakeReminderWorkflow(ctx workflow.Context, reminderDetails app.ReminderDetails) error {
	// RetryPolicy specifies how to automatically handle retries if an Activity fails.
	retrypolicy := &temporal.RetryPolicy{
		InitialInterval:    time.Second,
		BackoffCoefficient: 2.0,
		MaximumInterval:    time.Minute,
		MaximumAttempts:    500,
	}
	options := workflow.ActivityOptions{
		// Timeout options specify when to automatically timeout Activity functions.
		StartToCloseTimeout: time.Minute,
		// Optionally provide a customized RetryPolicy.
		// Temporal retries failures by default, this is just an example.
		RetryPolicy: retrypolicy,
	}

	ctx = workflow.WithActivityOptions(ctx, options)

	// Create a reminder
	err := workflow.ExecuteActivity(ctx, app.Create, reminderDetails).Get(ctx, nil)
	if err != nil {
		return err
	}

	// Dummy Update
	err = workflow.ExecuteActivity(ctx, app.Update, reminderDetails).Get(ctx, nil)
	if err != nil {
		return err
	}

	// Sleep until the reminder time has elapsed
	workflow.Sleep(ctx, reminderDetails.ReminderTime.Sub(reminderDetails.CreatedAt))

	return nil
}

// @@@SNIPEND
