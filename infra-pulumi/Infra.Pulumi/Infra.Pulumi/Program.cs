using System.ComponentModel.DataAnnotations;
using Infra.Pulumi.Resources;
using McMaster.Extensions.CommandLineUtils;
using Pulumi;
using Pulumi.Automation;
using Pulumi.Aws;
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
        
        var requiredEnvList = new List<string> { "stress_test_loader_allowed_cidr", "public_key",
            "s3_client_bucket_name", "s3_log_bucket_name" };
        foreach (var env in requiredEnvList)
        {
            if (Environment.GetEnvironmentVariable(env) == null)
            {
                Console.WriteLine($"[Error] Need to set the \"{env}\" environment variable");
                throw new NullReferenceException();
            }
        }
        
        var desiredCapacity = Environment.GetEnvironmentVariable("desired_capacity") ?? "2";
        var regions = Environment.GetEnvironmentVariable("regions") ?? "us-west-2";

        var cfg = new StressConfig
        {
            Environment = ProjectName,
            DesiredCapacity = int.Parse(desiredCapacity),
            PublicKey = Environment.GetEnvironmentVariable("public_key")!,
            AllowedCidrBlocks = Environment.GetEnvironmentVariable("stress_test_loader_allowed_cidr")!.Split(",").ToList(),
            S3ClientBucketName = Environment.GetEnvironmentVariable("s3_client_bucket_name")!,
            S3LogBucketName = Environment.GetEnvironmentVariable("s3_log_bucket_name")!,
            Regions = regions
        };
        
        // Delete IP.json if it exists
        var filePath = Environment.GetEnvironmentVariable("ip_json_path") ?? "/tmp/IP.json";
        if (File.Exists(filePath))
        {
            File.Delete(filePath);
        }

        #region Deploy
        var program = PulumiFn.Create(async () =>
        {
            var s3 = new S3(cfg);
            
            var regionList = regions.Split(',').ToList();
            
            foreach (var r in regionList)
            {
                var region = r.Trim();
                var provider = new Provider(region, new()
                {
                    Region = region,
                });
                cfg.CurrentRegion = region;
                var ami = new Ami(cfg, new ComponentResourceOptions { Provider = provider });
                var iam = new Iam(cfg, new ComponentResourceOptions { Provider = provider });
                var vpc = new Vpc(cfg, new ComponentResourceOptions { Provider = provider });
                var autoscaling = new Autoscaling(ami.AmiId,
                    iam.StressTestClientReadProfileName,vpc.MainVpcId, vpc.MainSubnetIds, cfg, 
                    new ComponentResourceOptions { Provider = provider });
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


