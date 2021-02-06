package internal

import (
	"carly_aws/pkg"
	"github.com/pulumi/pulumi-aws/sdk/v3/go/aws/apigateway"
	"github.com/pulumi/pulumi-aws/sdk/v3/go/aws/iam"
	"github.com/pulumi/pulumi-aws/sdk/v3/go/aws/lambda"
	"github.com/pulumi/pulumi/sdk/v2/go/pulumi"
)

func CreateAPI(ctx *pulumi.Context, _ ApiConfig) (ApiContext, error) {
	account, region, err := pkg.GetAwsMeta(ctx)
	if err != nil {
		return ApiContext{}, err
	}

	// Create an IAM role.
	role, err := iam.NewRole(ctx, pkg.GetResourceName("task-exec-role"), &iam.RoleArgs{
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
		return ApiContext{}, err
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
	lambdaFunction, err := pkg.BuildLambdaFunction(ctx, role, logPolicy, "handler")
	if err != nil {
		return ApiContext{}, err
	}

	// Create a new API Gateway.
	gateway, err := apigateway.NewRestApi(ctx, pkg.GetResourceName("UpperCaseGateway"), &apigateway.RestApiArgs{
		Name:        pulumi.String("UpperCaseGateway"),
		Description: pulumi.String("An API Gateway for the UpperCase function"),
		Policy: pulumi.String(`{
			"Version": "2012-10-17",
  			"Statement": [{
      			"Action": "sts:AssumeRole",
      			"Principal": {
        			"Service": "lambda.amazonaws.com"
      			},
      			"Effect": "Allow",
				"Sid": ""
			},
			{
			  "Action": "execute-api:Invoke",
			  "Resource": "*",
			  "Principal": "*",
			  "Effect": "Allow",
			  "Sid": ""
			}
		  ]
		}`)})
	if err != nil {
		return ApiContext{}, err
	}

	// Add a resource to the API Gateway.
	// This makes the API Gateway accept requests on "/{message}".
	apiResource, err := apigateway.NewResource(ctx, pkg.GetResourceName("UpperAPI"), &apigateway.ResourceArgs{
		RestApi:  gateway.ID(),
		PathPart: pulumi.String("{proxy+}"),
		ParentId: gateway.RootResourceId,
	}, pulumi.DependsOn([]pulumi.Resource{gateway}))
	if err != nil {
		return ApiContext{}, err
	}

	// Add a method to the API Gateway.
	_, err = apigateway.NewMethod(ctx, pkg.GetResourceName("AnyMethod"), &apigateway.MethodArgs{
		HttpMethod:    pulumi.String("ANY"),
		Authorization: pulumi.String("NONE"),
		RestApi:       gateway.ID(),
		ResourceId:    apiResource.ID(),
	}, pulumi.DependsOn([]pulumi.Resource{gateway, apiResource}))
	if err != nil {
		return ApiContext{}, err
	}

	// Add an integration to the API Gateway.
	// This makes communication between the API Gateway and the Lambda function work
	_, err = apigateway.NewIntegration(ctx, pkg.GetResourceName("LambdaIntegration"), &apigateway.IntegrationArgs{
		HttpMethod:            pulumi.String("ANY"),
		IntegrationHttpMethod: pulumi.String("POST"),
		ResourceId:            apiResource.ID(),
		RestApi:               gateway.ID(),
		Type:                  pulumi.String("AWS_PROXY"),
		Uri:                   lambdaFunction.InvokeArn,
	}, pulumi.DependsOn([]pulumi.Resource{gateway, apiResource, lambdaFunction}))
	if err != nil {
		return ApiContext{}, err
	}

	// Add a resource based policy to the Lambda function.
	// This is the final step and allows AWS API Gateway to communicate with the AWS Lambda function
	permission, err := lambda.NewPermission(ctx, pkg.GetResourceName("APIPermission"), &lambda.PermissionArgs{
		Action:    pulumi.String("lambda:InvokeFunction"),
		Function:  lambdaFunction.Name,
		Principal: pulumi.String("apigateway.amazonaws.com"),
		SourceArn: pulumi.Sprintf("arn:aws:execute-api:%s:%s:%s/*/*/*", region.Name, account.AccountId, gateway.ID()),
	}, pulumi.DependsOn([]pulumi.Resource{gateway, apiResource, lambdaFunction}))
	if err != nil {
		return ApiContext{}, err
	}

	// Create a new deployment
	_, err = apigateway.NewDeployment(ctx, pkg.GetResourceName("APIDeployment"), &apigateway.DeploymentArgs{
		Description:      pulumi.String("UpperCase API deployment"),
		RestApi:          gateway.ID(),
		StageDescription: pulumi.String("Production"),
		StageName:        pulumi.String("prod"),
	}, pulumi.DependsOn([]pulumi.Resource{gateway, apiResource, lambdaFunction, permission}))
	if err != nil {
		return ApiContext{}, err
	}

	ctx.Export("invocation URL", pulumi.Sprintf("https://%s.execute-api.%s.amazonaws.com/prod/{message}", gateway.ID(), region.Name))

	return ApiContext{}, nil
}


type ApiConfig struct {}

type ApiContext struct {}