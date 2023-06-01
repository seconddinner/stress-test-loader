using System.ComponentModel.DataAnnotations;
using Infra.Pulumi.Resources;
using McMaster.Extensions.CommandLineUtils;
using Pulumi;
using Pulumi.Automation;
using Pulumi.Aws;
using Config = Pulumi.Config;

#pragma warning disable CS1998

namespace Infra.Pulumi;

[Command("deploy", Description = "Deploy Pulumi stack")]
public class DeployPulumiCommand
{
    public static Task<int> Main(string[] args) => CommandLineApplication.ExecuteAsync<DeployPulumiCommand>(args);

    #region Command Options
    [Option("--destroy",
        "Destroy stack (optional)",
        CommandOptionType.NoValue)]
    private bool Destroy { get; set; }

    [Option("--no-refresh",
        "Do not refresh state store with current live state",
        CommandOptionType.NoValue)]
    private bool NoRefresh { get; set; }

    [Option("--preview",
        "Only preview changes, but do not make them",
        CommandOptionType.NoValue)]
    private bool Preview { get; set; }

    [Required]
    [Option("--project-name",
        "The project name in pulumi",
        CommandOptionType.SingleValue)]
    private string ProjectName { get; set; } = null!;

    [Required]
    [Option("--stack-name",
        "The stack name in pulumi",
        CommandOptionType.SingleValue)]
    private string StackName { get; set; } = null!;
    #endregion

    public async Task<int> OnExecuteAsync(CancellationToken cancellationToken)
    {
        var optimizationVars = new Dictionary<string, string>
        {
            {"PULUMI_EXPERIMENTAL", "1"},
            {"PULUMI_SKIP_CHECKPOINTS", "true"},
            {"PULUMI_OPTIMIZED_CHECKPOINT_PATCH", "true"},
        };
        foreach (var (k,v) in optimizationVars)
        {
            Environment.SetEnvironmentVariable(k, v);
        }

        const string clientS3BucketArn = "arn:aws:s3:::stresstest-client";
        const string clientLogsS3BucketArn = "arn:aws:s3:::stresstest-client-log";

        #region Deploy
        var program = PulumiFn.Create(async () =>
        {
            var config = new Config();
            List<string> regionList = new List<string>();
            if (config.Get("regions") != null)
            {
                string regions = config.Get("regions")!;
                // Remove the brackets and any whitespace
                string trimmedInput = regions.Trim('[', ']', ' ');
                // Split the string by commas to get individual elements
                string[] elements = trimmedInput.Split(',');
                regionList = new List<string>(elements);
                for (int i = 0; i < regionList.Count; i++)
                {
                    regionList[i] = regionList[i].Trim(' ');
                }
            }
            else
            {
                regionList.Add("us-west-2");
            }
            
            foreach (var region in regionList)
            {
                var provider = new Provider(region, new()
                {
                    Region = region,
                });

                var ami = new Ami($"stl-ami-{region}", provider, region);
                var iam = new Iam($"stl-iam-{region}", provider, region);
                var vpc = new Vpc($"stl-vpc-{region}", provider, region);
                var autoscaling = new Autoscaling($"stl-autoscaling-{region}", provider, region, ami.AmiId,
                    iam.StressTestClientReadProfileName,vpc.MainVpcId, vpc.DefaultSecurityGroupId, 
                    vpc.MainSubnetIds);
            }
        });

        if (Destroy)
        {
            return await DeployHelpers.DestroyPulumiAsync(program, ProjectName, StackName, NoRefresh, cancellationToken);
        }

        if (Preview)
        {
            return await DeployHelpers.PreviewPulumiAsync(program, ProjectName, StackName, NoRefresh, cancellationToken);
        }

        return await DeployHelpers.UpdatePulumiAsync(program, ProjectName, StackName, NoRefresh, cancellationToken);
        #endregion
    }


}


