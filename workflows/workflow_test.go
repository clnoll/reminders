package workflows

import (
	"reminders/app"
	"reminders/app/activities"
	"reminders/app/whatsapp"
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"go.temporal.io/sdk/testsuite"
)

func Test_Workflow(t *testing.T) {
	testSuite := &testsuite.WorkflowTestSuite{}
	env := testSuite.NewTestWorkflowEnvironment()
	// Mock activity implementation
	testDetails := app.ReminderDetails{
		ReminderText: "Book return flights from Jakarta",
		ReminderName: "Flights",
	}
	wc := whatsapp.MockWhatsappClient{}
	env.OnActivity(activities.Create, mock.Anything, testDetails).Return(nil)
	env.OnActivity(activities.Delete, mock.Anything, testDetails).Return(nil)
	env.ExecuteWorkflow(MakeReminderWorkflow, wc, testDetails)
	require.True(t, env.IsWorkflowCompleted())
	require.NoError(t, env.GetWorkflowError())
}
