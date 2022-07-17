package workflows

import (
	"log"
	"reminders/app"
	"reminders/app/activities"
	"reminders/app/whatsapp"
	"time"

	"go.temporal.io/sdk/temporal"
	"go.temporal.io/sdk/workflow"
)

func MakeReminderWorkflow(ctx workflow.Context, wc whatsapp.WhatsappClientDefinition, reminderDetails app.ReminderDetails) error {
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
	err := workflow.ExecuteActivity(ctx, activities.Create, reminderDetails).Get(ctx, nil)
	if err != nil {
		return err
	}

	// Handle any incoming updates and/or wait until the reminder time has elapsed
	var reminderUpdateVal app.UpdateReminderSignal
	updateReminderChannel := workflow.GetSignalChannel(ctx, app.UpdateReminderSignalChannelName)
	timerFired := false
	for !timerFired && ctx.Err() == nil {
		timerCtx, timerCancel := workflow.WithCancel(ctx)
		timeToReminder := reminderDetails.GetMinutesToReminder(timerCtx)
		timer := workflow.NewTimer(timerCtx, timeToReminder)
		log.Println("Remind in", timeToReminder, "minutes, at", reminderDetails.GetReminderTime().Format(app.TIME_FORMAT))
		workflow.NewSelector(timerCtx).
			AddFuture(timer, func(f workflow.Future) {
				err := f.Get(timerCtx, nil)
				_ = workflow.ExecuteActivity(timerCtx, activities.SendReminder, wc.GetWhatsappClient(), reminderDetails).Get(timerCtx, nil)
				if err == nil {
					log.Println("Reminder fired")
					timerFired = true
				} else if ctx.Err() != nil {
					// if a timer returned an error then it was canceled
					log.Println("Reminder canceled")
				}
			}).
			AddReceive(updateReminderChannel, func(c workflow.ReceiveChannel, more bool) {
				timerCancel() // Create a new timer even if the reminder time hasn't been updated
				c.Receive(timerCtx, &reminderUpdateVal)
				originalNMinutes := reminderDetails.NMinutes
				updated := updateReminderDetails(timerCtx, &reminderUpdateVal, &reminderDetails)
				log.Println("ReminderDetails updated: ", reminderDetails)

				if updated.NMinutes != originalNMinutes {
					log.Println("New reminder time set:", reminderDetails.ReminderTime.Format(app.TIME_FORMAT))
				}

			}).
			Select(timerCtx)
	}
	return ctx.Err()
}
