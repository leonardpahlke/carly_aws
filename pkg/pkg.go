package pkg

import (
	"fmt"
	"github.com/pulumi/pulumi-aws/sdk/v3/go/aws"
	"github.com/pulumi/pulumi-aws/sdk/v3/go/aws/iam"
	"github.com/pulumi/pulumi-aws/sdk/v3/go/aws/lambda"
	"github.com/pulumi/pulumi/sdk/v2/go/pulumi"
)

const ProjectName = "carly"
const DeploymentEnv = "dev"

func GetResourceName(name string) string {
	return fmt.Sprintf("%s-%s-%s", ProjectName, DeploymentEnv, name)
}

func GetTags(resourceName string) pulumi.StringMap {
	return pulumi.StringMap{
		"STAGE":    pulumi.String(DeploymentEnv) ,
		"RESOURCE": pulumi.String(resourceName),
		"CREATED_BY": pulumi.String("Pulumi"),
		"PROJECT": pulumi.String(ProjectName),
	}
}

func GetAwsMeta(ctx *pulumi.Context) (*aws.GetCallerIdentityResult, *aws.GetRegionResult, error) {
	account, err := aws.GetCallerIdentity(ctx)
	if err != nil {
		return &aws.GetCallerIdentityResult{}, &aws.GetRegionResult{}, err
	}

	region, err := aws.GetRegion(ctx, &aws.GetRegionArgs{})
	if err != nil {
		return &aws.GetCallerIdentityResult{}, &aws.GetRegionResult{}, err
	}

	return account, region, nil
}

// Create a lambda function
func BuildLambdaFunction(ctx *pulumi.Context, role *iam.Role, logPolicy *iam.RolePolicy, handlerFolder string) (*lambda.Function, error) {
	lambdaHandlerFileName := "handler"
	args := &lambda.FunctionArgs{
		Handler: pulumi.String(lambdaHandlerFileName),
		Role:    role.Arn,
		Runtime: pulumi.String("go1.x"),
		Code:    pulumi.NewFileArchive(fmt.Sprintf("./build/%s/%s.zip", handlerFolder, lambdaHandlerFileName)),
	}

	// Create the lambda using the args.
	lambdaFunction, err := lambda.NewFunction(
		ctx,
		GetResourceName(handlerFolder),
		args,
		pulumi.DependsOn([]pulumi.Resource{logPolicy}),
	)
	if err != nil {
		return &lambda.Function{}, err
	}

	return lambdaFunction, nil
}