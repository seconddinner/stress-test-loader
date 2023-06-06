using Pulumi;
using Aws = Pulumi.Aws;
using System.IO;
using System.Text.Json;

namespace Infra.Pulumi.Resources;

class Autoscaling : ComponentResource
{
    public Autoscaling(string name, string region, Input<string> amiId, 
        Input<string> stressTestClientReadProfileName, Input<string> mainVpcId,
        InputList<string> mainSubnetIds, StressConfig cfg, ComponentResourceOptions opts = null) 
        : base("stl:aws:Autoscaling", name, opts)
    {
        // Create an AutoScaling Group
        var sdStresstest = new Aws.Ec2.KeyPair("sdStresstest-" + region, new Aws.Ec2.KeyPairArgs
        {
            PublicKey = cfg.PublicKey ?? Environment.GetEnvironmentVariable("public_key"),
        }, new CustomResourceOptions
        {
            Parent = this
        });

        var cidrBlocks = new InputList<string>
        {
            cfg.AllowedCidrBlocks
        };
        var prometheusCirdBlocks = new InputList<string>
        {
            cfg.PrometheusAllowedCidrBlocks
        };
        var stresstestSecurityGroup = new Aws.Ec2.SecurityGroup("stresstestSecurityGroup-" + region, new Aws.Ec2.SecurityGroupArgs
        {
            VpcId = mainVpcId,
            Ingress = 
            {
                new Aws.Ec2.Inputs.SecurityGroupIngressArgs
                {
                    FromPort = 0,
                    ToPort = 8,
                    Protocol = "icmp",
                    CidrBlocks = cidrBlocks,
                },
                new Aws.Ec2.Inputs.SecurityGroupIngressArgs
                {
                    FromPort = 22,
                    ToPort = 22,
                    Protocol = "TCP",
                    CidrBlocks = cidrBlocks,
                },
                new Aws.Ec2.Inputs.SecurityGroupIngressArgs
                {
                    FromPort = cfg.StlPort,
                    ToPort = cfg.StlPort,
                    Protocol = "TCP",
                    CidrBlocks = cidrBlocks,
                },
                new Aws.Ec2.Inputs.SecurityGroupIngressArgs
                {
                    FromPort = 9100,
                    ToPort = 9100,
                    Protocol = "TCP",
                    CidrBlocks = prometheusCirdBlocks,
                },
                new Aws.Ec2.Inputs.SecurityGroupIngressArgs
                {
                    FromPort = 9301,
                    ToPort = 9301,
                    Protocol = "TCP",
                    CidrBlocks = prometheusCirdBlocks,
                },
                new Aws.Ec2.Inputs.SecurityGroupIngressArgs
                {
                    FromPort = cfg.StlPort,
                    ToPort = cfg.StlPort,
                    Protocol = "TCP",
                    CidrBlocks = 
                    {
                        cfg.CidrBlock
                    },
                },
            },
            Egress = 
            {
                new Aws.Ec2.Inputs.SecurityGroupEgressArgs
                {
                    FromPort = 0,
                    ToPort = 0,
                    Protocol = "-1",
                    CidrBlocks = 
                    {
                        cfg.EgressAllowedCidr
                    },
                },
            },
        }, new CustomResourceOptions
        {
            Parent = this
        });
        
        // Prepare user data
        var userData = File.ReadAllText("./Resources/stressclient-template.sh");
        userData = userData.Replace("${stress_test_loader_port}", cfg.StlPort.ToString());
        userData = userData.Replace("${environment}", cfg.Environment);
        userData = userData.Replace("${telegraf_username}", cfg.TelegrafUsername);
        userData = userData.Replace("${telegraf_password}", cfg.TelegrafPassword);
        userData = userData.Replace("${telegraf_url}", cfg.TelegrafUrl);

        var stressTestLoaderLaunchConfiguration = new Aws.Ec2.LaunchConfiguration("stressTestLoaderLaunchConfiguration-" + region, new Aws.Ec2.LaunchConfigurationArgs
        {
            NamePrefix = $"stress_test_loader-{cfg.Environment}",
            ImageId = amiId.Apply(x => x!.ToString())!,
            InstanceType = cfg.InstanceType,
            IamInstanceProfile = stressTestClientReadProfileName.Apply(x => x!.ToString())!,
            AssociatePublicIpAddress = true,
            KeyName = sdStresstest.KeyName,
            SecurityGroups = 
            {
                stresstestSecurityGroup.Id,
            },
            RootBlockDevice = new Aws.Ec2.Inputs.LaunchConfigurationRootBlockDeviceArgs
            {
                VolumeType = "gp2",
                VolumeSize = 30,
                DeleteOnTermination = true,
            },
            UserData = userData,
        }, new CustomResourceOptions
        {
            Parent = this
        });
        

        var stressTestLoaderGroup = new Aws.AutoScaling.Group("stressTestLoaderGroup-" + region, new Aws.AutoScaling.GroupArgs
        {
            LaunchConfiguration = stressTestLoaderLaunchConfiguration.Name,
            VpcZoneIdentifiers = mainSubnetIds,
            HealthCheckType = "EC2",
            MinSize = cfg.MinSize,
            MaxSize = cfg.MaxSize,
            DesiredCapacity = cfg.DesiredCapacity,
            Tags = 
            {
                new Aws.AutoScaling.Inputs.GroupTagArgs
                {
                    Key = "stress-test-nodes",
                    Value = "stress-test-nodes",
                    PropagateAtLaunch = true,
                },
            },
            InstanceRefresh = new Aws.AutoScaling.Inputs.GroupInstanceRefreshArgs
            {
                Strategy = "Rolling",
                Preferences = new Aws.AutoScaling.Inputs.GroupInstanceRefreshPreferencesArgs
                {
                    MinHealthyPercentage = 50,
                    InstanceWarmup = "300",
                },
                Triggers = new InputList<string> {"tag"},
            },
            EnabledMetrics = new InputList<string>
            {
                "GroupMinSize",
                "GroupMaxSize",
                "GroupDesiredCapacity",
                "GroupInServiceInstances",
                "GroupPendingInstances",
                "GroupStandbyInstances",
                "GroupTerminatingInstances",
                "GroupTotalInstances"
            },
        }, new CustomResourceOptions
        {
            Parent = this
        });
        StorePublicIpsToJson(stressTestLoaderGroup, amiId);
            
        // auto scale up policy
        var stressTestLoaderUpPolicy = new Aws.AutoScaling.Policy("stressTestLoaderUpPolicy-" + region, new Aws.AutoScaling.PolicyArgs
        {
            ScalingAdjustment = cfg.UpScalingAdjustment,
            AdjustmentType = "ChangeInCapacity",
            Cooldown = 300,
            AutoscalingGroupName = stressTestLoaderGroup.Name,
        }, new CustomResourceOptions
        {
            Parent = this
        });
        // auto scale down policy
        var stressTestLoaderDownPolicy = new Aws.AutoScaling.Policy("stressTestLoaderDownPolicy-" + region, new Aws.AutoScaling.PolicyArgs
        {
            ScalingAdjustment = cfg.DownScalingAdjustment,
            AdjustmentType = "ChangeInCapacity",
            Cooldown = 300,
            AutoscalingGroupName = stressTestLoaderGroup.Name,
        }, new CustomResourceOptions
        {
            Parent = this
        });
        var stressTestLoaderUpMetricAlarm = new Aws.CloudWatch.MetricAlarm("stressTestLoaderUpMetricAlarm-" + region, new Aws.CloudWatch.MetricAlarmArgs
        {
            ComparisonOperator = "GreaterThanOrEqualToThreshold",
            EvaluationPeriods = 2,
            MetricName = "CPUUtilization",
            Namespace = "AWS/EC2",
            Period = 120,
            Statistic = "Average",
            Dimensions = 
            {
                { "AutoScalingGroupName", stressTestLoaderGroup.Name },
            },
            Threshold = 120,
            AlarmDescription = "Check whether EC2 instance CPU utilisation is over 80% on average",
            AlarmActions = 
            {
                stressTestLoaderUpPolicy.Arn,
            },
        }, new CustomResourceOptions
        {
            Parent = this
        });
        var stressTestLoaderDownMetricAlarm = new Aws.CloudWatch.MetricAlarm("stressTestLoaderDownMetricAlarm-" + region, new Aws.CloudWatch.MetricAlarmArgs
        {
            ComparisonOperator = "LessThanOrEqualToThreshold",
            EvaluationPeriods = 120,
            MetricName = "CPUUtilization",
            Namespace = "AWS/EC2",
            Period = 120,
            Statistic = "Average",
            Threshold = 20,
            Dimensions = 
            {
                { "AutoScalingGroupName", stressTestLoaderGroup.Name },
            },
            AlarmDescription = "Check whether EC2 instance CPU utilisation is under 20% on average",
            AlarmActions = 
            {
                stressTestLoaderDownPolicy.Arn,
            },
        }, new CustomResourceOptions
        {
            Parent = this
        });
        
        RegisterOutputs();
    }

