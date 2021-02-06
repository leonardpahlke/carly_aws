package internal

import (
	"github.com/pulumi/pulumi-aws/sdk/v3/go/aws/ec2"
	"github.com/pulumi/pulumi/sdk/v2/go/pulumi"
)

func CreateCrawler(ctx *pulumi.Context, _ CrawlerConfig) (CrawlerData, error) {
	// eip
	_, err := ec2.NewEip(ctx, "myeip", &ec2.EipArgs{})
	if err != nil {
		return CrawlerData{}, err
	}
	return CrawlerData{}, nil
}

type CrawlerConfig struct {}

type CrawlerData struct {}