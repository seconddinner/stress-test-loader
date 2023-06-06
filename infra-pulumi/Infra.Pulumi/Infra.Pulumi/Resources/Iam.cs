using System.Text.Json;
using Amazon.Auth.AccessControlPolicy;
using Pulumi;
using Pulumi.Aws.Iam.Inputs;
using Aws = Pulumi.Aws;
using Policy = Amazon.Auth.AccessControlPolicy.Policy;
using Principal = Amazon.Auth.AccessControlPolicy.Principal;
using Resource = Amazon.Auth.AccessControlPolicy.Resource;
using Statement = Amazon.Auth.AccessControlPolicy.Statement;

namespace Infra.Pulumi.Resources;

class Iam : ComponentResource
{
    public Output<string> StressTestClientReadProfileName { get; set; }
    public Iam(StressConfig cfg, ComponentResourceOptions opts = null) : base("stl:aws:Iam", $"stl-iam-{cfg.CurrentRegion}", opts)
    {
        // Create an IAM
        var assumeRolePolicy = new Policy().WithStatements(
            new Statement(Statement.StatementEffect.Allow)
                .WithId("AssumeRoleStatement")
                .WithActionIdentifiers(new ActionIdentifier("sts:AssumeRole"))
                .WithPrincipals(new Principal(Principal.SERVICE_PROVIDER, "ec2.amazonaws.com")));
        var inlinePolicy = new Policy().WithStatements(
            new Statement(Statement.StatementEffect.Allow)
                .WithId("GetObjectsStatement")
                .WithActionIdentifiers(new ActionIdentifier("s3:GetObject"))
                .WithResources(new Resource($"arn:aws:s3:::{cfg.S3ClientBucketName}/*")),
            new Statement(Statement.StatementEffect.Allow)
                .WithId("PutObjectsStatement")
                .WithActionIdentifiers(new ActionIdentifier("s3:PutObject"))
                .WithResources(new Resource($"arn:aws:s3:::{cfg.S3LogBucketName}/*")));
        
        var stressTestClientReadRole = new Aws.Iam.Role("stressTestClientReadRole-" + cfg.CurrentRegion, new Aws.Iam.RoleArgs
        {
            AssumeRolePolicy = assumeRolePolicy.ToJson(),
            Tags = 
            {
                { "tag-key", "stress_test" },
            },
            InlinePolicies = new InputList<RoleInlinePolicyArgs>
            {
                new RoleInlinePolicyArgs
                {
                    Name = $"{cfg.CurrentRegion}-CrossRegionRole",
                    Policy = inlinePolicy.ToJson(),
                }
            }
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

        RegisterOutputs();
    }
}

