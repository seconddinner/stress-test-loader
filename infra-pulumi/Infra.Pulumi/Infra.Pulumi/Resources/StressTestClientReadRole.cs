using Amazon.Auth.AccessControlPolicy;
using Pulumi;
using Pulumi.Aws.Iam;
using Pulumi.Aws.Iam.Inputs;
using Policy = Amazon.Auth.AccessControlPolicy.Policy;
using Principal = Amazon.Auth.AccessControlPolicy.Principal;
using Resource = Amazon.Auth.AccessControlPolicy.Resource;
using Statement = Amazon.Auth.AccessControlPolicy.Statement;

namespace Infra.Pulumi.Resources;

public class StressTestClientReadRole : ComponentResource
{
    [Output] public Output<string> StressTestClientReadRoleName { get; set; }

    public StressTestClientReadRole(string clientS3BucketArn, string clientLogsS3BucketArn, ComponentResourceOptions options)
        : base("cube:aws:StressTestClientReadRole", "StressTestClientReadRole", options)
    {
        var role = new Role("StressTestClientReadRole", new RoleArgs
        {
            Name = "StressTestClientReadRole",
            AssumeRolePolicy = new Policy().WithStatements(
                    new Statement(Statement.StatementEffect.Allow)
                        .WithId("AssumeRoleStatement")
                        .WithActionIdentifiers(new ActionIdentifier("sts:AssumeRole"))
                        .WithPrincipals(new Principal(Principal.SERVICE_PROVIDER, "ec2.amazonaws.com")))
                .ToJson(),
            InlinePolicies = new InputList<RoleInlinePolicyArgs>
            {
                new RoleInlinePolicyArgs
                {
                    Name = "S3",
                    Policy = new Policy().WithStatements(
                            new Statement(Statement.StatementEffect.Allow)
                                .WithId("statement1")
                                .WithActionIdentifiers(new ActionIdentifier("s3:GetObject"))
                                .WithResources(new Resource($"{clientS3BucketArn}/*")),
                            new Statement(Statement.StatementEffect.Allow)
                                .WithId("statement2")
                                .WithActionIdentifiers(new ActionIdentifier("s3:PutObject"))
                                .WithResources(new Resource($"{clientLogsS3BucketArn}/*")))
                        .ToJson(),
                }
            },
        }, new CustomResourceOptions { Parent = this });

        StressTestClientReadRoleName = role.Name;
        RegisterOutputs();
    }
}