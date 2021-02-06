package internal

import (
	"carly_aws/pkg"
	"fmt"
	"github.com/pulumi/pulumi-aws/sdk/v3/go/aws/dynamodb"
	"github.com/pulumi/pulumi-aws/sdk/v3/go/aws/ec2"
	"github.com/pulumi/pulumi-aws/sdk/v3/go/aws/kms"
	"github.com/pulumi/pulumi/sdk/v2/go/pulumi"
)

func CreatePersistent(ctx *pulumi.Context, config PersistentConfig) (PersistentData, error) {
	docdbKeyName := pkg.GetResourceName(fmt.Sprintf("Key%s", config.Mongo.Name))
	_, err := kms.NewKey(ctx, docdbKeyName, &kms.KeyArgs{
		DeletionWindowInDays: pulumi.Int(10),
		Description:          pulumi.String(fmt.Sprintf("%s encryption key", config.Mongo.Name)),
		Tags: pkg.GetTags(docdbKeyName),
	})
	if err != nil {
		return PersistentData{}, err
	}

	// DynamoDB
	ddbTableArticleRef, err := dynamodb.NewTable(ctx, pkg.GetResourceName(config.DdbName), &dynamodb.TableArgs{
		Attributes: dynamodb.TableAttributeArray{
			&dynamodb.TableAttributeArgs{
				Name: pulumi.String("Id"),
				Type: pulumi.String("S"),
			},
		},
		HashKey:       pulumi.String("Id"),
		ReadCapacity:  pulumi.Int(1),
		WriteCapacity: pulumi.Int(1),
		Tags: pkg.GetTags(config.DdbName),
	})
	if err != nil {
		return PersistentData{}, err
	}

	mongoDbAmi, err := ec2.GetAmi(ctx, pkg.GetResourceName(config.Mongo.Name), pulumi.ID(config.Mongo.AmiId), &ec2.AmiState{

		Description: pulumi.String(fmt.Sprintf("%s Database", config.Mongo.Name)),
		Tags: 				   pkg.GetTags(config.Mongo.Name),
	})
	if err != nil {
		return PersistentData{}, err
	}


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
		DynamoDbArticleRefTable: ddbTableArticleRef,
		MongoDbArticleAmi:       mongoDbAmi,
	}, nil
}


type PersistentConfig struct {
	DdbName string
	Mongo   PersistentMongoConfig
}
type PersistentMongoConfig struct {
	AmiId string
	AmiArn string
	Name string
	MasterUsername string
	Port int
}

type PersistentData struct {
	DynamoDbArticleRefTable *dynamodb.Table
	MongoDbArticleAmi       *ec2.Ami
}