using System;
using System.Collections.Generic;
using System.Linq;
using System.Threading.Tasks;
using Pulumi;
using Aws = Pulumi.Aws;
using Pulumi.Automation;

namespace InlineProgram
{
    class Program
    {
        static async Task Main(string[] args)
        {

            var program = PulumiFn.Create(() =>
            {
                var euProvider = new Aws.Provider("euProvider", new()
                {
                    Region = "eu-west-1",
                });
                var usProvider = new Aws.Provider("usProvider", new()
                {
                    Region = "us-east-1",
                });

                var USGameApp = new GameApp("test", new CustomResourceOptions
                {
                    Provider = usProvider,
                });

                var EUGameApp = new GameApp("eutest", new CustomResourceOptions
                {
                    Provider = euProvider,
                });

                return new Dictionary<string, object?>
                {
                    ["us_api_url"] = USGameApp.GatewayURL,
                    ["eu_api_url"] = EUGameApp.GatewayURL,
                };

            });

            Console.WriteLine("hi");
            var destroy = args.Any() && args[0] == "destroy";

            var projectName = "stresstest-loader-cicd";
            var stackName = "stresstest-loader-cicd";
            var stackArgs = new InlineProgramArgs(projectName, stackName, program);
            var stack = await LocalWorkspace.CreateOrSelectStackAsync(stackArgs);
            Console.WriteLine("successfully initialized stack");

            Console.WriteLine("installing plugins...");
            await stack.Workspace.InstallPluginAsync("aws", "v4.24.1");
            Console.WriteLine("plugins installed");

            Console.WriteLine("setting up config...");
            await stack.SetConfigAsync("aws:region", new ConfigValue("us-west-2"));
            Console.WriteLine("config set");
            Console.WriteLine("refreshing stack...");
            await stack.RefreshAsync(new RefreshOptions { OnStandardOutput = Console.WriteLine });
            Console.WriteLine("refresh complete");
            if (destroy)
            {
                Console.WriteLine("destroying stack...");
                await stack.DestroyAsync(new DestroyOptions { OnStandardOutput = Console.WriteLine });
                Console.WriteLine("stack destroy complete");
            }
            else
            {
                Console.WriteLine("updating stack...");
                var result = await stack.UpAsync(new UpOptions { OnStandardOutput = Console.WriteLine });

                if (result.Summary.ResourceChanges != null)
                {
                    Console.WriteLine("update summary:");
                    foreach (var change in result.Summary.ResourceChanges)
                        Console.WriteLine($"    {change.Key}: {change.Value}");
                }
            }
        }

    }
}
