package dynamodbtk

import (
	"context"
	"log"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
)

func NewClient(region string) *dynamodb.Client {
	if region == "localhost" {
		return localhostClient()
	} else {
		return awsClient(region)
	}
}

func localhostClient() *dynamodb.Client {
	log.Printf("Using 'localhost' mode")

	cfg, err := config.LoadDefaultConfig(
		context.TODO(),
		config.WithEndpointResolver(aws.EndpointResolverFunc(
			func(service, region string) (aws.Endpoint, error) {
				return aws.Endpoint{
					URL:           "http://localhost:8000",
					SigningRegion: "localhost",
				}, nil
			})),
		config.WithCredentialsProvider(credentials.StaticCredentialsProvider{
			Value: aws.Credentials{
				AccessKeyID:     "dummy",
				SecretAccessKey: "dummy",
				SessionToken:    "dummy",
				Source:          "credentials values are irrelevant for local DynamoDB",
			},
		}),
	)

	if err != nil {
		log.Fatalf("unable to load config for local: %v", err)
	}

	return dynamodb.NewFromConfig(cfg)
}

func awsClient(region string) *dynamodb.Client {
	cfg, err := config.LoadDefaultConfig(context.TODO(), func(o *config.LoadOptions) error {
		o.Region = region
		return nil
	})
	if err != nil {
		log.Fatalf("unable to load config for aws: %v", err.Error())
	}

	return dynamodb.NewFromConfig(cfg)
}
