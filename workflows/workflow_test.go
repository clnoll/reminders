package workflows

import (
	"reminders/app"
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
		ReminderId:   "Test",
	}
	env.OnActivity(app.Create, mock.Anything, testDetails).Return(nil)
	env.OnActivity(app.Update, mock.Anything, testDetails).Return(nil)
	env.OnActivity(app.Delete, mock.Anything, testDetails).Return(nil)
	env.ExecuteWorkflow(MakeReminderWorkflow, testDetails)
	require.True(t, env.IsWorkflowCompleted())
	require.NoError(t, env.GetWorkflowError())
}