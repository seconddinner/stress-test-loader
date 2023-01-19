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
        // var basic_dynamodb_table = new Aws.DynamoDB.Table("basic-dynamodb-table", new()
        // {
        //     Attributes = new[]
        //     {
        //     new Aws.DynamoDB.Inputs.TableAttributeArgs
        //     {
        //         Name = "UserId",
        //         Type = "S",
        //     },
        //     new Aws.DynamoDB.Inputs.TableAttributeArgs
        //     {
        //         Name = "GameTitle",
        //         Type = "S",
        //     },
        //     new Aws.DynamoDB.Inputs.TableAttributeArgs
        //     {
        //         Name = "TopScore",
        //         Type = "N",
        //     },
        // },
        //     BillingMode = "PAY_PER_REQUEST",
        //     GlobalSecondaryIndexes = new[]
        //     {
        //     new Aws.DynamoDB.Inputs.TableGlobalSecondaryIndexArgs
        //     {
        //         HashKey = "GameTitle",
        //         Name = "GameTitleIndex",
        //         NonKeyAttributes = new[]
        //         {
        //             "UserId",
        //         },
        //         ProjectionType = "INCLUDE",
        //         RangeKey = "TopScore",
        //         ReadCapacity = 1,
        //         WriteCapacity = 1,
        //     },
        // },
        //     HashKey = "UserId",
        //     RangeKey = "GameTitle",
        //     Tags =
        // {
        //     { "Environment", "production" },
        //     { "Name", "dynamodb-table-1" },
        // },
        //     Ttl = new Aws.DynamoDB.Inputs.TableTtlArgs
        //     {
        //         AttributeName = "TimeToExist",
        //         Enabled = false,
        //     },
        // });

        var CubeGameTable = new Aws.DynamoDB.Table("testTable", new()
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
            Name = "testTable",
            PointInTimeRecovery = new Aws.DynamoDB.Inputs.TablePointInTimeRecoveryArgs
            {
                Enabled = true,
            },
            RangeKey = "SortKey",
            Tags =
        {
            { "CubeResourceType", "Latest" },
        },
            Ttl = new Aws.DynamoDB.Inputs.TableTtlArgs
            {
                AttributeName = "Ttl",
                Enabled = true,
            },
        });

        var role = new Aws.Iam.Role("role", new()
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
        });

        var fn = new Aws.Lambda.Function("fn", new()
        {
            Runtime = "python3.9",
            Handler = "handler.handler",
            Role = role.Arn,
            Code = new FileArchive("./function"),
        });

        var api = new AwsApiGateway.RestAPI("api", new()
        {
            Routes =
        {
            new AwsApiGateway.Inputs.RouteArgs
            {
                Path = "/",
                LocalPath = "www",
            },
            new AwsApiGateway.Inputs.RouteArgs
            {
                Path = "/date",
                Method = AwsApiGateway.Method.GET,
                EventHandler = fn,
            },
        },
        });
        this.ApiEndpoint = api.Url;



        var lambdaFunction = new Aws.Lambda.Function("PulumiWebApiGateway_LambdaFunction", new Aws.Lambda.FunctionArgs
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
        });
        System.Console.WriteLine(Aws.Iam.ManagedPolicy.AWSLambdaBasicExecutionRole.ToString());


        var httpApiGateway = new Pulumi.Aws.ApiGatewayV2.Api("PulumiWebApiGateway_ApiGateway", new Pulumi.Aws.ApiGatewayV2.ApiArgs
        {
            ProtocolType = "HTTP",
            RouteSelectionExpression = "${request.method} ${request.path}",
        });

        var httpApiGateway_LambdaIntegration = new Pulumi.Aws.ApiGatewayV2.Integration("PulumiWebApiGateway_ApiGatewayIntegration", new Pulumi.Aws.ApiGatewayV2.IntegrationArgs
        {
            ApiId = httpApiGateway.Id,
            IntegrationType = "AWS_PROXY",
            IntegrationMethod = "POST",
            IntegrationUri = lambdaFunction.Arn,
            PayloadFormatVersion = "2.0",
            TimeoutMilliseconds = 30000,
        });

        var httpApiGatewayRoute = new Pulumi.Aws.ApiGatewayV2.Route("PulumiWebApiGateway_ApiGatewayRoute", new Pulumi.Aws.ApiGatewayV2.RouteArgs
        {
            ApiId = httpApiGateway.Id,
            RouteKey = "$default",
            Target = httpApiGateway_LambdaIntegration.Id.Apply(id => $"integrations/{id}"),
        });

        var httpApiGatewayStage = new Pulumi.Aws.ApiGatewayV2.Stage("PulumiWebApiGateway_ApiGatewayStage", new Pulumi.Aws.ApiGatewayV2.StageArgs
        {
            ApiId = httpApiGateway.Id,
            AutoDeploy = true,
            Name = "$default",
        });

        var lambdaPermissionsForApiGateway = new Aws.Lambda.Permission("PulumiWebApiGateway_LambdaPermission", new Aws.Lambda.PermissionArgs
        {
            Action = "lambda:InvokeFunction",
            Function = lambdaFunction.Name,
            Principal = "apigateway.amazonaws.com",
            SourceArn = Output.Format($"{httpApiGateway.ExecutionArn}/*") // note it's the ExecutionArn.
                                                                          // SourceArn = httpApiGateway.ExecutionArn.Apply(arn => $"{arn}/*") // this is another way of doing the same thing
        });

        this.dotnetApiEndpoint = httpApiGateway.ApiEndpoint.Apply(endpoint => $"{endpoint}/api/values");
    }

    [Output] public Output<string> ApiEndpoint { get; set; }

    [Output] public Output<string> dotnetApiEndpoint { get; set; }
}