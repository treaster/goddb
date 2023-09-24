package ecstk

import (
	"context"
	"log"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ecs"
)

func NewClient(ctx context.Context, region string) *ecs.Client {
	cfg, err := config.LoadDefaultConfig(
		context.TODO(),
		func(o *config.LoadOptions) error {
			o.Region = region
			return nil
		})
	if err != nil {
		log.Fatalf("unable to load config for aws: %v", err.Error())
	}

	return ecs.NewFromConfig(cfg)
}
