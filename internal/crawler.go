package internal

import (
	"carly_aws/pkg"
	"github.com/pulumi/pulumi-aws/sdk/v3/go/aws/ec2"
	"github.com/pulumi/pulumi/sdk/v2/go/pulumi"
)

func CreateCrawler(ctx *pulumi.Context, config CrawlerConfig) (CrawlerData, error) {
	amazon2AmiHvm, err := ec2.GetAmi(
		ctx,
		pkg.GetResourceName("amzn2-ami-hvm"),
		pulumi.ID("ami-0a6dc7529cd559185"),
		&ec2.AmiState{
			Arn: pulumi.String("ami-0a6dc7529cd559185"),
		})
	if err != nil {
		return CrawlerData{}, err
	}
	crawlerInstance, err := ec2.NewInstance(ctx, pkg.GetResourceName("ec2-crawler"), &ec2.InstanceArgs{
		Ami:          amazon2AmiHvm.ID(),
		InstanceType: pulumi.String("t2.micro"),
		SubnetId: config.CrawlerSubnet.ID(),
		VpcSecurityGroupIds: config.CrawlerVpcSecurityGroups,
		Tags: pkg.GetTags("Crawler"),
	})
	if err != nil {
		return CrawlerData{}, err
	}

	return CrawlerData{
		CrawlerInstance: crawlerInstance,
	}, nil
}

type CrawlerConfig struct {
	CrawlerSubnet *ec2.Subnet
	CrawlerVpcSecurityGroups pulumi.StringArray
}

type CrawlerData struct {
	CrawlerInstance *ec2.Instance
}