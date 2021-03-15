package internal

import (
	"carly_aws/pkg"

	"github.com/pulumi/pulumi-aws/sdk/v3/go/aws/dynamodb"
	"github.com/pulumi/pulumi-aws/sdk/v3/go/aws/s3"
	"github.com/pulumi/pulumi/sdk/v2/go/pulumi"
)

func CreatePersistent(ctx *pulumi.Context, config PersistentConfig) (PersistentData, error) {
	//docdbKeyName := pkg.GetResourceName(fmt.Sprintf("Key%s", config.Mongo.Name))
	//_, err := kms.NewKey(ctx, docdbKeyName, &kms.KeyArgs{
	//	DeletionWindowInDays: pulumi.Int(10),
	//	Description:          pulumi.String(fmt.Sprintf("%s encryption key", config.Mongo.Name)),
	//	Tags: pkg.GetTags(docdbKeyName),
	//})
	//if err != nil {
	//	return PersistentData{}, err
	//}

	// Article - DynamoDB
	ddbTableArticleRef, err := dynamodb.NewTable(ctx, pkg.GetResourceName(pkg.DdbArticleTableName), &dynamodb.TableArgs{
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
		Name:          pulumi.String(pkg.GetResourceName(pkg.DdbArticleTableName)),
		Tags:          pkg.GetTags(pkg.DdbArticleTableName),
	})
	if err != nil {
		return PersistentData{}, err
	}

	// Article Dom S3-Bucket
	s3ArticleDomBucket, err := s3.NewBucket(ctx, pkg.GetResourceName(pkg.S3BucketArticleDomName), &s3.BucketArgs{
		Bucket: pulumi.String(pkg.GetResourceName(pkg.S3BucketArticleDomName)),
		Tags:   pkg.GetTags(pkg.S3BucketArticleDomName),
	})
	if err != nil {
		return PersistentData{}, err
	}

	// Article Analytics S3-Bucket
	s3ArticleAnalyticsBucket, err := s3.NewBucket(ctx, pkg.GetResourceName(pkg.S3BucketArticleAnalyticsName), &s3.BucketArgs{
		Bucket: pulumi.String(pkg.GetResourceName(pkg.S3BucketArticleAnalyticsName)),
		Tags:   pkg.GetTags(pkg.S3BucketArticleAnalyticsName),
	})
	if err != nil {
		return PersistentData{}, err
	}

	//mongoDbAmi, err := ec2.GetAmi(ctx, pkg.GetResourceName(config.Mongo.Name), pulumi.ID(config.Mongo.AmiId), &ec2.AmiState{
	//
	//	Description: pulumi.String(fmt.Sprintf("%s Database", config.Mongo.Name)),
	//	Tags: 				   pkg.GetTags(config.Mongo.Name),
	//})
	//if err != nil {
	//	return PersistentData{}, err
	//}

	//docdbArticle, err := docdb.NewCluster(ctx, pkg.GetResourceName(config.Mongo.Name), &docdb.ClusterArgs{
	//	BackupRetentionPeriod: pulumi.Int(config.Mongo.BackupRetentionPeriod),
	//	ClusterIdentifier:     pulumi.String(pkg.GetResourceName(fmt.Sprintf("%sCluster", config.Mongo.Name))),
	//	MasterPassword:        pulumi.String("mustbeeightchars"),
	//	MasterUsername:        pulumi.String(config.Mongo.MasterUsername),
	//	KmsKeyId: 			   docdbKey.KeyId,
	//	ClusterMembers: 	   pulumi.StringArray{},
	//	PreferredBackupWindow: pulumi.String("07:00-09:00"),
	//	Port:				   pulumi.Int(config.Mongo.Port),
	//	StorageEncrypted: 	   pulumi.Bool(true),
	//	SkipFinalSnapshot:     pulumi.Bool(true),
	//	VpcSecurityGroupIds:   pulumi.StringArray{config.Mongo.VpcSecurityGroupIds},
	//	Tags: 				   pkg.GetTags(config.Mongo.Name),
	//})
	//if err != nil {
	//	return PersistentData{}, err
	//}
	//
	//var clusterInstances []*docdb.ClusterInstance
	//for i := 1; i <= config.Mongo.InstanceCount; i++ {
	//	__res, err := docdb.NewClusterInstance(ctx, fmt.Sprintf("clusterInstances-%v", i), &docdb.ClusterInstanceArgs{
	//		Identifier:        pulumi.String(fmt.Sprintf("%v%v", "docdb-cluster-demo-", i)),
	//		ClusterIdentifier: _default.ID(),
	//		InstanceClass:     pulumi.String("db.t3.medium"),
	//	})
	//	if err != nil {
	//		return PersistentData{}, err
	//	}
	//	clusterInstances = append(clusterInstances, __res)
	//}

	return PersistentData{
		DdbArticleTable:          ddbTableArticleRef,
		S3ArticleDomBucket:       s3ArticleDomBucket,
		S3ArticleAnalyticsBucket: s3ArticleAnalyticsBucket,
		// MongoDbArticleAmi:       mongoDbAmi,
	}, nil
}

type PersistentConfig struct {
	// DdbArticleTableName    string
	// S3BucketArticleDomName string
	// S3BucketArticleAnalyticsName string
	// Mongo   PersistentMongoConfig
}

//type PersistentMongoConfig struct {
//	AmiId string
//	AmiArn string
//	Name string
//	MasterUsername string
//	Port int
//}

type PersistentData struct {
	DdbArticleTable          *dynamodb.Table
	S3ArticleDomBucket       *s3.Bucket
	S3ArticleAnalyticsBucket *s3.Bucket
	// MongoDbArticleAmi *ec2.Ami
}
