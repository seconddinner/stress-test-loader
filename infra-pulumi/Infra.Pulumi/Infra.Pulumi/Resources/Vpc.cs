﻿using Pulumi;
using Aws = Pulumi.Aws;
using System.Collections.Immutable;

namespace Infra.Pulumi.Resources;

class Vpc : ComponentResource
{
    [Output]
    public Output<string> MainVpcId { get; set; }
    [Output]
    public Output<string> DefaultSecurityGroupId { get; set; }
    [Output]
    public Output<ImmutableArray<string>> MainSubnetIds { get; set; }

    public Vpc(string name, ComponentResourceOptions opts = null) : base("stl:aws:Vpc", name, null)
    {
        // Set up Config
        var config = new Config();

        // Set up a VPC
        var available = Output.Create(Aws.GetAvailabilityZones.InvokeAsync());
        var mainVpc = new Aws.Ec2.Vpc("mainVpc", new Aws.Ec2.VpcArgs
        {
            CidrBlock = config.Require("cidr_block"), // default 10.10.0.0/16
        });
        this.MainVpcId = mainVpc.Id;

        // Create a default security group for the VPC
        var defaultSecurityGroup = new Aws.Ec2.DefaultSecurityGroup(string.Format("stress_test_loader-instance-{0}", config.Require("environment")), new Aws.Ec2.DefaultSecurityGroupArgs
        {
            VpcId = mainVpc.Id,
            Ingress = 
            {
                new Aws.Ec2.Inputs.DefaultSecurityGroupIngressArgs
                {
                    Protocol = "-1",
                    Self = true,
                    FromPort = 0,
                    ToPort = 0,
                },
            },
            Egress = 
            {
                new Aws.Ec2.Inputs.DefaultSecurityGroupEgressArgs
                {
                    FromPort = 0,
                    ToPort = 0,
                    Protocol = "-1",
                    CidrBlocks = 
                    {
                        "0.0.0.0/0",
                    },
                },
            },
        });
        this.DefaultSecurityGroupId = defaultSecurityGroup.Id;
        // defaultSecurityGroup.Id.Apply(id => {
        //     Console.WriteLine($"Default Security Group ID: {id}");
        //     return id;
        // });

        var mainSubnet = new List<Aws.Ec2.Subnet>();
        InputList<string> mainSubnetIds = Output.Create(new List<string>());
        for (var rangeIndex = 0; rangeIndex < int.Parse(config.Require("az_count")); rangeIndex++)
        {
            var range = new { Value = rangeIndex };
            var subnet = new Aws.Ec2.Subnet($"mainSubnet-{range.Value}", new Aws.Ec2.SubnetArgs
            {
                AvailabilityZone = available.Apply(available => available.Names[range.Value]),
                VpcId = mainVpc.Id,
                MapPublicIpOnLaunch = bool.Parse(config.Require("public_ip_on_launch")),
                CidrBlock = ReplaceIPandCIDR(mainVpc.CidrBlock, rangeIndex), // e.g., 10.10.0.0/22, 10.10.4.0/22
            });
            mainSubnet.Add(subnet);
            mainSubnetIds.Add(subnet.Id);
        }
        this.MainSubnetIds = mainSubnetIds;
        
        var gw = new Aws.Ec2.InternetGateway("gw", new Aws.Ec2.InternetGatewayArgs
        {
            VpcId = mainVpc.Id,
        });
        var routeTable = new Aws.Ec2.RouteTable("routeTable", new Aws.Ec2.RouteTableArgs
        {
            VpcId = mainVpc.Id,
            Routes = 
            {
                new Aws.Ec2.Inputs.RouteTableRouteArgs
                {
                    CidrBlock = "0.0.0.0/0",
                    GatewayId = gw.Id,
                },
            },
        });
        var routeTableAssociation = new List<Aws.Ec2.RouteTableAssociation>();
        for (var rangeIndex = 0; rangeIndex < int.Parse(config.Require("az_count")); rangeIndex++)
        {
            var range = new { Value = rangeIndex };
            routeTableAssociation.Add(new Aws.Ec2.RouteTableAssociation($"routeTableAssociation-{range.Value}", new Aws.Ec2.RouteTableAssociationArgs
            {
                SubnetId = mainSubnet.Select(__item => __item.Id).ToList()[range.Value],
                RouteTableId = routeTable.Id,
            }));
        }
    }

    public Output<string> ReplaceIPandCIDR(Output<string> input, int rangeIndex)
    {
        return input.Apply(value =>
        {
            // Split the IP address and CIDR notation
            string[] parts = value.Split('.');
            string[] cidrParts = parts[3].Split('/');

            // Replace the third number in the IP address
            parts[2] = (rangeIndex * 4).ToString();

            // Replace the number after the slash
            cidrParts[1] = "22";

            // Combine the parts back into a string
            string result = string.Join(".", parts[0], parts[1], parts[2], cidrParts[0]) + "/" + cidrParts[1];

            return result;
        });
    }

}