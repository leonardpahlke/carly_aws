package shared

import (
	"fmt"

	"github.com/pulumi/pulumi/sdk/v2/go/pulumi"
)

const ProjectName = "carly"
const DeploymentEnv = "dev"

func GetResourceName(name string) string {
	return fmt.Sprintf("%s-%s-%s", ProjectName, DeploymentEnv, name)
}

func GetTags(resourceName string) pulumi.StringMapInput {
	return pulumi.StringMapInput(pulumi.StringMap{
		"STAGE":      pulumi.String(DeploymentEnv),
		"RESOURCE":   pulumi.String(resourceName),
		"CREATED_BY": pulumi.String("Pulumi"),
		"PROJECT":    pulumi.String(ProjectName),
	})
}
