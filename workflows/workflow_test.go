package workflows

import (
	"reminders/app/activities"
	"reminders/app/utils"
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"go.temporal.io/sdk/testsuite"
)

func Test_Workflow(t *testing.T) {
	testSuite := &testsuite.WorkflowTestSuite{}
	env := testSuite.NewTestWorkflowEnvironment()
	// Mock activity implementation
	testDetails := utils.ReminderDetails{
		ReminderText: "Book return flights from Jakarta",
		ReminderName: "Flights",
	}
	wc := utils.MockWhatsappClient{}
	env.OnActivity(activities.Create, mock.Anything, testDetails).Return(nil)
	env.OnActivity(activities.Delete, mock.Anything, testDetails).Return(nil)
	env.ExecuteWorkflow(MakeReminderWorkflow, wc, testDetails)
	require.True(t, env.IsWorkflowCompleted())
	require.NoError(t, env.GetWorkflowError())
}
