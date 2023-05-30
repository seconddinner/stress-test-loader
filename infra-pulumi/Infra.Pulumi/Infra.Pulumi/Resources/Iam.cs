using Pulumi;
using Aws = Pulumi.Aws;

namespace Infra.Pulumi.Resources;

class Iam : ComponentResource
{
    [Output]
    public Output<string> StressTestClientReadProfileName { get; set; }
    public Iam(string name, ComponentResourceOptions opts = null) : base("stl:aws:Iam", name, opts)
    {
        // Create an IAM
        var currentCallerIdentity = Output.Create(Aws.GetCallerIdentity.InvokeAsync());
        var currentRegion = Output.Create(Aws.GetRegion.InvokeAsync());
        var stressTestClientReadRole = new Aws.Iam.Role("stressTestClientReadRole", new Aws.Iam.RoleArgs
        {
            AssumeRolePolicy = @"{
  ""Version"": ""2012-10-17"",
  ""Statement"": [
    {
      ""Action"": ""sts:AssumeRole"",
      ""Principal"": {
        ""Service"": ""ec2.amazonaws.com""
      },
      ""Effect"": ""Allow"",
      ""Sid"": """"
    }
  ]
}
",
            Tags = 
            {
                { "tag-key", "stress_test" },
            },
        });
        var stressTestClientReadProfile = new Aws.Iam.InstanceProfile("stressTestClientReadProfile", new Aws.Iam.InstanceProfileArgs
        {
            Role = stressTestClientReadRole.Name,
        });
        this.StressTestClientReadProfileName = stressTestClientReadProfile.Name;
        var stressTestClientRead = new Aws.Iam.RolePolicy("stressTestClientRead", new Aws.Iam.RolePolicyArgs
        {
            Role = stressTestClientReadRole.Id,
            Policy = @"{
  ""Version"": ""2012-10-17"",
  ""Statement"": [
    {
      ""Action"": [
        ""s3:GetObject""
      ],
      ""Effect"": ""Allow"",
      ""Resource"": ""arn:aws:s3:::stresstest-client/*""
    },
    {
      ""Action"": [
        ""s3:PutObject""
      ],
      ""Effect"": ""Allow"",
      ""Resource"": ""arn:aws:s3:::stresstest-client-log/*""
    }
  ]
}
",
        });

    }
}

