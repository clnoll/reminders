package activities

import (
	"context"
	"fmt"
	"reminders/app"
)

func Create(ctx context.Context, reminderDetails app.ReminderDetails) error {
	fmt.Printf(
		"\nCreating reminder %s (%s) to alert at %s. ReferenceId: %s\n",
		reminderDetails.ReminderName,
		reminderDetails.ReminderText,
		app.GetReminderTime(reminderDetails.CreatedAt, reminderDetails.NMinutes).Format(app.TIME_FORMAT),
		reminderDetails.ReminderId,
	)
	return nil
}

func Update(ctx context.Context, reminderDetails app.ReminderDetails) error {
	fmt.Printf(
		"\nSnoozing reminder %s (%s) until %s. ReferenceId: %s\n",
		reminderDetails.ReminderName,
		reminderDetails.ReminderText,
		app.GetReminderTime(reminderDetails.CreatedAt, reminderDetails.NMinutes).Format(app.TIME_FORMAT),
		reminderDetails.ReminderId,
	)
	return nil
}

func Delete(ctx context.Context, reminderDetails app.ReminderDetails) error {
	fmt.Printf(
		"\nDismissing reminder %s (%s). ReferenceId: %s\n",
		reminderDetails.ReminderName,
		reminderDetails.ReminderText,
		reminderDetails.ReminderId,
	)
	return nil
}

func SendReminder(ctx context.Context, reminderDetails app.ReminderDetails) error {
	fmt.Printf(
		"\nSending reminder to %s: %s (%s)! ReferenceId: %s\n",
		reminderDetails.Phone,
		reminderDetails.ReminderName,
		reminderDetails.ReminderText,
		reminderDetails.ReminderId,
	)
	return app.SendViaWhatsapp(reminderDetails)
}
