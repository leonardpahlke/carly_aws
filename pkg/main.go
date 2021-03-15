package pkg

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"strings"

	sdk_aws "github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
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

func ListContains(element string, list []string) bool {
	for _, elem := range list {
		if element == elem {
			return true
		}
	}
	return false
}

func GetBucketFileName(bucketName string, newspaper string, filename string, fileEnding string) string {
	return fmt.Sprintf("%s/%s/%s.%s", bucketName, newspaper, filename, fileEnding)
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
}

func CreateKeyValuePairs(m map[string]string) string {
	b := new(bytes.Buffer)
	valEscapeString := ""
	_, _ = fmt.Fprint(b, "{\n  ")
	for key, value := range m {
		if strings.ContainsAny(value, "{}") {
			valEscapeString = ""
		} else {
			valEscapeString = "\""
		}
		_, _ = fmt.Fprintf(b, "  \"%s\": %s%s%s\n  ", key, valEscapeString, value, valEscapeString)
	}
	_, _ = fmt.Fprint(b, "}")
	return b.String()
}

func MarshalStruct(structIn interface{}) []byte {
	b, err := json.Marshal(structIn)
	if err != nil {
		fmt.Println(err)
		return nil
	}
	return b
}

// store file in article analytics bucket
func StoreFileInArticleAnalyticsBucket(fileInfoIn StoreFileStruct, fileMetaIn StoreFileMetaStruct) *s3manager.UploadOutput {
	result, err := fileMetaIn.Uploader.Upload(&s3manager.UploadInput{
		Bucket: sdk_aws.String(fileMetaIn.BucketName),
		Key: sdk_aws.String(fmt.Sprintf("%s/%s/%s.%s",
			fileMetaIn.Newspaper,
			fileMetaIn.ArticleReference,
			fileInfoIn.Filename,
			fileInfoIn.FileEnding)),
		Body: strings.NewReader(fileInfoIn.File),
	})
	if err != nil {
		LogError(fileMetaIn.SpiderName, "s3 upload error", err)
		return &s3manager.UploadOutput{}
	}
	return result
}

func GetBucketKeyForAnalyticsBucket(newspaper string, articleReference string, fileName string, fileEnding string) string {
	return fmt.Sprintf("%s/%s/%s.%s",
		newspaper,
		articleReference,
		fileName,
		fileEnding)
}

func GetBucketUriForKey(key string) string {
	return fmt.Sprintf("s3://%s", key)
}

/*
	result, err := uploader.Upload(&s3manager.UploadInput{
		Bucket: aws.String(s3BucketName),
		Key:    aws.String(fmt.Sprintf("%s/%s", event.Newspaper, fileName)),
		Body:   f,
	})
*/

// clean string removes
func TrimStringAry(inputStrings []string) []string {
	var cleanedString []string
	for _, str := range inputStrings {
		str = strings.TrimSpace(str)
		if str != "" {
			cleanedString = append(cleanedString, str)
		}
	}
	LogInfo("pkg.TrimStringAry", fmt.Sprintf("cleanedString.. %v", cleanedString))
	return cleanedString
}
