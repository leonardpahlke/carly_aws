package internal

import (
	"github.com/pulumi/pulumi/sdk/v2/go/pulumi"
)

func CreateCrawlerHub(ctx *pulumi.Context, config CrawlerHubConfig) (CrawlerHubData, error) {
	return CrawlerHubData{}, nil
}

type CrawlerHubConfig struct{}

type CrawlerHubData struct{}
