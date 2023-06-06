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

        var config = new Dictionary<string, ConfigValue>
        {
            { "aws:region", new ConfigValue("us-west-2") },
            { "aws-native:region", new ConfigValue("us-west-2") },
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
