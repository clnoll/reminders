package utils

import (
	"context"
	"errors"
	"reminders/app/whatsapp"

	"go.temporal.io/api/workflowservice/v1"
	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/converter"
)

// MOCKING UTILITIES

// Workflow Client

type MockWorkflowClient struct {
	client.Client
}

func (f MockWorkflowClient) GetClient() (client.Client, error) {
	return MockWorkflowClient{}, nil
}

type MockEncodedValue struct {
	value string
}

func (b MockEncodedValue) Get(valuePtr interface{}) error {
	if !b.HasValue() {
		return errors.New("No value!")
	}
	valuePtr = b.value
	return nil
}

func (b MockEncodedValue) HasValue() bool {
	return true
}

func (f MockWorkflowClient) QueryWorkflow(ctx context.Context, workflowID string, runID string, queryType string, args ...interface{}) (converter.EncodedValue, error) {
	if queryType == "getPhone" {
		ev := MockEncodedValue{"12345"}
		return ev, nil
	}
	return nil, nil
}

func (f MockWorkflowClient) DescribeWorkflowExecution(ctx context.Context, workflowID, runID string) (*workflowservice.DescribeWorkflowExecutionResponse, error) {
	c, _ := client.NewClient(client.Options{})
	return c.DescribeWorkflowExecution(ctx, workflowID, runID)
}

func (f MockWorkflowClient) CancelWorkflow(ctx context.Context, workflowID string, runID string) error {
	c, _ := client.NewClient(client.Options{})
	return c.CancelWorkflow(ctx, workflowID, runID)
}

func (f MockWorkflowClient) Close() {
	c, _ := client.NewClient(client.Options{})
	c.Close()
}

// Whatsapp Client

type MockWhatsappClient struct {
	whatsapp.IWhatsappClient
}

func (f MockWhatsappClient) GetWhatsappClient() whatsapp.IWhatsappClient {
	return MockWhatsappClient{}
}

func (f MockWhatsappClient) SendMessage(toPhone string, message string) error {
	return nil
}
