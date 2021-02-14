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

	// spiderDownloaderRole
	spiderDownloaderRole, err := iam.NewRole(ctx, pkg.GetResourceName(LambdaSpiderDownloaderFolderName + "-role"), &iam.RoleArgs{
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

	spiderParserRole, err := iam.NewRole(ctx, pkg.GetResourceName("spider-parser-role"), &iam.RoleArgs{
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

	spiderMlRole, err := iam.NewRole(ctx, pkg.GetResourceName(LambdaSpiderMlFolderName + "-role"), &iam.RoleArgs{
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

	_, err = iam.NewRolePolicy(ctx, pkg.GetResourceName("s3-log-policy"), &iam.RolePolicyArgs{
		Role: spiderDownloaderRole.Name,
		Policy: pulumi.String(`{
                "Version": "2012-10-17",
                "Statement": [{
                    "Effect": "Allow",
                    "Action": "s3:PutObject",
                    "Resource": "arn:aws:s3:::carly-dev-bucket-article-dom-store/*"
                }]
            }`),
	})


	if err != nil {
		return SpidersData{}, err
	}

	// Attach a policy to allow writing logs to CloudWatch
	logPolicySpiderMl, err := iam.NewRolePolicy(ctx, pkg.GetResourceName(LambdaSpiderMlFolderName + "-log-policy"), &iam.RolePolicyArgs{
		Role: spiderMlRole.Name,
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
	if err != nil {
		return SpidersData{}, err
	}

	// Attach a policy to allow writing logs to CloudWatch
	logPolicySpiderParser, err := iam.NewRolePolicy(ctx, pkg.GetResourceName("spider-parser-log-policy"), &iam.RolePolicyArgs{
		Role: spiderParserRole.Name,
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
	if err != nil {
		return SpidersData{}, err
	}

	// Attach a policy to spider-downloader to allow to store objects in article s3 bucket
	logPolicySpiderDownloader, err := iam.NewRolePolicy(ctx, pkg.GetResourceName(LambdaSpiderDownloaderFolderName + "-log-policy"), &iam.RolePolicyArgs{
		Role: spiderDownloaderRole.Name,
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
	if err != nil {
		return SpidersData{}, err
	}

	// SPIDER-ML
	lambdaSpiderMl, err := pkg.BuildLambdaFunction(ctx, pkg.BuildLambdaConfig{
		Role:      spiderMlRole,
		LogPolicy: logPolicySpiderMl,
		Env: pulumi.StringMap{
			pkg.EnvSpiderName: pulumi.String("SpiderMl"),
		},
		HandlerFolder: LambdaSpiderMlFolderName,
		Timeout:       pkg.DefaultLambdaTimeout,
		//VpcId:           config.NetworkData.Vpc.ID(),
		//SecurityGroupId: config.NetworkData.CrawlerSecurityGroup.ID(),
		//SubnetId:        config.NetworkData.PublicSubnet.ID(),
	})
	if err != nil {
		return SpidersData{}, err
	}

	// SPIDER-TAZ-PARSER
	lambdaSpiderTazParser, err := pkg.BuildLambdaFunction(ctx, pkg.BuildLambdaConfig{
		Role:      spiderParserRole,
		LogPolicy: logPolicySpiderParser,
		Env: pulumi.StringMap{
			pkg.EnvSpiderName: pulumi.String("SpiderTazParser"),
		},
		HandlerFolder: LambdaSpiderTazParserFolderName,
		Timeout:       pkg.DefaultLambdaTimeout,
		//VpcId:           config.NetworkData.Vpc.ID(),
		//SecurityGroupId: config.NetworkData.CrawlerSecurityGroup.ID(),
		//SubnetId:        config.NetworkData.PublicSubnet.ID(),
	})
	if err != nil {
		return SpidersData{}, err
	}

	// SPIDER-DOWNLOADER
	lambdaSpiderDownloader, err := pkg.BuildLambdaFunction(ctx, pkg.BuildLambdaConfig{
		Role:      spiderDownloaderRole,
		LogPolicy: logPolicySpiderDownloader,
		Env: pulumi.StringMap{
			pkg.EnvSpiderName:    pulumi.String("SpiderDownloader"),
			pkg.EnvArticleBucket: config.ArticleBucket.Bucket,
		},
		HandlerFolder: LambdaSpiderDownloaderFolderName,
		Timeout:       pkg.DefaultLambdaTimeout,
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
