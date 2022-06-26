package app

import (
	"testing"

	"github.com/stretchr/testify/mock"

	"github.com/stretchr/testify/require"
	"go.temporal.io/sdk/testsuite"
)

func Test_Workflow(t *testing.T) {
	testSuite := &testsuite.WorkflowTestSuite{}
	env := testSuite.NewTestWorkflowEnvironment()
	// Mock activity implementation
	testDetails := ReminderDetails{
		CreatedAt:    1656191661,
		ReminderTime: 1656278061,
		ReminderText: "Book return flights from Jakarta",
		ReminderName: "Flights",
		ReminderId:   "Test",
	}
	env.OnActivity(Snooze, mock.Anything, testDetails).Return(nil)
	env.OnActivity(Dismiss, mock.Anything, testDetails).Return(nil)
	env.ExecuteWorkflow(CreateReminder, testDetails)
	require.True(t, env.IsWorkflowCompleted())
	require.NoError(t, env.GetWorkflowError())
}
