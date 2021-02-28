package main

import (
	"carly_aws/internal"
	"github.com/pulumi/pulumi/sdk/v2/go/pulumi"
)

func main() {
	pulumi.Run(func(ctx *pulumi.Context) error {
		// NETWORK - VPC, Subnets, Security Groups
		networkData, err := internal.CreateNetwork(ctx, internal.NetworkConfig{})
		if err != nil {
			return err
		}

		// PERSISTENT - DynamoDB DB, Document DB
		persistentData, err := internal.CreatePersistent(
			ctx,
			internal.PersistentConfig{
				DdbArticleTableName:    "ddbArticleTableName",
				S3BucketArticleDomName: "bucket-article-dom-store",
				S3BucketArticleAnalyticsName: "bucket-article-analytics-store",
			},
		)
		if err != nil {
			return err
		}

		//// API - Gateway, Lambda Handler
		//_, err = handler.CreateAPI(ctx, handler.ApiConfig{})
		//if err != nil {
		//	return err
		//}

		//// STATIC WEBSITE - S3
		//_, err = handler.CreateStaticWebsite(ctx, handler.StaticWebsiteConfig{})
		//if err != nil {
		//	return err
		//}

		//// CI CD Website - CodePipeline
		//_, err = handler.CreateCiCdWebsite(ctx, handler.CiCdWebsiteConfig{})
		//if err != nil {
		//	return err
		//}

		// Spiders - Lambdas
		_, err = internal.CreateSpiders(ctx, internal.SpidersConfig{
			ArticleBucket: *persistentData.S3ArticleDomBucket,
			ArticleBucketAnalytics: *persistentData.S3ArticleAnalyticsBucket,
			ArticleTable:  *persistentData.DdbArticleTable,
			NetworkData:   networkData,
		})
		if err != nil {
			return err
		}

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
