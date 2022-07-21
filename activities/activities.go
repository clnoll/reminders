package activities

import (
	"context"
	"fmt"

	"reminders/app"
	"reminders/app/utils"
	"reminders/app/whatsapp"
)

func Create(ctx context.Context, reminderDetails utils.ReminderDetails) error {
	fmt.Printf(
		"\nCreating reminder %s (%s) to alert at %s. workflowId=%s runId=%s\n",
		reminderDetails.ReminderName,
		reminderDetails.ReminderText,
		utils.GetReminderTime(reminderDetails.FromTime, reminderDetails.NMinutes).Format(app.TIME_FORMAT),
		reminderDetails.WorkflowId,
		reminderDetails.RunId,
	)
	return nil
}

func Update(ctx context.Context, reminderDetails utils.ReminderDetails) error {
	fmt.Printf(
		"\nSnoozing reminder %s (%s) until %s. workflowId=%s runId=%s\n",
		reminderDetails.ReminderName,
		reminderDetails.ReminderText,
		utils.GetReminderTime(reminderDetails.FromTime, reminderDetails.NMinutes).Format(app.TIME_FORMAT),
		reminderDetails.WorkflowId,
		reminderDetails.RunId,
	)
	return nil
}

func Delete(ctx context.Context, reminderDetails utils.ReminderDetails) error {
	fmt.Printf(
		"\nDismissing reminder %s (%s). workflowId=%s runId=%s\n",
		reminderDetails.ReminderName,
		reminderDetails.ReminderText,
		reminderDetails.WorkflowId,
		reminderDetails.RunId,
	)
	return nil
}

func SendReminder(ctx context.Context, wc whatsapp.IWhatsappClient, reminderDetails utils.ReminderDetails) error {
	fmt.Printf(
		"\nSending reminder to %s: %s (%s)! workflowId=%s runId=%s\n",
		reminderDetails.Phone,
		reminderDetails.ReminderName,
		reminderDetails.ReminderText,
		reminderDetails.WorkflowId,
		reminderDetails.RunId,
	)
	message := makeReminderMessage(reminderDetails)
	return wc.SendMessage(reminderDetails.Phone, message)
}

func makeReminderMessage(reminderDetails utils.ReminderDetails) string {
	return fmt.Sprintf(
		"%s: %s\nReference ID: %s",
		reminderDetails.ReminderName,
		reminderDetails.ReminderText,
		reminderDetails.ReferenceId,
	)
}
