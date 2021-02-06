package internal

import (
	"carly_aws/pkg"
	"github.com/pulumi/pulumi-aws/sdk/v3/go/aws/s3"
	"github.com/pulumi/pulumi/sdk/v2/go/pulumi"
)

// https://github.com/pulumi/examples/blob/master/aws-go-s3-folder/main.go

func CreateStaticWebsite(ctx *pulumi.Context, _ StaticWebsiteConfig) (StaticWebsiteData, error) {
	// Create a bucket and expose a website index document
	siteBucket, err := s3.NewBucket(ctx, pkg.GetResourceName("S3WebsiteBucket"), &s3.BucketArgs{
		Website: s3.BucketWebsiteArgs{
			IndexDocument: pulumi.String("index.html"),
		},
	})
	if err != nil {
		return StaticWebsiteData{}, err
	}

	// Set the access policy for the bucket so all objects are readable.
	if _, err := s3.NewBucketPolicy(ctx, pkg.GetResourceName("BucketPolicy"), &s3.BucketPolicyArgs{
		Bucket: siteBucket.ID(), // refer to the bucket created earlier
		Policy: pulumi.Any(map[string]interface{}{
			"Version": "2012-10-17",
			"Statement": []map[string]interface{}{
				{
					"Effect":    "Allow",
					"Principal": "*",
					"Action": []interface{}{
						"s3:GetObject",
					},
					"Resource": []interface{}{
						pulumi.Sprintf("arn:aws:s3:::%s/*", siteBucket.ID()), // policy refers to bucket name explicitly
					},
				},
			},
		}),
	}); err != nil {
		return StaticWebsiteData{},err
	}

	ctx.Export("bucketName", siteBucket.ID())
	ctx.Export("websiteUrl", siteBucket.WebsiteEndpoint)

	return StaticWebsiteData{
		SiteBucket: siteBucket,
	}, nil
}

type StaticWebsiteConfig struct {}

type StaticWebsiteData struct {
	SiteBucket *s3.Bucket
}