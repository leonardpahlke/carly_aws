package internal

import (
	"carly_aws/pkg"
	"github.com/pulumi/pulumi-aws/sdk/v3/go/aws/iam"
	"github.com/pulumi/pulumi-aws/sdk/v3/go/aws/lambda"
	"github.com/pulumi/pulumi/sdk/v2/go/pulumi"
)

const LambdaSpiderMlFolderName = "spider-ml"
const LambdaSpiderParserFolderName = "spider-parser"

func CreateSpiders(ctx *pulumi.Context, _ SpidersConfig) (SpidersData, error) {
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
	lambdaSpiderMl, err := pkg.BuildLambdaFunction(role, logPolicy, LambdaSpiderMlFolderName)
	if err != nil {
		return SpidersData{}, err
	}

	// SPIDER-PARSER
	lambdaSpiderParser, err := pkg.BuildLambdaFunction(role, logPolicy, LambdaSpiderParserFolderName)
	if err != nil {
		return SpidersData{}, err
	}

	return SpidersData{
		LambdaSpiderMl: *lambdaSpiderMl,
		LambdaSpiderParser: *lambdaSpiderParser,
	}, nil
}


type SpidersConfig struct {

}

type SpidersData struct {
	LambdaSpiderMl lambda.Function
	LambdaSpiderParser lambda.Function
}