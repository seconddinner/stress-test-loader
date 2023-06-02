using Pulumi;
using Aws = Pulumi.Aws;

namespace Infra.Pulumi.Resources;

class Iam : ComponentResource
{
    [Output]
    public Output<string> StressTestClientReadProfileName { get; set; }
    public Iam(string name, Aws.Provider provider, string region, ComponentResourceOptions opts = null) : base("stl:aws:Iam", name, opts)
    {
        // Set up Config
        var config = new Config();
        var s3_client_bucket_name = config.Require("s3_client_bucket_name");
        var s3_log_bucket_name = config.Require("s3_log_bucket_name");
        
        // Create an IAM
        var currentCallerIdentity = Output.Create(Aws.GetCallerIdentity.InvokeAsync());
        var currentRegion = Output.Create(Aws.GetRegion.InvokeAsync());
        var stressTestClientReadRole = new Aws.Iam.Role("stressTestClientReadRole-" + region, new Aws.Iam.RoleArgs
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
        }, new CustomResourceOptions
        {
          Provider = provider,
          Parent = this
        });
        var stressTestClientReadProfile = new Aws.Iam.InstanceProfile("stressTestClientReadProfile-" + region, new Aws.Iam.InstanceProfileArgs
        {
            Role = stressTestClientReadRole.Name,
        }, new CustomResourceOptions
        {
          Provider = provider,
          Parent = this
        });
        this.StressTestClientReadProfileName = stressTestClientReadProfile.Name;
        var stressTestClientRead = new Aws.Iam.RolePolicy("stressTestClientRead-" + region, new Aws.Iam.RolePolicyArgs
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
      ""Resource"": ""arn:aws:s3:::" + s3_client_bucket_name + @"/*""
    },
    {
      ""Action"": [
        ""s3:PutObject""
      ],
      ""Effect"": ""Allow"",
      ""Resource"": ""arn:aws:s3:::" + s3_log_bucket_name + @"/*""
    }
  ]
}
",
        }, new CustomResourceOptions
        {
          Provider = provider,
          Parent = this
        });
        
        RegisterOutputs();
    }
}

