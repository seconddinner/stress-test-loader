using Pulumi;
using Aws = Pulumi.Aws;
using System.IO;
using System.Text.Json;

namespace Infra.Pulumi.Resources;

class Autoscaling : ComponentResource
{
    [Output]
    public Output<string> WebAutoScalingGroupName { get; set; }
    [Output]
    public Output<string> KeyPairName { get; set; }

    public Autoscaling(string name, Input<string> amiId, Input<string> stressTestClientReadProfileName,
        Input<string> mainVpcId, Input<string> defaultSecurityGroupId,
        InputList<string> mainSubnetIds, ComponentResourceOptions opts = null) : base("stl:aws:Autoscaling", name, null)
    {
        // Set up Config
        var config = new Config();
        
        // Create an AutoScaling Group
        var sdStresstest = new Aws.Ec2.KeyPair("sdStresstest", new Aws.Ec2.KeyPairArgs
        {
            PublicKey = config.Get("public_key") ?? Environment.GetEnvironmentVariable("public_key"),
        });
        this.KeyPairName = sdStresstest.KeyName;

        var stress_test_loader_allowed_cidr = config.Require("stress_test_loader_allowed_cidr");
        var cirdBlocks = new InputList<string>
        {
            stress_test_loader_allowed_cidr,
        };
        var prometheus_node_allowed_cidr = config.Require("prometheus_node_allowed_cidr");
        var prometheusCirdBlocks = new InputList<string>
        {
            prometheus_node_allowed_cidr,
        };
        var stresstestSecurityGroup = new Aws.Ec2.SecurityGroup("stresstestSecurityGroup", new Aws.Ec2.SecurityGroupArgs
        {
            VpcId = mainVpcId,
            Ingress = 
            {
                new Aws.Ec2.Inputs.SecurityGroupIngressArgs
                {
                    FromPort = 0,
                    ToPort = 8,
                    Protocol = "icmp",
                    CidrBlocks = cirdBlocks,
                },
                new Aws.Ec2.Inputs.SecurityGroupIngressArgs
                {
                    FromPort = 22,
                    ToPort = 22,
                    Protocol = "TCP",
                    CidrBlocks = cirdBlocks,
                },
                new Aws.Ec2.Inputs.SecurityGroupIngressArgs
                {
                    FromPort = int.Parse(config.Require("stress_test_loader_port")),
                    ToPort = int.Parse(config.Require("stress_test_loader_port")),
                    Protocol = "TCP",
                    CidrBlocks = cirdBlocks,
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
                    FromPort = int.Parse(config.Require("stress_test_loader_port")),
                    ToPort = int.Parse(config.Require("stress_test_loader_port")),
                    Protocol = "TCP",
                    CidrBlocks = 
                    {
                        config.Require("cidr_block"),
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
                        config.Require("egress_allowed_cidr"),
                    },
                },
            },
        });
        
        // Prepare user data
        var userData = File.ReadAllText("./Resources/stressclient-template.sh");
        userData = userData.Replace("${stress_test_loader_allowed_cidr}", config.Require("stress_test_loader_allowed_cidr"));
        userData = userData.Replace("${stress_test_loader_port}", config.Require("stress_test_loader_port"));
        userData = userData.Replace("${environment}", config.Require("environment"));
        userData = userData.Replace("${telegraf_username}", config.Require("telegraf_username"));
        userData = userData.Replace("${telegraf_password}", config.Require("telegraf_password"));
        userData = userData.Replace("${telegraf_url}", config.Require("telegraf_url"));

        var stressTestLoaderLaunchConfiguration = new Aws.Ec2.LaunchConfiguration("stressTestLoaderLaunchConfiguration", new Aws.Ec2.LaunchConfigurationArgs
        {
            NamePrefix = $"stress_test_loader-{config.Require("environment")}",
            ImageId = amiId.Apply(x => x!.ToString())!,
            InstanceType = config.Require("instance_type"),
            // DONE: iam
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
            // DONE: user data
            UserData = userData,
        });

        // var vpcZoneIdentifiers = new InputList<string>();
        // vpcZoneIdentifiers.AddRange(mainSubnet.Select(__item => __item.Id).ToList());

        var stressTestLoaderGroup = new Aws.AutoScaling.Group("stressTestLoaderGroup", new Aws.AutoScaling.GroupArgs
        {
            LaunchConfiguration = stressTestLoaderLaunchConfiguration.Name,
            VpcZoneIdentifiers = mainSubnetIds,
            HealthCheckType = "EC2",
            MinSize = int.Parse(config.Require("min_size")),
            MaxSize = int.Parse(config.Require("max_size")),
            DesiredCapacity = int.Parse(config.Require("desired_capacity")),
            Tags = 
            {
                new Aws.AutoScaling.Inputs.GroupTagArgs
                {
                    Key = "stress-test-nodes",
                    Value = "stress-test-nodes",
                    PropagateAtLaunch = true,
                },
            },
            // DONE: instance_refresh
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
            // DONE: enabled_metrics
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
        });
        this.WebAutoScalingGroupName = stressTestLoaderGroup.Name;

        // auto scale up policy
        var stressTestLoaderUpPolicy = new Aws.AutoScaling.Policy("stressTestLoaderUpPolicy", new Aws.AutoScaling.PolicyArgs
        {
            ScalingAdjustment = int.Parse(config.Require("up_scaling_adjustment")),
            AdjustmentType = "ChangeInCapacity",
            Cooldown = 300,
            AutoscalingGroupName = stressTestLoaderGroup.Name,
        });
        // auto scale down policy
        var stressTestLoaderDownPolicy = new Aws.AutoScaling.Policy("stressTestLoaderDownPolicy", new Aws.AutoScaling.PolicyArgs
        {
            ScalingAdjustment = int.Parse(config.Require("down_scaling_adjustment")),
            AdjustmentType = "ChangeInCapacity",
            Cooldown = 300,
            AutoscalingGroupName = stressTestLoaderGroup.Name,
        });
        var stressTestLoaderUpMetricAlarm = new Aws.CloudWatch.MetricAlarm("stressTestLoaderUpMetricAlarm", new Aws.CloudWatch.MetricAlarmArgs
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
        });
        var stressTestLoaderDownMetricAlarm = new Aws.CloudWatch.MetricAlarm("stressTestLoaderDownMetricAlarm", new Aws.CloudWatch.MetricAlarmArgs
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
        });
        
        var ec2Instances = Aws.Ec2.GetInstances.Invoke(new()
        {
            Filters = new[]
            {
                new Aws.Ec2.Inputs.GetInstancesFilterInputArgs
                {
                    Name = "image-id",
                    Values = new[]
                    {
                        amiId.Apply(x => x!.ToString())!
                    }
                }
            }
        });
   
        var ec2PublicIps = ec2Instances.Apply(instance => instance.PublicIps);
        var instancesList = new List<List<Dictionary<string, string>>>();

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
      
            // Convert the instances to JSON
            var json = JsonSerializer.Serialize(instancesList);

            // Save the JSON to a file
            var filePath = "/tmp/IP.json";
            File.WriteAllText(filePath, json);

            Console.WriteLine($"IP addresses of EC2 instances saved to '{filePath}'.");

            return ips;
        });
    }

}