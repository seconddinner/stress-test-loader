using Pulumi.Automation;

namespace Infra.Pulumi;

public static class DeployHelpers
{
    public static async Task<int> UpdatePulumiAsync(PulumiFn program, string projectName, string stackName, bool noRefresh, CancellationToken cancellationToken)
    {
        var stack = await SetupAsync(projectName, stackName, program, noRefresh, cancellationToken);

        try
        {
            await stack.UpAsync(new UpOptions
            {
                OnStandardOutput = Console.WriteLine,
                OnStandardError = Console.Error.WriteLine,
                LogToStdErr = true,
                Parallel = int.MaxValue,
                Color = "always",
            }, cancellationToken);
        }
        catch
        {
            return 1;
        }

        return 0;
    }

    public static async Task<int> PreviewPulumiAsync(PulumiFn program, string projectName, string stackName, bool noRefresh, CancellationToken cancellationToken)
    {
        var stack = await SetupAsync(projectName, stackName, program, noRefresh, cancellationToken);

        try
        {
            await stack.PreviewAsync(new()
            {
                OnStandardOutput = Console.WriteLine,
                OnStandardError = Console.Error.WriteLine,
                LogToStdErr = true,
                Parallel = int.MaxValue,
                Color = "always",
                Diff = true,
            }, cancellationToken);
        }
        catch
        {
            return 1;
        }

        return 0;
    }

    public static async Task<int> DestroyPulumiAsync(PulumiFn program, string projectName, string stackName, bool noRefresh, CancellationToken cancellationToken)
    {
        var stack = await SetupAsync(projectName, stackName, program, noRefresh, cancellationToken);

        try
        {
            await stack.DestroyAsync(new DestroyOptions
            {
                OnStandardOutput = Console.WriteLine,
                OnStandardError = Console.Error.WriteLine,
                LogToStdErr = true,
                Parallel = int.MaxValue,
                Color = "always",
            }, cancellationToken);

            await stack.Workspace.RemoveStackAsync(stackName, cancellationToken);

            return 0;
        }
        catch
        {
            return 1;
        }
    }

    private static async Task<WorkspaceStack> SetupAsync(string projectName, string stackName, PulumiFn program, bool noRefresh, CancellationToken cancellationToken)
    {
        var stackArgs = new InlineProgramArgs(projectName, stackName, program);
        var plugins = new Dictionary<string, string>
        {
            { "aws", "v4.24.1" },
        };
        string sourceAmiId = Environment.GetEnvironmentVariable("source_ami_id");
        string localPublicIp = Environment.GetEnvironmentVariable("stress_test_loader_allowed_cidr");
        string publicKey = Environment.GetEnvironmentVariable("public_key");
        string environment = Environment.GetEnvironmentVariable("environment") ?? "stresstest-git-action";
        var config = new Dictionary<string, ConfigValue>
        {
            { "aws:region", new ConfigValue("us-west-2") },
            { "aws-native:region", new ConfigValue("us-west-2") },
            { "stress-test-loader-pulumi:az_count", new ConfigValue("2") },
            { "stress-test-loader-pulumi:cidr_block", new ConfigValue("10.10.0.0/16") },
            { "stress-test-loader-pulumi:desired_capacity", new ConfigValue("1") },
            { "stress-test-loader-pulumi:up_scaling_adjustment", new ConfigValue("-1") },
            { "stress-test-loader-pulumi:down_scaling_adjustment", new ConfigValue("-1") },
            { "stress-test-loader-pulumi:egress_allowed_cidr", new ConfigValue("0.0.0.0/0") },
            { "stress-test-loader-pulumi:environment", new ConfigValue(environment) },
            { "stress-test-loader-pulumi:iam_name", new ConfigValue("stress_test_client_read_profile-stress_test_loader") },
            { "stress-test-loader-pulumi:instance_type", new ConfigValue("t4g.nano") },
            { "stress-test-loader-pulumi:max_size", new ConfigValue("3") },
            { "stress-test-loader-pulumi:min_size", new ConfigValue("1") },
            { "stress-test-loader-pulumi:name", new ConfigValue("stl-stress_test_loader-cluster") },
            { "stress-test-loader-pulumi:prometheus_node_allowed_cidr", new ConfigValue("0.0.0.0/0") },
            { "stress-test-loader-pulumi:public_ip_on_launch", new ConfigValue("true") },
            // NOTE: make public key, ami id, and public ip env var.
            { "stress-test-loader-pulumi:source_ami_id", new ConfigValue(sourceAmiId) },
            { "stress-test-loader-pulumi:source_ami_region", new ConfigValue("us-west-2") },
            { "stress-test-loader-pulumi:ssh_allowed_cidr", new ConfigValue($"[{localPublicIp}]") },
            { "stress-test-loader-pulumi:stress_test_loader_port", new ConfigValue("9005") },
            { "stress-test-loader-pulumi:telegraf_password", new ConfigValue("") },
            { "stress-test-loader-pulumi:telegraf_url", new ConfigValue("") },
            { "stress-test-loader-pulumi:telegraf_username", new ConfigValue("admin") },
            { "stress-test-loader-pulumi:public_key", new ConfigValue(publicKey) },
            { "stress-test-loader-pulumi:stress_test_loader_allowed_cidr", new ConfigValue(localPublicIp) },
        };

        var stack = await LocalWorkspace.CreateOrSelectStackAsync(stackArgs, cancellationToken);
        Console.WriteLine("successfully initialized stack");

        Console.WriteLine("installing plugins...");
        foreach (var kvp in plugins)
        {
            await stack.Workspace.InstallPluginAsync(kvp.Key, kvp.Value, cancellationToken: cancellationToken);
        }
        Console.WriteLine("plugins installed");

        Console.WriteLine("setting up config...");
        foreach (var kvp in config)
        {
            await stack.SetConfigAsync(kvp.Key, kvp.Value, cancellationToken: cancellationToken);
        }
        Console.WriteLine("config set");

        if (!noRefresh)
        {
            Console.WriteLine("refreshing stack...");
            await stack.RefreshAsync(new RefreshOptions { OnStandardOutput = Console.WriteLine }, cancellationToken);
            Console.WriteLine("refresh complete");
        }

        return stack;
    }
}
