using Pulumi;
using Aws = Pulumi.Aws;
using System.IO;
using System.Text.Json;

namespace Infra.Pulumi.Resources;

class GetPublicIp : ComponentResource
{
    [Output] public Output<string> WebAutoScalingGroupName { get; set; }
    [Output] public Output<string> KeyPairName { get; set; }

    public GetPublicIp(string name, Input<string> amiId, ComponentResourceOptions opts = null) : base("stl:aws:GetPublicIp", name, null)
    {
        var ec2Instances = Aws.Ec2.GetInstances.Invoke(new()
        {
            Filters = new[]
            {
                new Aws.Ec2.Inputs.GetInstancesFilterInputArgs
                {
                    Name = "image-id",
                    Values = new[]
                    {
                        amiId.Apply(x => x!.ToString())!
                    }
                }
            }
        });
   
        var ec2PublicIps = ec2Instances.Apply(instance => instance.PublicIps);
        var instancesList = new List<List<Dictionary<string, string>>>();

        // Iterate over the instances and add their public IP addresses to the list
        ec2PublicIps.Apply(ips =>
        {
            foreach (var publicIp in ips)
            {

                // Create a dictionary for each instance with the "public_ip" property
                var instanceDict = new Dictionary<string, string>
                {
                    { "public_ip", publicIp }
                };

                // Add the instance dictionary to the instances list
                instancesList.Add(new List<Dictionary<string, string>> { instanceDict });
            }
      
            // Convert the instances to JSON
            var json = JsonSerializer.Serialize(instancesList);

            // Save the JSON to a file
            var filePath = "/tmp/IP.json";
            File.WriteAllText(filePath, json);

            Console.WriteLine($"IP addresses of EC2 instances saved to '{filePath}'.");

            return ips;
        });
    }
}

