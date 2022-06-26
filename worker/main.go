package main

import (
	"log"

	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/worker"

	"reminders/app"
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
	w := worker.New(c, app.ReminderTaskQueue, worker.Options{})
	w.RegisterWorkflow(app.MakeReminderWorkflow)
	w.RegisterActivity(app.Create)
	w.RegisterActivity(app.Update)
	w.RegisterActivity(app.Delete)
	// Start listening to the Task Queue
	err = w.Run(worker.InterruptCh())
	if err != nil {
		log.Fatalln("unable to start Worker", err)
	}
}

// @@@SNIPEND
