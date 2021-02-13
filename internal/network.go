package internal

import (
	"carly_aws/pkg"
	"github.com/pulumi/pulumi-aws/sdk/v3/go/aws/ec2"
	"github.com/pulumi/pulumi/sdk/v2/go/pulumi"
)

func CreateNetwork(ctx *pulumi.Context, _ NetworkConfig) (NetworkData, error) {
	// VPC
	vpc, err := ec2.NewVpc(ctx, pkg.GetResourceName("vpc"), &ec2.VpcArgs{
		CidrBlock: pulumi.String("10.0.0.0/16"),
		Tags:      pkg.GetTags("Vpc"),
	})
	if err != nil {
		return NetworkData{}, err
	}

	// Security Group
	crawlerSgName := pkg.GetResourceName("crawler-security-group")
	crawlerSecurityGroup, err := ec2.NewSecurityGroup(ctx, crawlerSgName, &ec2.SecurityGroupArgs{
		VpcId: vpc.ID(),
		Ingress: ec2.SecurityGroupIngressArray{
			ec2.SecurityGroupIngressArgs{
				Protocol:   pulumi.String("tcp"),
				FromPort:   pulumi.Int(80),
				ToPort:     pulumi.Int(80),
				CidrBlocks: pulumi.StringArray{pulumi.String("10.0.2.0/24")},
			},
		},
		Tags: pkg.GetTags("SecurityGroup"),
		Name: pulumi.String(crawlerSgName),
	})
	if err != nil {
		return NetworkData{}, err
	}

	// Internet Gateway
	igw, err := ec2.NewInternetGateway(ctx, pkg.GetResourceName("igw"), &ec2.InternetGatewayArgs{
		VpcId: vpc.ID(),
		Tags:  pkg.GetTags("InternetGateway"),
	})
	if err != nil {
		return NetworkData{}, err
	}

	publicSubnet, err := ec2.NewSubnet(ctx, pkg.GetResourceName("public-subnet"), &ec2.SubnetArgs{
		VpcId:     vpc.ID(),
		CidrBlock: pulumi.String("10.0.1.0/24"),
		Tags:      pkg.GetTags("PublicSubnet"),
	})
	if err != nil {
		return NetworkData{}, err
	}

	privateSubnet, err := ec2.NewSubnet(ctx, pkg.GetResourceName("private-subnet"), &ec2.SubnetArgs{
		VpcId:     vpc.ID(),
		CidrBlock: pulumi.String("10.0.2.0/24"),
		Tags:      pkg.GetTags("PrivateSubnet"),
	})
	if err != nil {
		return NetworkData{}, err
	}

	// Route Table
	_, err = ec2.NewRouteTable(ctx, pkg.GetResourceName("route-table"), &ec2.RouteTableArgs{
		Routes: ec2.RouteTableRouteArray{
			ec2.RouteTableRouteArgs{
				CidrBlock: pulumi.String("0.0.0.0/0"),
				GatewayId: igw.ID(),
			},
		},
		VpcId: vpc.ID(),
		Tags:  pkg.GetTags("RouteTable"),
	})
	if err != nil {
		return NetworkData{}, err
	}

	return NetworkData{
		Vpc: vpc,
		PublicSubnet: publicSubnet,
		PrivateSubnet: privateSubnet,
		CrawlerSecurityGroup: crawlerSecurityGroup,
	}, nil
}

type NetworkConfig struct {}

type NetworkData struct {
	Vpc *ec2.Vpc
	PublicSubnet *ec2.Subnet
	PrivateSubnet *ec2.Subnet
	CrawlerSecurityGroup *ec2.SecurityGroup
}