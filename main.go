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

		// PERSISTENT - DynamoDB DB, Document DB
		persistentData, err := internal.CreatePersistent(
			ctx,
			internal.PersistentConfig{
				DdbArticleTableName: "ddbArticleTableName",
				S3BucketArticleDomName:  "bucket-article-dom-store",
			},
		)
		if err != nil {
			return err
		}

		//// API - Gateway, Lambda Handler
		//_, err = internal.CreateAPI(ctx, internal.ApiConfig{})
		//if err != nil {
		//	return err
		//}

		//// STATIC WEBSITE - S3
		//_, err = internal.CreateStaticWebsite(ctx, internal.StaticWebsiteConfig{})
		//if err != nil {
		//	return err
		//}

		//// CI CD Website - CodePipeline
		//_, err = internal.CreateCiCdWebsite(ctx, internal.CiCdWebsiteConfig{})
		//if err != nil {
		//	return err
		//}

		// Spiders - Lambdas
		_, err = internal.CreateSpiders(ctx, internal.SpidersConfig{
			ArticleBucket: *persistentData.S3ArticleDomBucket,
			ArticleTable: *persistentData.DdbArticleTable,
		})
		if err != nil {
			return err
		}

		//// Crawler - EC2
		//_, err = internal.CreateCrawler(ctx, internal.CrawlerConfig{
		//	CrawlerSubnet: networkData.PrivateSubnet,
		//	CrawlerVpcSecurityGroups: pulumi.StringArray{networkData.CrawlerSecurityGroup.ID()},
		//})
		//if err != nil {
		//	return err
		//}

		return nil
	})
}
