package activities

import (
	"context"
	"fmt"
	"reminders/app"
	"reminders/app/whatsapp"
)

func Create(ctx context.Context, reminderDetails app.ReminderDetails) error {
	fmt.Printf(
		"\nCreating reminder %s (%s) to alert at %s. workflowId=%s runId=%s\n",
		reminderDetails.ReminderName,
		reminderDetails.ReminderText,
		app.GetReminderTime(reminderDetails.FromTime, reminderDetails.NMinutes).Format(app.TIME_FORMAT),
		reminderDetails.WorkflowId,
		reminderDetails.RunId,
	)
	return nil
}

func Update(ctx context.Context, reminderDetails app.ReminderDetails) error {
	fmt.Printf(
		"\nSnoozing reminder %s (%s) until %s. workflowId=%s runId=%s\n",
		reminderDetails.ReminderName,
		reminderDetails.ReminderText,
		app.GetReminderTime(reminderDetails.FromTime, reminderDetails.NMinutes).Format(app.TIME_FORMAT),
		reminderDetails.WorkflowId,
		reminderDetails.RunId,
	)
	return nil
}

func Delete(ctx context.Context, reminderDetails app.ReminderDetails) error {
	fmt.Printf(
		"\nDismissing reminder %s (%s). workflowId=%s runId=%s\n",
		reminderDetails.ReminderName,
		reminderDetails.ReminderText,
		reminderDetails.WorkflowId,
		reminderDetails.RunId,
	)
	return nil
}

func SendReminder(ctx context.Context, reminderDetails app.ReminderDetails) error {
	fmt.Printf(
		"\nSending reminder to %s: %s (%s)! workflowId=%s runId=%s\n",
		reminderDetails.Phone,
		reminderDetails.ReminderName,
		reminderDetails.ReminderText,
		reminderDetails.WorkflowId,
		reminderDetails.RunId,
	)
	message := makeReminderMessage(reminderDetails)
	return whatsapp.SendMessage(reminderDetails.Phone, message)
}

func makeReminderMessage(reminderDetails app.ReminderDetails) string {
	return fmt.Sprintf(
		"%s: %s\nReference ID: %s_%s",
		reminderDetails.ReminderName,
		reminderDetails.ReminderText,
		reminderDetails.WorkflowId,
		reminderDetails.RunId,
	)
}
