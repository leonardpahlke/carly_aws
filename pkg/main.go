package pkg

import (
	"fmt"
	"github.com/pulumi/pulumi-aws/sdk/v3/go/aws"
	"github.com/pulumi/pulumi-aws/sdk/v3/go/aws/iam"
	"github.com/pulumi/pulumi-aws/sdk/v3/go/aws/lambda"
	"github.com/pulumi/pulumi/sdk/v2/go/pulumi"
	"os"
)

const ProjectName = "carly"
const DeploymentEnv = "dev"

func GetResourceName(name string) string {
	return fmt.Sprintf("%s-%s-%s", ProjectName, DeploymentEnv, name)
}

func GetTags(resourceName string) pulumi.StringMap {
	return pulumi.StringMap{
		"STAGE":      pulumi.String(DeploymentEnv),
		"RESOURCE":   pulumi.String(resourceName),
		"CREATED_BY": pulumi.String("Pulumi"),
		"PROJECT":    pulumi.String(ProjectName),
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

func CheckEnvNotEmpty(key string) (string, bool) {
	envString, isEmpty := CheckEnv(key, "")
	if isEmpty {
		LogWarning("CheckEnvNotEmpty", fmt.Sprintf("environment variable %s is empty", key))
	}
	return envString, isEmpty
}

func CheckEnv(key string, expectedVal string) (string, bool) {
	env := os.Getenv(key)
	return env, env == expectedVal
}

const DefaultLambdaTimeout = 3

// Create a lambda function
func BuildLambdaFunction(ctx *pulumi.Context, config BuildLambdaConfig) (*lambda.Function, error) {
	lambdaHandlerFileName := "handler"
	args := &lambda.FunctionArgs{
		Handler:     pulumi.String(lambdaHandlerFileName),
		Role:        config.Role.Arn,
		Runtime:     pulumi.String("go1.x"),
		Code:        pulumi.NewFileArchive(fmt.Sprintf("./build/%s/%s.zip", config.HandlerFolder, lambdaHandlerFileName)),
		Environment: lambda.FunctionEnvironmentArgs{Variables: config.Env},
		//VpcConfig: lambda.FunctionVpcConfigArgs{
		//	SecurityGroupIds: pulumi.StringArray{config.SecurityGroupId},
		//	SubnetIds:        pulumi.StringArray{config.SubnetId},
		//	VpcId:            config.VpcId,
		//},
		Timeout: pulumi.Int(config.Timeout),
		Tags:    GetTags(lambdaHandlerFileName),
	}

	// Create the lambda using the args.
	lambdaFunction, err := lambda.NewFunction(
		ctx,
		GetResourceName(config.HandlerFolder),
		args,
		pulumi.DependsOn([]pulumi.Resource{config.LogPolicy}),
	)
	if err != nil {
		return &lambda.Function{}, err
	}

	return lambdaFunction, nil
}

type BuildLambdaConfig struct {
	Role          *iam.Role
	LogPolicy     *iam.RolePolicy
	Env           pulumi.StringMap
	HandlerFolder string
	Timeout       int
	// network configuration
	//VpcId pulumi.IDOutput
	//SecurityGroupId pulumi.IDOutput
	//SubnetId pulumi.IDOutput
}
