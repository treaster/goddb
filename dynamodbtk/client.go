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

	cfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithRegion("localhost"),
		config.WithEndpointResolver(aws.EndpointResolverFunc(
			func(service, region string) (aws.Endpoint, error) {
				return aws.Endpoint{
					URL:           "http://localhost:8000",
					SigningRegion: "localhost",
				}, nil // The SigningRegion key was what's was missing! D'oh.
			})),
	)

	if err != nil {
		log.Fatalf("unable to load config for local: %v", err)
	}

	client := dynamodb.NewFromConfig(cfg, func(o *dynamodb.Options) {
		o.Credentials = credentials.NewStaticCredentialsProvider(
			"fakeMyKeyId",
			"fakeSecretAccessKey",
			"")
	})

	{
		// For debugging connectivity when running locally
		out, err := client.ListTables(context.TODO(), &dynamodb.ListTablesInput{})
		if err != nil {
			log.Fatalf("error listing tables: %s", err.Error())
		}
		log.Printf("Tables: %v", out.TableNames)
	}

	return client
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
