package workflows

import (
	"log"
	"reminders/app"
	"time"

	"go.temporal.io/sdk/temporal"
	"go.temporal.io/sdk/workflow"
)

type UpdateReminderSignal struct {
	NMinutes     int
	ReminderText string
	ReminderName string
	Phone        string
}

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

	// Handle any incoming updates and/or wait until the reminder time has elapsed
	childCtx, cancelHandler := workflow.WithCancel(ctx)
	selector := workflow.NewSelector(ctx)
	processed := false

	// Create a timer whose handler will send the reminder at the specified time
	timerFuture := workflow.NewTimer(childCtx, reminderDetails.NMinutes)
	log.Println("Created timer for", reminderDetails.NMinutes, "minutes")
	selector.AddFuture(timerFuture, func(f workflow.Future) {
		_ = workflow.ExecuteActivity(ctx, app.SendReminder).Get(ctx, nil)
		processed = true
	})
	// Watch for signals to update the reminder
	for {
		if processed {
			return nil
		}
		var reminderUpdateVal UpdateReminderSignal
		channel := workflow.GetSignalChannel(ctx, app.UpdateReminderSignal)
		selector.AddReceive(channel, func(c workflow.ReceiveChannel, more bool) {
			log.Println("Received signal on channel", app.UpdateReminderSignal)
			c.Receive(ctx, &reminderUpdateVal)
			updated := updateReminderDetails(reminderUpdateVal, reminderDetails)
			log.Println("ReminderDetails updated: ", updated)

			// If the reminder time was updated, cancel the existing timer and create a new one
			if updated.NMinutes != reminderDetails.NMinutes {
				cancelHandler()
				log.Println("Cancelled existing timer.")
				childCtx, cancelHandler = workflow.WithCancel(ctx)
				timerFuture := workflow.NewTimer(childCtx, reminderDetails.NMinutes)
				selector.AddFuture(timerFuture, func(f workflow.Future) {
					_ = workflow.ExecuteActivity(ctx, app.SendReminder).Get(ctx, nil)
				})
				log.Println("New Timer created.")
			}
		})
		// Wait the timer or the update
		selector.Select(ctx)
	}
}
