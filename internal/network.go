package internal

import (
	"github.com/pulumi/pulumi-aws/sdk/v3/go/aws/ec2"
	"github.com/pulumi/pulumi/sdk/v2/go/pulumi"
)

func CreateNetwork(ctx *pulumi.Context, _ NetworkConfig) (NetworkData, error) {
	// Security Group
	_, err := ec2.NewSecurityGroup(ctx, "mysecuritygroup", &ec2.SecurityGroupArgs{
		Ingress: ec2.SecurityGroupIngressArray{
			ec2.SecurityGroupIngressArgs{
				Protocol:   pulumi.String("tcp"),
				FromPort:   pulumi.Int(80),
				ToPort:     pulumi.Int(80),
				CidrBlocks: pulumi.StringArray{pulumi.String("0.0.0.0/0")},
			},
		},
	})
	if err != nil {
		return NetworkData{}, err
	}

	// VPC
	vpc, err := ec2.NewVpc(ctx, "myvpc", &ec2.VpcArgs{
		CidrBlock: pulumi.String("10.0.0.0/16"),
	})
	if err != nil {
		return NetworkData{}, err
	}

	// Internet Gateway
	igw, err := ec2.NewInternetGateway(ctx, "myinternetgateway", &ec2.InternetGatewayArgs{
		VpcId: vpc.ID(),
	})
	if err != nil {
		return NetworkData{}, err
	}

	// Route Table
	_, err = ec2.NewRouteTable(ctx, "myroutetable", &ec2.RouteTableArgs{
		Routes: ec2.RouteTableRouteArray{
			ec2.RouteTableRouteArgs{
				CidrBlock: pulumi.String("0.0.0.0/0"),
				GatewayId: igw.ID(),
			},
		},
		VpcId: vpc.ID(),
	})
	if err != nil {
		return NetworkData{}, err
	}

	return NetworkData{}, nil
}

type NetworkConfig struct {}

type NetworkData struct {}