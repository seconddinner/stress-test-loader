using System.Collections.Generic;
using System.Text.Json;
using Pulumi;
using Aws = Pulumi.Aws;
using AwsApiGateway = Pulumi.AwsApiGateway;

return await Deployment.RunAsync(() =>
{

    // add comments to trigger build
    var basic_dynamodb_table = new Aws.DynamoDB.Table("basic-dynamodb-table", new()
    {
        Attributes = new[]
        {
            new Aws.DynamoDB.Inputs.TableAttributeArgs
            {
                Name = "UserId",
                Type = "S",
            },
            new Aws.DynamoDB.Inputs.TableAttributeArgs
            {
                Name = "GameTitle",
                Type = "S",
            },
            new Aws.DynamoDB.Inputs.TableAttributeArgs
            {
                Name = "TopScore",
                Type = "N",
            },
        },
        BillingMode = "PAY_PER_REQUEST",
        GlobalSecondaryIndexes = new[]
        {
            new Aws.DynamoDB.Inputs.TableGlobalSecondaryIndexArgs
            {
                HashKey = "GameTitle",
                Name = "GameTitleIndex",
                NonKeyAttributes = new[]
                {
                    "UserId",
                },
                ProjectionType = "INCLUDE",
                RangeKey = "TopScore",
                ReadCapacity = 1,
                WriteCapacity = 1,
            },
        },
        HashKey = "UserId",
        RangeKey = "GameTitle",
        Tags =
        {
            { "Environment", "production" },
            { "Name", "dynamodb-table-1" },
        },
        Ttl = new Aws.DynamoDB.Inputs.TableTtlArgs
        {
            AttributeName = "TimeToExist",
            Enabled = false,
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

    return new Dictionary<string, object?>
    {
        ["url"] = api.Url,
    };
});

