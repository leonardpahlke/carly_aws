package internal

import (
	"carly_aws/pkg"
	"fmt"

	"github.com/pulumi/pulumi-aws/sdk/v3/go/aws/dynamodb"
	"github.com/pulumi/pulumi-aws/sdk/v3/go/aws/iam"
	"github.com/pulumi/pulumi-aws/sdk/v3/go/aws/lambda"
	"github.com/pulumi/pulumi-aws/sdk/v3/go/aws/s3"
	"github.com/pulumi/pulumi/sdk/v2/go/pulumi"
)

const LambdaSpiderMlFolderName = "spider-ml"
const LambdaSpiderParserFolderName = "spider-parser"
const LambdaSpiderTranslatorFolderName = "spider-translator"
const LambdaSpiderDownloaderFolderName = "spider-downloader"

func CreateSpiders(ctx *pulumi.Context, config SpidersConfig) (SpidersData, error) {

	/*
		Create Roles for Spiders
	*/
	// Policies
	policyWriteLogs := pulumi.String(`{
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
	}`)
	/*
		pkg.IamCreatePolicyString_Actions_Resource(
			pulumi.StringArray{
				pulumi.String("logs:CreateLogGroup"),
				pulumi.String("logs:CreateLogStream"),
				pulumi.String("logs:PutLogEvents"),
			},
			pulumi.String("arn:aws:logs:*:*:*"),
		)
	*/

	// Roles
	// Spider-Downloader
	spiderDownloaderRole, spiderDownloaderLogPolicy := createSpiderRolesAndPolicies(
		ctx,
		LambdaSpiderDownloaderFolderName,
		policyWriteLogs,
		iam.RoleInlinePolicyArray{
			pkg.CreateInlinePolicyStatement(
				"article-dom-bucket-s3-put-obj",
				pulumi.Sprintf(`{
					"Version": "2012-10-17",
					"Statement": [{
						"Effect": "Allow",
						"Action": "s3:PutObject",
						"Resource": "%s/*"
					}]
				}`, config.ArticleBucket.Arn),
			),
		},
	)

	// Spider-Parser
	spiderParserRole, spiderParserLogPolicy := createSpiderRolesAndPolicies(
		ctx,
		LambdaSpiderParserFolderName,
		policyWriteLogs,
		iam.RoleInlinePolicyArray{},
	)

	// Spider-Translator
	spiderTranslatorRole, spiderTranslatorLogPolicy := createSpiderRolesAndPolicies(
		ctx,
		LambdaSpiderTranslatorFolderName,
		policyWriteLogs,
		iam.RoleInlinePolicyArray{
			pkg.CreateInlinePolicyStatement(
				"article-analytics-bucket-s3-get",
				pulumi.Sprintf(`{
					"Version": "2012-10-17",
					"Statement": [{
						"Effect": "Allow",
						"Action": "translate:TranslateText",
						"Resource": "%s/*"
					}]
				}`, config.ArticleBucketAnalytics.Arn),
			),
			pkg.CreateInlinePolicyStatement(
				"article-analytics-bucket-s3-get",
				pulumi.Sprintf(`{
					"Version": "2012-10-17",
					"Statement": [{
						"Effect": "Allow",
						"Action": "s3:PutObject",
						"Resource": "%s/*"
					}]
				}`, config.ArticleBucketAnalytics.Arn),
			),
		},
	)

	// Spider-ML
	spiderMlRole, spiderMlLogPolicy := createSpiderRolesAndPolicies(
		ctx,
		LambdaSpiderMlFolderName,
		policyWriteLogs,
		iam.RoleInlinePolicyArray{
			pkg.CreateInlinePolicyStatement(
				"article-analytics-bucket-s3-get",
				pulumi.Sprintf(`{
					"Version": "2012-10-17",
					"Statement": [{
						"Effect": "Allow",
						"Action": "s3:GetObject",
						"Resource": "%s/*"
					}]
				}`, config.ArticleBucketAnalytics.Arn),
			),
			pkg.CreateInlinePolicyStatement(
				"article-analytics-bucket-s3-put",
				pulumi.Sprintf(`{
					"Version": "2012-10-17",
					"Statement": [{
						"Effect": "Allow",
						"Action": "s3:PutObject",
						"Resource": "%s/*"
					}]
				}`, config.ArticleBucketAnalytics.Arn),
			),
		},
	)

	/*
		Lambda functions
	*/

	// SPIDER-DOWNLOADER
	lambdaSpiderDownloader, err := pkg.BuildLambdaFunction(ctx, pkg.BuildLambdaConfig{
		Role:      spiderDownloaderRole,
		LogPolicy: spiderDownloaderLogPolicy,
		Env: pulumi.StringMap{
			pkg.EnvSpiderName:    pulumi.String(pkg.SpiderNameDownloader),
			pkg.EnvArticleBucket: config.ArticleBucket.Bucket,
		},
		HandlerFolder: LambdaSpiderDownloaderFolderName,
		Timeout:       pkg.DefaultLambdaTimeout,
	})
	if err != nil {
		return SpidersData{}, err
	}

	// SPIDER-PARSER
	lambdaSpiderParser, err := pkg.BuildLambdaFunction(ctx, pkg.BuildLambdaConfig{
		Role:      spiderParserRole,
		LogPolicy: spiderParserLogPolicy,
		Env: pulumi.StringMap{
			pkg.EnvSpiderName: pulumi.String(pkg.SpiderNameParser),
		},
		HandlerFolder: LambdaSpiderParserFolderName,
		Timeout:       pkg.DefaultLambdaTimeout,
	})
	if err != nil {
		return SpidersData{}, err
	}

	// SPIDER-TRANSLATOR
	lambdaSpiderTranslator, err := pkg.BuildLambdaFunction(ctx, pkg.BuildLambdaConfig{
		Role:      spiderTranslatorRole,
		LogPolicy: spiderTranslatorLogPolicy,
		Env: pulumi.StringMap{
			pkg.EnvSpiderName:             pulumi.String(pkg.SpiderNameTranslator),
			pkg.EnvArticleBucketAnalytics: config.ArticleBucketAnalytics.Bucket,
		},
		HandlerFolder: LambdaSpiderTranslatorFolderName,
		Timeout:       pkg.DefaultLambdaTimeout,
	})
	if err != nil {
		return SpidersData{}, err
	}

	// SPIDER-ML
	lambdaSpiderMl, err := pkg.BuildLambdaFunction(ctx, pkg.BuildLambdaConfig{
		Role:      spiderMlRole,
		LogPolicy: spiderMlLogPolicy,
		Env: pulumi.StringMap{
			pkg.EnvSpiderName:             pulumi.String(pkg.SpiderNameMl),
			pkg.EnvArticleBucketAnalytics: config.ArticleBucketAnalytics.Bucket,
		},
		HandlerFolder: LambdaSpiderMlFolderName,
		Timeout:       pkg.DefaultLambdaTimeout,
	})
	if err != nil {
		return SpidersData{}, err
	}

	return SpidersData{
		LambdaSpiderMl:         *lambdaSpiderMl,
		LambdaSpiderParser:     *lambdaSpiderParser,
		LambdaSpiderDownloader: *lambdaSpiderDownloader,
		LambdaSpiderTranslator: *lambdaSpiderTranslator,
	}, nil
}

