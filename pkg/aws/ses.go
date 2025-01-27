package aws

import (
	"context"
	"email-go/config"
	"github.com/aws/aws-sdk-go-v2/aws"
	awsConfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/sesv2"
	"github.com/aws/aws-sdk-go-v2/service/sesv2/types"
)

// GetClient creates an email client
func GetClient(ctx context.Context) (*sesv2.Client, error) {
	// Load AWS Config
	AccessKeyId := config.GetEnv("AWS_ACCESS_KEY_ID")
	SecretAccessKey := config.GetEnv("AWS_SECRET_ACCESS_KEY")
	cfg, err := awsConfig.LoadDefaultConfig(
		ctx,
		awsConfig.WithRegion("ap-northeast-2"),
		awsConfig.WithCredentialsProvider(
			credentials.NewStaticCredentialsProvider(AccessKeyId, SecretAccessKey, ""),
		),
	)
	if err != nil {
		return nil, err
	}
	return sesv2.NewFromConfig(cfg), nil
}

// SendSESEmail sends an email
func SendSESEmail(client *sesv2.Client, ctx context.Context, subject, body *string, receivers *[]string) (string, error) {
	sender := config.GetEnv("EMAIL_SENDER")
	input := &sesv2.SendEmailInput{
		FromEmailAddress: aws.String(sender),
		Destination: &types.Destination{
			ToAddresses: *receivers,
		},
		Content: &types.EmailContent{
			Simple: &types.Message{
				Subject: &types.Content{
					Data: aws.String(*subject),
				},
				Body: &types.Body{
					Html: &types.Content{
						Data: aws.String(*body),
					},
				},
			},
		},
	}
	result, err := client.SendEmail(ctx, input)
	if err != nil {
		return "", err
	}
	return *result.MessageId, nil
}