    private void StorePublicIpsToJson(Aws.AutoScaling.Group stressTestLoaderGroup, Input<string> amiId)
    {
        stressTestLoaderGroup.Name.Apply(n =>
        {
            InputList<string> ec2PublicIps = new InputList<string>();
            var ec2Instances = Aws.Ec2.GetInstances.Invoke(new()
            {
                Filters = new[]
                {
                    new Aws.Ec2.Inputs.GetInstancesFilterInputArgs
                    {
                        Name = "image-id",
                        Values = amiId
                    }
                }
            }, new InvokeOptions()
            {
                Parent = this
            });
            ec2PublicIps.AddRange(ec2Instances.Apply(instance => instance.PublicIps));

            var filePath = "/tmp/IP.json";
            var instancesList = new List<List<Dictionary<string, string>>>();
            if (File.Exists(filePath))
            {
                // If the file exists, read and parse it
                string existingJson = File.ReadAllText(filePath);
                instancesList = JsonSerializer.Deserialize<List<List<Dictionary<string, string>>>>(existingJson);
            }

            // Iterate over the instances and add their public IP addresses to the list
            ec2PublicIps.Apply(ips =>
            {
                foreach (var publicIp in ips)
                {
                    // Create a dictionary for each instance with the "public_ip" property
                    var instanceDict = new Dictionary<string, string>
                    {
                        { "public_ip", publicIp }
                    };
                    // Add the instance dictionary to the instances list
                    instancesList.Add(new List<Dictionary<string, string>> { instanceDict });
                }
                // Convert the instances to JSON & save the JSON to a file
                File.WriteAllText(filePath, JsonSerializer.Serialize(instancesList));
                Console.WriteLine($"IP addresses of EC2 instances saved to '{filePath}'.");
                return ips;
            });
            return n;
        });
    }

}