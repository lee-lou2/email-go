package aws_test

import (
	"aws-ses-sender-go/pkg/aws"
	"context"
	"errors"
	"reflect"
	"testing"

	awsConfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/sesv2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

// MockSESClient struct provides a mock implementation of the SES client
type MockSESClient struct {
	mock.Mock
}

// SendEmail method is the mock implementation of the SendEmail method
func (m *MockSESClient) SendEmail(ctx context.Context, params *sesv2.SendEmailInput, optFns ...func(*sesv2.Options)) (*sesv2.SendEmailOutput, error) {
	args := m.Called(ctx, params, optFns)
	return args.Get(0).(*sesv2.SendEmailOutput), args.Error(1)
}

// TestGetClient_Success tests if the GetClient function successfully returns a client
func TestGetClient_Success(t *testing.T) {
	// Create context
	ctx := context.TODO()

	// Expected result
	expectedCfg, _ := awsConfig.LoadDefaultConfig(
		ctx,
		awsConfig.WithRegion("ap-northeast-2"),
		awsConfig.WithCredentialsProvider(
			credentials.NewStaticCredentialsProvider("test_key", "test_secret", ""),
		),
	)
	expectedClient := sesv2.NewFromConfig(expectedCfg)

	// Execute test
	client, err := aws.GetClient(ctx)

	// Validate results
	require.NoError(t, err, "unexpected error while getting AWS client")
	assert.Equal(t, reflect.TypeOf(expectedClient), reflect.TypeOf(client), "client type mismatch")
}

// TestSendSESEmail_Success tests if the SendSESEmail function successfully sends an email
func TestSendSESEmail_Success(t *testing.T) {
	// Set up mock client and input data
	mockClient := new(MockSESClient)
	ctx := context.TODO()
	subject := "Test Subject"
	body := "Test Body"
	receivers := []string{"test@example.com"}
	expectedMessageId := "12345"

	// Set up mock SendEmail method
	mockClient.On("SendEmail", ctx, mock.AnythingOfType("*sesv2.SendEmailInput"), mock.AnythingOfType("[]func(*sesv2.Options)")).Return(&sesv2.SendEmailOutput{
		MessageId: &expectedMessageId,
	}, nil)

	// Execute test
	messageId, err := aws.SendSESEmail(mockClient, ctx, &subject, &body, &receivers)

	// Validate results
	require.NoError(t, err, "unexpected error while sending email")
	assert.Equal(t, expectedMessageId, messageId, "messageId mismatch")
	mockClient.AssertExpectations(t)
}

// TestSendSESEmail_SendError tests if the SendSESEmail function returns an error when an email send operation fails
func TestSendSESEmail_SendError(t *testing.T) {
	// Set up mock client and input data
	mockClient := new(MockSESClient)
	ctx := context.TODO()
	subject := "Test Subject"
	body := "Test Body"
	receivers := []string{"test@example.com"}
	expectedError := errors.New("send email error")

	// Set up mock SendEmail method
	mockClient.On("SendEmail", ctx, mock.AnythingOfType("*sesv2.SendEmailInput"), mock.AnythingOfType("[]func(*sesv2.Options)")).Return((*sesv2.SendEmailOutput)(nil), expectedError)

	// Execute test
	_, err := aws.SendSESEmail(mockClient, ctx, &subject, &body, &receivers)

	// Validate results
	require.Error(t, err, "expected error but got nil")
	assert.EqualError(t, err, expectedError.Error(), "error message mismatch")
	mockClient.AssertExpectations(t)
}
