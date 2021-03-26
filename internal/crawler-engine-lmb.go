package internal

import (
	"carly_aws/shared"
	"fmt"

	pkg "github.com/leonardpahlke/carly_pkg"

	"github.com/pulumi/pulumi-aws/sdk/v3/go/aws/dynamodb"
	"github.com/pulumi/pulumi-aws/sdk/v3/go/aws/iam"
	"github.com/pulumi/pulumi-aws/sdk/v3/go/aws/lambda"
	"github.com/pulumi/pulumi-aws/sdk/v3/go/aws/s3"
	"github.com/pulumi/pulumi/sdk/v2/go/pulumi"
)

const CrawlerEngineLmbMlFolderName = "crawler-lmb-ml"
const CrawlerEngineLmbParserFolderName = "crawler-lmb-parser"
const CrawlerEngineLmbTranslatorFolderName = "crawler-lmb-translator"
const CrawlerEngineLmbDownloaderFolderName = "crawler-lmb-downloader"

const DefaultLambdaTimeout = 3

func CreateCrawlerEngineLmb(ctx *pulumi.Context, config CrawlerEngineLmbConfig) (CrawlerEngineLmbData, error) {

	/*
		Create Roles for crawler engine lambdas
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
		ROLES
	*/
	// DOWNLOADER
	spiderDownloaderRole, spiderDownloaderLogPolicy := createSpiderRolesAndPolicies(
		ctx,
		CrawlerEngineLmbDownloaderFolderName,
		policyWriteLogs,
		iam.RoleInlinePolicyArray{
			shared.CreateInlinePolicyStatement(
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

	// PARSER
	spiderParserRole, spiderParserLogPolicy := createSpiderRolesAndPolicies(
		ctx,
		CrawlerEngineLmbParserFolderName,
		policyWriteLogs,
		iam.RoleInlinePolicyArray{},
	)

	// TRANSLATOR
	spiderTranslatorRole, spiderTranslatorLogPolicy := createSpiderRolesAndPolicies(
		ctx,
		CrawlerEngineLmbTranslatorFolderName,
		policyWriteLogs,
		iam.RoleInlinePolicyArray{
			shared.CreateInlinePolicyStatement(
				"article-analytics-bucket-s3-get",
				pulumi.Sprintf(`{
					"Version": "2012-10-17",
					"Statement": [{
						"Effect": "Allow",
						"Action": "translate:TranslateText",
						"Resource": "%s/*"
					}]
				}`, config.ArticleBucket.Arn),
			),
			shared.CreateInlinePolicyStatement(
				"article-analytics-bucket-s3-get",
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

	// ML
	spiderMlRole, spiderMlLogPolicy := createSpiderRolesAndPolicies(
		ctx,
		CrawlerEngineLmbMlFolderName,
		policyWriteLogs,
		iam.RoleInlinePolicyArray{
			shared.CreateInlinePolicyStatement(
				"article-analytics-bucket-s3-get",
				pulumi.Sprintf(`{
					"Version": "2012-10-17",
					"Statement": [{
						"Effect": "Allow",
						"Action": "s3:GetObject",
						"Resource": "%s/*"
					}]
				}`, config.ArticleBucket.Arn),
			),
			shared.CreateInlinePolicyStatement(
				"article-analytics-bucket-s3-put",
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

	/*
		Lambda functions
	*/

	// DOWNLOADER
	lambdaSpiderDownloader, err := buildLambdaFunction(ctx, buildLambdaConfig{
		Role:      spiderDownloaderRole,
		LogPolicy: spiderDownloaderLogPolicy,
		Env: pulumi.StringMap{
			pkg.EnvSpiderName:    pulumi.String(pkg.CARLY_ENGINE_LMB_DOWNLOADER.Name),
			pkg.EnvArticleBucket: config.ArticleBucket.Bucket,
		},
		HandlerFolder: CrawlerEngineLmbDownloaderFolderName,
		Timeout:       DefaultLambdaTimeout,
	})
	if err != nil {
		return CrawlerEngineLmbData{}, err
	}

	// PARSER
	lambdaSpiderParser, err := buildLambdaFunction(ctx, buildLambdaConfig{
		Role:      spiderParserRole,
		LogPolicy: spiderParserLogPolicy,
		Env: pulumi.StringMap{
			pkg.EnvSpiderName: pulumi.String(pkg.CARLY_ENGINE_LMB_PARSER.Name),
		},
		HandlerFolder: CrawlerEngineLmbParserFolderName,
		Timeout:       DefaultLambdaTimeout,
	})
	if err != nil {
		return CrawlerEngineLmbData{}, err
	}

	// TRANSLATOR
	lambdaSpiderTranslator, err := buildLambdaFunction(ctx, buildLambdaConfig{
		Role:      spiderTranslatorRole,
		LogPolicy: spiderTranslatorLogPolicy,
		Env: pulumi.StringMap{
			pkg.EnvSpiderName:             pulumi.String(pkg.CARLY_ENGINE_LMB_TRANSLATOR.Name),
			pkg.EnvArticleBucketAnalytics: config.ArticleBucket.Bucket,
		},
		HandlerFolder: CrawlerEngineLmbTranslatorFolderName,
		Timeout:       DefaultLambdaTimeout,
	})
	if err != nil {
		return CrawlerEngineLmbData{}, err
	}

	// ML
	lambdaSpiderMl, err := buildLambdaFunction(ctx, buildLambdaConfig{
		Role:      spiderMlRole,
		LogPolicy: spiderMlLogPolicy,
		Env: pulumi.StringMap{
			pkg.EnvSpiderName:             pulumi.String(pkg.CARLY_ENGINE_LMB_ML.Name),
			pkg.EnvArticleBucketAnalytics: config.ArticleBucket.Bucket,
		},
		HandlerFolder: CrawlerEngineLmbMlFolderName,
		Timeout:       DefaultLambdaTimeout,
	})
	if err != nil {
		return CrawlerEngineLmbData{}, err
	}

	return CrawlerEngineLmbData{
		CrawlerEngineLmbMl:         *lambdaSpiderMl,
		CrawlerEngineLmbParser:     *lambdaSpiderParser,
		CrawlerEngineLmbDownloader: *lambdaSpiderDownloader,
		CrawlerEngineLmbTranslator: *lambdaSpiderTranslator,
	}, nil
}

type CrawlerEngineLmbConfig struct {
	ArticleBucket    s3.Bucket
	LambdaCodeBucket s3.Bucket
	ArticleTable     dynamodb.Table
	NetworkData      NetworkData
}
type CrawlerEngineLmbData struct {
	CrawlerEngineLmbMl         lambda.Function
	CrawlerEngineLmbParser     lambda.Function
	CrawlerEngineLmbDownloader lambda.Function
	CrawlerEngineLmbTranslator lambda.Function
}

func getRolePolicyForLogging(ctx *pulumi.Context, spiderName string, policyDocument pulumi.String, role iam.Role) *iam.RolePolicy {
	logPolicy, err := iam.NewRolePolicy(ctx, shared.GetResourceName(fmt.Sprintf("%s-log-policy", spiderName)), &iam.RolePolicyArgs{
		Role:   role.Name,
		Policy: policyDocument,
	})
	if err != nil {
		pkg.LogError("GetRolePolicyForLogging", "could not create iam.NewRolePolicy for lambda function", err)
	}
	return logPolicy
}

func createSpiderRolesAndPolicies(ctx *pulumi.Context, spiderName string, policyToWriteLogs pulumi.String, inlinePolicyStatements iam.RoleInlinePolicyArray) (*iam.Role, *iam.RolePolicy) {
	policyAssumeLambda := shared.IamCreatePolicyString_Assume_Policy(shared.IamPolicy_Service_Lambda)
	spiderRole := shared.CreateRole(
		ctx,
		fmt.Sprintf("%s-role", spiderName),
		policyAssumeLambda,
		inlinePolicyStatements,
	)
	spiderLogPolicy := getRolePolicyForLogging(ctx, spiderName, policyToWriteLogs, *spiderRole)
	return spiderRole, spiderLogPolicy
}

// Create a lambda function
func buildLambdaFunction(ctx *pulumi.Context, config buildLambdaConfig) (*lambda.Function, error) {
	lambdaHandlerFileName := "handler"
	args := &lambda.FunctionArgs{
		Handler:     pulumi.String(lambdaHandlerFileName),
		Role:        config.Role.Arn,
		Runtime:     pulumi.String("go1.x"),
		S3Key:       pulumi.Sprintf("%s/%s/latest.zip", shared.GetResourceName(StorageS3LambdaCode), config.HandlerFolder),
		Environment: lambda.FunctionEnvironmentArgs{Variables: config.Env},
		Timeout:     pulumi.Int(config.Timeout),
		Tags:        shared.GetTags(lambdaHandlerFileName),
	}

	// Create the lambda using the args.
	lambdaFunction, err := lambda.NewFunction(
		ctx,
		shared.GetResourceName(config.HandlerFolder),
		args,
		pulumi.DependsOn([]pulumi.Resource{config.LogPolicy}),
	)
	if err != nil {
		return &lambda.Function{}, err
	}

	return lambdaFunction, nil
}

type buildLambdaConfig struct {
	Role          *iam.Role
	LogPolicy     *iam.RolePolicy
	Env           pulumi.StringMap
	HandlerFolder string
	Timeout       int
}
