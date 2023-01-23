using Pulumi;
using System.Collections.Generic;
using System.Text.Json;
using Aws = Pulumi.Aws;


public class ServiceDeploymentArgs
{
    public Input<string> SomeInput { get; set; } = null!;
}

public class GameApp : Pulumi.ComponentResource
{
    public Output<string> GatewayURL { get; private set; } = null!;
    public GameApp(string name, ServiceDeploymentArgs opts, CustomResourceOptions? cropts = null, ComponentResourceOptions? options = null)
        : base("examples:aws:GameApp", name, options)
    {
        var TestTable = new Aws.DynamoDB.Table($"{name}-testTable", new()
        {
            Attributes = new[]
            {
                new Aws.DynamoDB.Inputs.TableAttributeArgs
                {
                    Name = "PartitionKey",
                    Type = "S",
                },
                new Aws.DynamoDB.Inputs.TableAttributeArgs
                {
                    Name = "SortKey",
                    Type = "S",
                },
            },
            BillingMode = "PAY_PER_REQUEST",
            HashKey = "PartitionKey",
            PointInTimeRecovery = new Aws.DynamoDB.Inputs.TablePointInTimeRecoveryArgs
            {
                Enabled = true,
            },
            RangeKey = "SortKey",
            Tags =
            {
                { "ResourceType", "Latest" },
            },
            Ttl = new Aws.DynamoDB.Inputs.TableTtlArgs
            {
                AttributeName = "Ttl",
                Enabled = true,
            },
        }, cropts);


        var role = new Aws.Iam.Role($"{name}-role", new()
        {
            AssumeRolePolicy = JsonSerializer.Serialize(new Dictionary<string, object?>
            {
                ["Version"] = "2012-10-17",
                ["Statement"] = new[]
                {
                new Dictionary<string, object?>
                {
                    ["Action"] = "sts:AssumeRole",
                    ["Effect"] = "Allow",
                    ["Principal"] = new Dictionary<string, object?>
                    {
                        ["Service"] = "lambda.amazonaws.com",
                    },
                },
            },
            }),
            ManagedPolicyArns = new[]
            {
            "arn:aws:iam::aws:policy/service-role/AWSLambdaBasicExecutionRole",
        },
        }, cropts);


        var fn = new Aws.Lambda.Function($"{name}-fn", new()
        {
            Runtime = "python3.9",
            Handler = "handler.handler",
            Role = role.Arn,
            Code = new FileArchive("../lambda-py-function"),
        }, cropts);

        // can't api rest api gateway? 
        // Argument type 'Pulumi.CustomResourceOptions' is not assignable to parameter type 'Pulumi.ComponentResourceOptions?'
        var api = new Pulumi.AwsApiGateway.RestAPI("api", new()
        {
            Routes =
            {
                new AwsApiGateway.Inputs.RouteArgs
                {
                    Path = "/",
                    LocalPath = "../www",
                },
                new AwsApiGateway.Inputs.RouteArgs
                {
                    Path = "/date",
                    Method = AwsApiGateway.Method.GET,
                    EventHandler = fn,
                },
            },
        }, cropts);

        var lambdaFunction = new Aws.Lambda.Function($"{name}-PulumiWebApiGateway_LambdaFunction", new Aws.Lambda.FunctionArgs
        {
            Architectures = new[]
            {
                "arm64",
            },
            EphemeralStorage = new Aws.Lambda.Inputs.FunctionEphemeralStorageArgs
            {
                Size = 512,
            },
            Handler = "WebAPILambda::WebAPILambda.LambdaEntryPoint::FunctionHandlerAsync",
            MemorySize = 1024,
            Publish = false,
            ReservedConcurrentExecutions = -1,
            Role = role.Arn,
            Runtime = Aws.Lambda.Runtime.Dotnet6,
            Timeout = 10,
            Code = new FileArchive("../WebAPILambda/bin/Release/net6.0/linux-arm64/lambda.zip"),
            TracingConfig = new Aws.Lambda.Inputs.FunctionTracingConfigArgs
            {
                Mode = "PassThrough",
            },
        }, cropts);



        var httpApiGateway = new Pulumi.Aws.ApiGatewayV2.Api($"{name}-PulumiWebApiGateway_ApiGateway", new Pulumi.Aws.ApiGatewayV2.ApiArgs
        {
            ProtocolType = "HTTP",
            RouteSelectionExpression = "${request.method} ${request.path}",
        }, cropts);

        var httpApiGateway_LambdaIntegration = new Pulumi.Aws.ApiGatewayV2.Integration($"{name}-PulumiWebApiGateway_ApiGatewayIntegration", new Pulumi.Aws.ApiGatewayV2.IntegrationArgs
        {
            ApiId = httpApiGateway.Id,
            IntegrationType = "AWS_PROXY",
            IntegrationMethod = "POST",
            IntegrationUri = lambdaFunction.Arn,
            PayloadFormatVersion = "2.0",
            TimeoutMilliseconds = 30000,
        }, cropts);

        var httpApiGatewayRoute = new Pulumi.Aws.ApiGatewayV2.Route($"{name}-PulumiWebApiGateway_ApiGatewayRoute", new Pulumi.Aws.ApiGatewayV2.RouteArgs
        {
            ApiId = httpApiGateway.Id,
            RouteKey = "$default",
            Target = httpApiGateway_LambdaIntegration.Id.Apply(id => $"integrations/{id}"),
        }, cropts);

        var httpApiGatewayStage = new Pulumi.Aws.ApiGatewayV2.Stage($"{name}-PulumiWebApiGateway_ApiGatewayStage", new Pulumi.Aws.ApiGatewayV2.StageArgs
        {
            ApiId = httpApiGateway.Id,
            AutoDeploy = true,
            Name = "$default",
        }, cropts);

        var lambdaPermissionsForApiGateway = new Aws.Lambda.Permission($"{name}-PulumiWebApiGateway_LambdaPermission", new Aws.Lambda.PermissionArgs
        {
            Action = "lambda:InvokeFunction",
            Function = lambdaFunction.Name,
            Principal = "apigateway.amazonaws.com",
            SourceArn = Output.Format($"{httpApiGateway.ExecutionArn}/*") // note it's the ExecutionArn.
                                                                          // SourceArn = httpApiGateway.ExecutionArn.Apply(arn => $"{arn}/*") // this is another way of doing the same thing
        }, cropts);

        this.GatewayURL = httpApiGateway.ApiEndpoint.Apply(endpoint => $"{endpoint}/api/values");
        this.RegisterOutputs();
    }
}