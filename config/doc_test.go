package config_test

import (
	"context"
	"fmt"
	"log"

	"github.com/mvdatacenter/mvdata-sdk-go/aws"
	"github.com/mvdatacenter/mvdata-sdk-go/config"
	"github.com/mvdatacenter/mvdata-sdk-go/service/sts"
)

func Example() {
	ctx := context.TODO()
	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		log.Fatal(err)
	}

	client := sts.NewFromConfig(cfg)
	identity, err := client.GetCallerIdentity(ctx, &sts.GetCallerIdentityInput{})
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Account: %s, Arn: %s", aws.ToString(identity.Account), aws.ToString(identity.Arn))
}

func Example_custom_config() {
	ctx := context.TODO()

	// Config sources can be passed to LoadDefaultConfig, these sources can implement one or more
	// provider interfaces. These sources take priority over the standard environment and shared configuration values.
	cfg, err := config.LoadDefaultConfig(ctx,
		config.WithRegion("us-west-2"),
		config.WithSharedConfigProfile("customProfile"),
	)
	if err != nil {
		log.Fatal(err)
	}

	client := sts.NewFromConfig(cfg)
	identity, err := client.GetCallerIdentity(ctx, &sts.GetCallerIdentityInput{})
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Account: %s, Arn: %s", aws.ToString(identity.Account), aws.ToString(identity.Arn))
}
