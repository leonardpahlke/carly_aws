package internal

import (
	"carly_aws/shared"

	pkg "github.com/leonardpahlke/carly_pkg"

	"github.com/pulumi/pulumi-aws/sdk/v3/go/aws/dynamodb"
	"github.com/pulumi/pulumi-aws/sdk/v3/go/aws/s3"
	"github.com/pulumi/pulumi/sdk/v2/go/pulumi"
)

const StorageS3LambdaCode = "bucket-lambda-code"

func CreateStorage(ctx *pulumi.Context, config StorageConfig) (StorageData, error) {
	// Article - DynamoDB
	ddbTableArticleRef, err := dynamodb.NewTable(ctx, shared.GetResourceName(pkg.DdbArticleTableName), &dynamodb.TableArgs{
		Attributes: dynamodb.TableAttributeArray{
			&dynamodb.TableAttributeArgs{
				Name: pulumi.String(pkg.DdbPrimaryKeyArticleRef),
				Type: pulumi.String("S"),
			},
			&dynamodb.TableAttributeArgs{
				Name: pulumi.String(pkg.DdbSortKeyNewspaper),
				Type: pulumi.String("S"),
			},
		},
		HashKey:       pulumi.String(pkg.DdbPrimaryKeyArticleRef),
		RangeKey:      pulumi.String(pkg.DdbSortKeyNewspaper),
		ReadCapacity:  pulumi.Int(1),
		WriteCapacity: pulumi.Int(1),
		Name:          pulumi.String(shared.GetResourceName(pkg.DdbArticleTableName)),
		Tags:          shared.GetTags(pkg.DdbArticleTableName),
	})
	if err != nil {
		return StorageData{}, err
	}

	// Article Dom S3-Bucket
	s3ArticleBucket := createSimpleBucket(ctx, "article-bucket")

	// Lambda Code Bucket
	s3LambdaCodeBucket := createSimpleBucket(ctx, StorageS3LambdaCode)

	return StorageData{
		DdbArticleTable:    ddbTableArticleRef,
		S3ArticleBucket:    s3ArticleBucket,
		S3LambdaCodeBucket: s3LambdaCodeBucket,
	}, nil
}

// create an s3 bucket
func createSimpleBucket(ctx *pulumi.Context, bucketName string) *s3.Bucket {
	bucket, err := s3.NewBucket(ctx, shared.GetResourceName(bucketName), &s3.BucketArgs{
		Bucket: pulumi.String(shared.GetResourceName(bucketName)),
		Tags:   shared.GetTags(bucketName),
	})
	if err != nil {
		pkg.LogError("storage.createSimpleBucket", "could not create bucket", err)
	}
	return bucket
}

type StorageConfig struct{}

type StorageData struct {
	DdbArticleTable    *dynamodb.Table
	S3ArticleBucket    *s3.Bucket
	S3LambdaCodeBucket *s3.Bucket
}
