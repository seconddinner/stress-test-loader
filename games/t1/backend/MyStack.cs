using System.Collections.Generic;
using System.Text.Json;
using Pulumi;
using Aws = Pulumi.Aws;
using AwsApiGateway = Pulumi.AwsApiGateway;

class MyStack : Stack
{
    public MyStack()
    {

        // add comments to trigger build

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

        var eugameApp = new GameApp("eutest", new CustomResourceOptions
        {
            Provider = euProvider,
        });
        this.TestURL = USGameApp.GatewayURL;
        this.EuTestURL = eugameApp.GatewayURL;
    }

    [Output] public Output<string> TestURL { get; set; }

    [Output] public Output<string> EuTestURL { get; set; }
}