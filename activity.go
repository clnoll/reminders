package app

import (
	"context"
	"fmt"
)

func Create(ctx context.Context, reminderDetails ReminderDetails) error {
	fmt.Printf(
		"\nCreating reminder %s (%s) to alert at %s. ReferenceId: %s\n",
		reminderDetails.ReminderName,
		reminderDetails.ReminderText,
		GetReminderTime(reminderDetails.CreatedAt, reminderDetails.NMinutes).Format(TIME_FORMAT),
		reminderDetails.ReminderId,
	)
	return nil
}

func Update(ctx context.Context, reminderDetails ReminderDetails) error {
	fmt.Printf(
		"\nSnoozing reminder %s (%s) until %s. ReferenceId: %s\n",
		reminderDetails.ReminderName,
		reminderDetails.ReminderText,
		GetReminderTime(reminderDetails.CreatedAt, reminderDetails.NMinutes).Format(TIME_FORMAT),
		reminderDetails.ReminderId,
	)
	return nil
}

func Delete(ctx context.Context, reminderDetails ReminderDetails) error {
	fmt.Printf(
		"\nDismissing reminder %s (%s). ReferenceId: %s\n",
		reminderDetails.ReminderName,
		reminderDetails.ReminderText,
		reminderDetails.ReminderId,
	)
	return nil
}