type SpidersConfig struct {
	ArticleBucket          s3.Bucket
	ArticleBucketAnalytics s3.Bucket
	ArticleTable           dynamodb.Table
	NetworkData            NetworkData
}

type SpidersData struct {
	LambdaSpiderMl         lambda.Function
	LambdaSpiderParser     lambda.Function
	LambdaSpiderDownloader lambda.Function
	LambdaSpiderTranslator lambda.Function
}

func getRolePolicyForLogging(ctx *pulumi.Context, spiderName string, policyDocument pulumi.String, role iam.Role) *iam.RolePolicy {
	logPolicy, err := iam.NewRolePolicy(ctx, pkg.GetResourceName(fmt.Sprintf("%s-log-policy", spiderName)), &iam.RolePolicyArgs{
		Role:   role.Name,
		Policy: policyDocument,
	})
	if err != nil {
		pkg.LogError("GetRolePolicyForLogging", "could not create iam.NewRolePolicy for lambda function", err)
	}
	return logPolicy
}

func createSpiderRolesAndPolicies(ctx *pulumi.Context, spiderName string, policyToWriteLogs pulumi.String, inlinePolicyStatements iam.RoleInlinePolicyArray) (*iam.Role, *iam.RolePolicy) {
	policyAssumeLambda := pkg.IamCreatePolicyString_Assume_Policy(pkg.IamPolicy_Service_Lambda)
	spiderRole := pkg.CreateRole(
		ctx,
		fmt.Sprintf("%s-role", spiderName),
		policyAssumeLambda,
		inlinePolicyStatements,
	)
	spiderLogPolicy := getRolePolicyForLogging(ctx, spiderName, policyToWriteLogs, *spiderRole)
	return spiderRole, spiderLogPolicy
}
