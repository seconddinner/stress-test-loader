using Pulumi;
using Aws = Pulumi.Aws;
using System.Collections.Immutable;
using Pulumi.Aws.Inputs;

namespace Infra.Pulumi.Resources;

class Vpc : ComponentResource
{
    public Output<string> MainVpcId { get; set; }
    
    public Output<ImmutableArray<string>> MainSubnetIds { get; set; }

    public Vpc(StressConfig cfg, ComponentResourceOptions opts = null) : base("stl:aws:Vpc", $"stl-vpc-{cfg.CurrentRegion}", opts)
    {
        // Set up a VPC
        var available = Output.Create(Aws.GetAvailabilityZones.InvokeAsync(new Aws.GetAvailabilityZonesArgs()
        {
            Filters = 
            {
                new GetAvailabilityZonesFilterArgs()
                {
                    Name = "region-name",
                    Values = { cfg.CurrentRegion }
                }
            }
        }, new InvokeOptions
        {
            Parent = this
        }));
        
        var mainVpc = new Aws.Ec2.Vpc("mainVpc-" + cfg.CurrentRegion, new Aws.Ec2.VpcArgs
        {
            CidrBlock = cfg.CidrBlock
        }, new CustomResourceOptions
        {
            Parent = this
        });
        this.MainVpcId = mainVpc.Id;

        // Create a default security group for the VPC
        var defaultSecurityGroup = new Aws.Ec2.DefaultSecurityGroup(
            $"stress_test_loader-instance-{cfg.Environment}-{cfg.CurrentRegion}", 
            new Aws.Ec2.DefaultSecurityGroupArgs
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
        }, new CustomResourceOptions
        {
            Parent = this
        });

        var mainSubnet = new List<Aws.Ec2.Subnet>();
        InputList<string> mainSubnetIds = Output.Create(new List<string>());
        for (var rangeIndex = 0; rangeIndex < cfg.AzCount; rangeIndex++)
        {
            var range = new { Value = rangeIndex };
            var subnet = new Aws.Ec2.Subnet($"mainSubnet-{range.Value}-" + cfg.CurrentRegion, new Aws.Ec2.SubnetArgs
            {
                AvailabilityZone = available.Apply(av => av.Names[range.Value]),
                VpcId = mainVpc.Id,
                MapPublicIpOnLaunch = cfg.PublicIpOnLaunch,
                CidrBlock = ReplaceIPandCIDR(mainVpc.CidrBlock, rangeIndex), // e.g., 10.10.0.0/22, 10.10.4.0/22
            }, new CustomResourceOptions
            {
                Parent = this
            });
            mainSubnet.Add(subnet);
            mainSubnetIds.Add(subnet.Id);
        }
        this.MainSubnetIds = mainSubnetIds;
        
        var gw = new Aws.Ec2.InternetGateway("gw-" + cfg.CurrentRegion, new Aws.Ec2.InternetGatewayArgs
        {
            VpcId = mainVpc.Id,
        }, new CustomResourceOptions
        {
            Parent = this
        });
        var routeTable = new Aws.Ec2.RouteTable("routeTable-" + cfg.CurrentRegion, new Aws.Ec2.RouteTableArgs
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
        }, new CustomResourceOptions
        {
            Parent = this
        });
        var routeTableAssociation = new List<Aws.Ec2.RouteTableAssociation>();
        for (var rangeIndex = 0; rangeIndex < cfg.AzCount; rangeIndex++)
        {
            var range = new { Value = rangeIndex };
            routeTableAssociation.Add(new Aws.Ec2.RouteTableAssociation(
                $"routeTableAssociation-{range.Value}-" + cfg.CurrentRegion, new Aws.Ec2.RouteTableAssociationArgs
            {
                SubnetId = mainSubnet.Select(__item => __item.Id).ToList()[range.Value],
                RouteTableId = routeTable.Id,
            }, new CustomResourceOptions
            {
                Parent = this
            }));
        }
        
        RegisterOutputs();
    }

    private Output<string> ReplaceIPandCIDR(Output<string> input, int rangeIndex)
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