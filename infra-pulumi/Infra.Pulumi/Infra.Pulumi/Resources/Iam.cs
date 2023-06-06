using System.Text.Json;
using Amazon.Auth.AccessControlPolicy;
using Pulumi;
using Aws = Pulumi.Aws;

namespace Infra.Pulumi.Resources;

class Iam : ComponentResource
{
    public Output<string> StressTestClientReadProfileName { get; set; }
    public Iam(StressConfig cfg, ComponentResourceOptions opts = null) : base("stl:aws:Iam", $"stl-iam-{cfg.CurrentRegion}", opts)
    {
        // Create an IAM
        var stressTestClientReadRole = new Aws.Iam.Role("stressTestClientReadRole-" + cfg.CurrentRegion, new Aws.Iam.RoleArgs
        {
            AssumeRolePolicy = new Policy().WithStatements(
              new Statement(Statement.StatementEffect.Allow)
                .WithId(null)
                .WithActionIdentifiers(new ActionIdentifier("sts:AssumeRole"))
                .WithPrincipals(new Principal(Principal.SERVICE_PROVIDER, "ec2.amazonaws.com"))
                .WithConditions())
            .ToJson(),
            Tags = 
            {
                { "tag-key", "stress_test" },
            },
        }, new CustomResourceOptions
        {
          Parent = this
        });
        var stressTestClientReadProfile = new Aws.Iam.InstanceProfile("stressTestClientReadProfile-" + cfg.CurrentRegion, new Aws.Iam.InstanceProfileArgs
        {
            Role = stressTestClientReadRole.Name,
        }, new CustomResourceOptions
        {
          Parent = this
        });
        this.StressTestClientReadProfileName = stressTestClientReadProfile.Name;
        var stressTestClientRead = new Aws.Iam.RolePolicy("stressTestClientRead-" + cfg.CurrentRegion, new()
        {
            Role = stressTestClientReadRole.Id,
            Policy = JsonSerializer.Serialize(new Dictionary<string, object?>
            {
                ["Version"] = "2012-10-17",
                ["Statement"] = new[]
                {
                    new Dictionary<string, object?>
                    {
                        ["Action"] = new[]
                        {
                            "s3:GetObject",
                        },
                        ["Effect"] = "Allow",
                        ["Resource"] = $"arn:aws:s3:::{cfg.S3ClientBucketName}/*",
                    },
                    new Dictionary<string, object?>
                    {
                        ["Action"] = new[]
                        {
                            "s3:PutObject",
                        },
                        ["Effect"] = "Allow",
                        ["Resource"] = $"arn:aws:s3:::{cfg.S3LogBucketName}/*",
                    },
                },
            }),
        }, new CustomResourceOptions
        {
          Parent = this
        });
        
        RegisterOutputs();
    }
}

