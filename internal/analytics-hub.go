package internal

import (
	"carly_aws/shared"

	"github.com/pulumi/pulumi-aws/sdk/v3/go/aws/ec2"
	"github.com/pulumi/pulumi/sdk/v2/go/pulumi"
)

func CreateAnalyticsHub(ctx *pulumi.Context, config AnalyticsHubConfig) (AnalyticsHubData, error) {
	amazon2AmiHvm, err := ec2.GetAmi(
		ctx,
		shared.GetResourceName("amzn2-ami-hvm"),
		pulumi.ID("ami-0a6dc7529cd559185"),
		&ec2.AmiState{
			Arn: pulumi.String("ami-0a6dc7529cd559185"),
		})
	if err != nil {
		return AnalyticsHubData{}, err
	}
	crawlerInstance, err := ec2.NewInstance(ctx, shared.GetResourceName("ec2-crawler"), &ec2.InstanceArgs{
		Ami:                 amazon2AmiHvm.ID(),
		InstanceType:        pulumi.String("t2.micro"),
		SubnetId:            config.CrawlerSubnet.ID(),
		VpcSecurityGroupIds: config.CrawlerVpcSecurityGroups,
		Tags:                shared.GetTags("Crawler"),
	})
	if err != nil {
		return AnalyticsHubData{}, err
	}

	return AnalyticsHubData{
		CrawlerInstance: crawlerInstance,
	}, nil
}

type AnalyticsHubConfig struct {
	CrawlerSubnet            *ec2.Subnet
	CrawlerVpcSecurityGroups pulumi.StringArray
}

type AnalyticsHubData struct {
	CrawlerInstance *ec2.Instance
}
