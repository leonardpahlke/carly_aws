package main

import (
	"carly_aws/internal"

	"github.com/pulumi/pulumi/sdk/v2/go/pulumi"
)

func main() {
	pulumi.Run(func(ctx *pulumi.Context) error {
		// NETWORK - VPC, Subnets, Security Groups
		_, err := internal.CreateNetwork(ctx, internal.NetworkConfig{})
		if err != nil {
			return err
		}

		// STORAGE - DynamoDB DB, Document DB
		_, err = internal.CreateStorage(ctx, internal.StorageConfig{})
		if err != nil {
			return err
		}

		// CRAWLER-ENGINE-LMB
		/*
			_, err = internal.CreateCrawlerEngineLmb(ctx, internal.CrawlerEngineLmbConfig{
				ArticleBucket:    *storageData.S3ArticleDomBucket,
				LambdaCodeBucket: *storageData.S3LambdaCodeBucket,
				ArticleTable:     *storageData.DdbArticleTable,
				NetworkData:      networkData,
			})
			if err != nil {
				return err
			}
		*/

		//// Crawler - EC2
		//_, err = handler.CreateCrawler(ctx, handler.CrawlerConfig{
		//	CrawlerSubnet: networkData.PrivateSubnet,
		//	CrawlerVpcSecurityGroups: pulumi.StringArray{networkData.CrawlerSecurityGroup.ID()},
		//})
		//if err != nil {
		//	return err
		//}

		return nil
	})
}
