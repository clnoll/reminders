package main

import (
	"log"

	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/worker"

	"reminders/app"
	"reminders/app/activities"
	"reminders/app/workflows"
)

// @@@SNIPSTART reminders-worker
func main() {
	// Create the client object just once per process
	c, err := client.NewClient(client.Options{})
	if err != nil {
		log.Fatalln("unable to create Temporal client", err)
	}
	defer c.Close()
	// This worker hosts both Workflow and Activity functions
	w := worker.New(c, app.ReminderTaskQueueName, worker.Options{})
	w.RegisterWorkflow(workflows.MakeReminderWorkflow)
	w.RegisterActivity(activities.Create)
	w.RegisterActivity(activities.Delete)
	w.RegisterActivity(activities.SendReminder)
	// Start listening to the Task Queue
	err = w.Run(worker.InterruptCh())
	if err != nil {
		log.Fatalln("unable to start Worker", err)
	}
}

// @@@SNIPEND
