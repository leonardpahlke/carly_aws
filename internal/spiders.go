package internal

import (
	"carly_aws/pkg"
	"github.com/pulumi/pulumi-aws/sdk/v3/go/aws/dynamodb"
	"github.com/pulumi/pulumi-aws/sdk/v3/go/aws/iam"
	"github.com/pulumi/pulumi-aws/sdk/v3/go/aws/lambda"
	"github.com/pulumi/pulumi-aws/sdk/v3/go/aws/s3"
	"github.com/pulumi/pulumi/sdk/v2/go/pulumi"
)

const LambdaSpiderMlFolderName = "spider-ml"
const LambdaSpiderTazParserFolderName = "spider-taz-parser"
const LambdaSpiderDownloaderFolderName = "spider-downloader"

func CreateSpiders(ctx *pulumi.Context, config SpidersConfig) (SpidersData, error) {
	// Create an IAM role.
	role, err := iam.NewRole(ctx, "task-exec-role", &iam.RoleArgs{
		AssumeRolePolicy: pulumi.String(`{
				"Version": "2012-10-17",
				"Statement": [{
					"Sid": "",
					"Effect": "Allow",
					"Principal": {
						"Service": "lambda.amazonaws.com"
					},
					"Action": "sts:AssumeRole"
				}]
			}`),
	})
	if err != nil {
		return SpidersData{}, err
	}

	// Attach a policy to allow writing logs to CloudWatch
	logPolicy, err := iam.NewRolePolicy(ctx, pkg.GetResourceName("lambda-log-policy"), &iam.RolePolicyArgs{
		Role: role.Name,
		Policy: pulumi.String(`{
                "Version": "2012-10-17",
                "Statement": [{
                    "Effect": "Allow",
                    "Action": [
                        "logs:CreateLogGroup",
                        "logs:CreateLogStream",
                        "logs:PutLogEvents"
                    ],
                    "Resource": "arn:aws:logs:*:*:*"
                }]
            }`),
	})

	// SPIDER-ML
	lambdaSpiderMl, err := pkg.BuildLambdaFunction(ctx, pkg.BuildLambdaConfig{
		Role:            role,
		LogPolicy:       logPolicy,
		Env:             pulumi.StringMap{
			pkg.EnvSpiderName: pulumi.String("SpiderMl"),
		},
		HandlerFolder:   LambdaSpiderMlFolderName,
		Timeout:         pkg.DefaultLambdaTimeout,
		//VpcId:           config.NetworkData.Vpc.ID(),
		//SecurityGroupId: config.NetworkData.CrawlerSecurityGroup.ID(),
		//SubnetId:        config.NetworkData.PublicSubnet.ID(),
	})
	if err != nil {
		return SpidersData{}, err
	}

	// SPIDER-TAZ-PARSER
	lambdaSpiderTazParser, err := pkg.BuildLambdaFunction(ctx, pkg.BuildLambdaConfig{
		Role:            role,
		LogPolicy:       logPolicy,
		Env:             pulumi.StringMap{
			pkg.EnvSpiderName: pulumi.String("SpiderTazParser"),
		},
		HandlerFolder:   LambdaSpiderTazParserFolderName,
		Timeout:         pkg.DefaultLambdaTimeout,
		//VpcId:           config.NetworkData.Vpc.ID(),
		//SecurityGroupId: config.NetworkData.CrawlerSecurityGroup.ID(),
		//SubnetId:        config.NetworkData.PublicSubnet.ID(),
	})
	if err != nil {
		return SpidersData{}, err
	}

	// SPIDER-DOWNLOADER
	lambdaSpiderDownloader, err := pkg.BuildLambdaFunction(ctx, pkg.BuildLambdaConfig{
		Role:            role,
		LogPolicy:       logPolicy,
		Env: 			 pulumi.StringMap{
			pkg.EnvSpiderName: pulumi.String("SpiderDownloader"),
			pkg.EnvArticleBucket: config.ArticleBucket.Bucket,
		},
		HandlerFolder:   LambdaSpiderDownloaderFolderName,
		Timeout:         pkg.DefaultLambdaTimeout,
		//VpcId:           config.NetworkData.Vpc.ID(),
		//SecurityGroupId: config.NetworkData.CrawlerSecurityGroup.ID(),
		//SubnetId:        config.NetworkData.PublicSubnet.ID(),
	})
	if err != nil {
		return SpidersData{}, err
	}

	return SpidersData{
		LambdaSpiderMl:         *lambdaSpiderMl,
		LambdaSpiderTazParser:  *lambdaSpiderTazParser,
		LambdaSpiderDownloader: *lambdaSpiderDownloader,
	}, nil
}


type SpidersConfig struct {
	ArticleBucket s3.Bucket
	ArticleTable  dynamodb.Table
	NetworkData   NetworkData
}

type SpidersData struct {
	LambdaSpiderMl         lambda.Function
	LambdaSpiderTazParser  lambda.Function
	LambdaSpiderDownloader lambda.Function
}