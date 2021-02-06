package internal

import "github.com/pulumi/pulumi/sdk/v2/go/pulumi"

func CreateCiCdWebsite(ctx *pulumi.Context, _ CiCdWebsiteConfig) (CiCdWebsiteContext, error) {
	return CiCdWebsiteContext{}, nil
}

type CiCdWebsiteConfig struct {}

type CiCdWebsiteContext struct {}
