package config

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
)

// GetQueue retrieves SQS client
func GetQueue() (*sqs.Client, error) {
	awsRegion := GetEnv("AWS_REGION", "us-east-1")
	sqsUrl := GetEnv("SQS_URL", "http://localhost:9324")
	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion(awsRegion))
	if err != nil {
		return nil, fmt.Errorf("failed to load config: %w", err)
	}
	sqsClient := sqs.NewFromConfig(cfg, func(o *sqs.Options) {
		o.EndpointResolver = sqs.EndpointResolverFunc(
			func(region string, options sqs.EndpointResolverOptions) (aws.Endpoint, error) {
				return aws.Endpoint{
					PartitionID:   "aws",
					URL:           sqsUrl,
					SigningRegion: awsRegion,
				}, nil
			},
		)
	})
	return sqsClient, nil
}
