package app

import (
	"context"
	"fmt"
)

func Create(ctx context.Context, reminderDetails ReminderDetails) error {
	fmt.Printf(
		"\nCreating reminder %s (%s) to alert at %f. ReferenceId: %s\n",
		reminderDetails.ReminderName,
		reminderDetails.ReminderText,
		reminderDetails.ReminderTime,
		reminderDetails.ReminderId,
	)
	return nil
}

func Update(ctx context.Context, reminderDetails ReminderDetails) error {
	fmt.Printf(
		"\nSnoozing reminder %s (%s) until %f. ReferenceId: %s\n",
		reminderDetails.ReminderName,
		reminderDetails.ReminderText,
		reminderDetails.ReminderTime,
		reminderDetails.ReminderId,
	)
	return nil
}

// @@@SNIPSTART reminders-activity
func Delete(ctx context.Context, reminderDetails ReminderDetails) error {
	fmt.Printf(
		"\nDismissing reminder %s (%s). ReferenceId: %s\n",
		reminderDetails.ReminderName,
		reminderDetails.ReminderText,
		reminderDetails.ReminderId,
	)
	// Switch out comments on the return statements to simulate an error
	//return fmt.Errorf("deposit did not occur due to an issue")
	return nil
}

// @@@SNIPEND"
